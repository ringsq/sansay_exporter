// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/ringsq/sansay_exporter/models"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hooklift/gowsdl/soap"
	"github.com/prometheus/client_golang/prometheus"
)

var realtimeMetrics = []string{"NumOrig",
	"NumTerm",
	"Cps",
	"NumPeak",
	"TotalCLZ",
	"NumCLZCps",
	"TotalLimit",
	"CpsLimit"}

type Sansay struct {
	XMLName  xml.Name `xml:"mysqldump"`
	Text     string   `xml:",chardata"`
	Database struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		Table []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			Row  []struct {
				Text  string `xml:",chardata"`
				Field []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"field"`
			} `xml:"row"`
		} `xml:"table"`
	} `xml:"database"`
}

type XBMediaServerRealTimeStatList struct {
	XMLName                   xml.Name `xml:"XBMediaServerRealTimeStatList"`
	Text                      string   `xml:",chardata"`
	XBMediaServerRealTimeStat []struct {
		Text              string `xml:",chardata"`
		MediaSrvIndex     string `xml:"mediaSrvIndex"`
		PublicIP          string `xml:"publicIP"`
		MaxConnections    string `xml:"maxConnections"`
		Priority          string `xml:"priority"`
		Alias             string `xml:"alias"`
		SwitchType        string `xml:"switchType"`
		Status            string `xml:"status"`
		NumActiveSessions string `xml:"numActiveSessions"`
	} `xml:"XBMediaServerRealTimeStat"`
}

var trunkfields = map[string]string{
	"1st15mins_call_attempt":     "Fifteen_Calls_Attempt",
	"1st15mins_call_answer":      "Fifteen_Calls_Answer",
	"1st15mins_call_fail":        "Fifteen_Calls_Fail",
	"1h_call_attempt":            "Hour_Calls_Attempt",
	"1h_call_answer":             "Hour_Calls_Answer",
	"1h_call_fail":               "Hour_Calls_Fail",
	"24h_call_attempt":           "Day_Calls_Attempt",
	"24h_call_answer":            "Day_Calls_Answer",
	"24h_call_fail":              "Day_Calls_Fail",
	"1st15mins_call_durationSec": "Fifteen_Duration",
	"1h_call_durationSec":        "Hour_Duration",
	"24h_call_durationSec":       "Day_Duration",
	"1st15mins_pdd_ms":           "Fifteen_PDD",
	"1h_pdd_ms":                  "Hour_PDD",
	"24h_pdd_ms":                 "Day_PDD",
}
var resourceMetrics = make([]string, 0, len(trunkfields))

type Trunk struct {
	TrunkId               string
	Alias                 string
	Fqdn                  string
	NumOrig               string
	NumTerm               string
	Cps                   string
	NumPeak               string
	TotalCLZ              string
	NumCLZCps             string
	TotalLimit            string
	CpsLimit              string
	Fifteen_Calls_Attempt string
	Fifteen_Calls_Answer  string
	Fifteen_Calls_Fail    string
	Hour_Calls_Attempt    string
	Hour_Calls_Answer     string
	Hour_Calls_Fail       string
	Day_Calls_Attempt     string
	Day_Calls_Answer      string
	Day_Calls_Fail        string
	Fifteen_Duration      string
	Hour_Duration         string
	Day_Duration          string
	Fifteen_PDD           string
	Hour_PDD              string
	Day_PDD               string
	Direction             string
}
type collector struct {
	target     string
	targetPath string
	username   string
	password   string
	logger     log.Logger
	useSoap    bool
}

func init() {
	for _, value := range trunkfields {
		resourceMetrics = append(resourceMetrics, value)
	}
}

// Describe implements Prometheus.Collector.
func (c collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

// Collect implements Prometheus.Collector.
func (c collector) Collect(ch chan<- prometheus.Metric) {
	paths := []string{"stats/realtime", "stats/resource", "stats/media_server", "download/resource"}
	var wg sync.WaitGroup
	var err error
	start := time.Now()
	results := make(chan interface{})
	defer close(results)
	for _, path := range paths {
		wg.Add(1)
		go ScrapeTarget(c, path, results, &wg)
	}
	for i := 0; i < len(paths); i++ {
		result := <-results
		switch obj := result.(type) {
		case Sansay:
			err = nil
			c.processCollection(ch, obj)
		case XBMediaServerRealTimeStatList:
			err = nil
			c.processMediaCollection(ch, obj)
		case models.XBResourceList:
			err = nil
			c.processXBResourceList(ch, obj)
		case error:
			err = obj
		default:
			err = errors.New("Invalid type returned from target")
		}
		if err != nil {
			level.Info(c.logger).Log("msg", "Error scraping target", "err", err)
			ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
		}
	}
	wg.Wait()
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("sansay_scrape_duration_seconds", "Total sansay time scrape took (walk and processing).", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())

}

// processMediaCollection creates the metrics for the media server statistics.  The media server stats are
// a totally different format than then other endpoints.
func (c collector) processMediaCollection(ch chan<- prometheus.Metric, media XBMediaServerRealTimeStatList) {
	for _, mediaServer := range media.XBMediaServerRealTimeStat {
		var msType string
		words := strings.Split(mediaServer.SwitchType, " ")
		msType = words[len(words)-1]
		if strings.LastIndex(msType, "-") > 0 {
			msType = msType[strings.LastIndex(msType, "-")+1:]
		}
		labels := []string{"server", "server_ip", "type"}
		labelValues := []string{mediaServer.Alias, mediaServer.PublicIP, msType}
		status := "0"
		if mediaServer.Status == "up" {
			status = "1"
		}
		addLabeledMetric(ch, "mediaserver_up", status, labels, labelValues)
		addLabeledMetric(ch, "mediaserver_sessions_limit", mediaServer.MaxConnections, labels, labelValues)
		addLabeledMetric(ch, "mediaserver_sessions", mediaServer.NumActiveSessions, labels, labelValues)
	}
}

// processXBResourceList creates the metrics for the resource configurations.
func (c collector) processXBResourceList(ch chan<- prometheus.Metric, resources models.XBResourceList) {
	labels := []string{"trunkgroup", "alias"}
	var labelValues []string
	for _, resource := range resources.XBResource {
		labelValues = []string{resource.TrunkId, resource.Name}
		addLabeledMetric(ch, "config_trunk_sessions_max", resource.Capacity, labels, labelValues)
		addLabeledMetric(ch, "config_trunk_cps_max", resource.CpsLimit, labels, labelValues)
	}
}

func (c collector) processCollection(ch chan<- prometheus.Metric, sansay Sansay) {
	for _, table := range sansay.Database.Table {
		var direction string
		switch table.Name {
		case "system_stat":
			for _, row := range table.Row {
				for _, field := range row.Field {
					switch field.Name {
					case "ha_pre_state":
					case "ha_current_state":
					default:
						addMetric(ch, field.Name, field.Text)
					}
				}
			}
		case "XBResourceRealTimeStatList":
			for _, row := range table.Row {
				trunk := Trunk{}
				for _, field := range row.Field {
					err := setField(&trunk, field.Name, field.Text)
					if err != nil {
						ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
					}
				}
				if trunk.Fqdn == "Group" {
					err := addTrunkMetrics(ch, trunk, realtimeMetrics)
					if err != nil {
						ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
					}
				}
			}
			// Resource tables
		case "ingress_stat":
			direction = "ingress"
			fallthrough
		case "gw_egress_stat":
			if direction == "" {
				direction = "egress"
			}
			for _, row := range table.Row {
				trunk := Trunk{}
				trunk.Direction = direction
				for _, field := range row.Field {
					if field.Name == "trunk_id" {
						field.Name = "trunkId"
					}
					fieldName, ok := trunkfields[field.Name]
					if !ok {
						fieldName = field.Name
					}
					err := setField(&trunk, fieldName, field.Text)
					if err != nil {
						ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
					}
				}
				err := addTrunkMetrics(ch, trunk, resourceMetrics)
				if err != nil {
					ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
				}
			}
		}
	}
}

// ScrapeTarget scrapes the Sansay API
func ScrapeTarget(c collector, path string, result chan<- interface{}, wg *sync.WaitGroup) {
	logger := c.logger
	var obj interface{}
	var sansay Sansay
	var media XBMediaServerRealTimeStatList
	var resourceList models.XBResourceList
	var body []byte
	var err error

	if c.useSoap {
		body, err = callSoapAPI(c, path)
		if err != nil {
			result <- err
			wg.Done()
			return
		}
	} else {
		body, err = callRestAPI(c, path)
		if err != nil {
			result <- err
			wg.Done()
			return
		}
	}
	if strings.HasSuffix(path, "media_server") {
		err = xml.Unmarshal(body, &media)
		obj = media
	} else if strings.HasSuffix(path, "download/resource") {
		err = xml.Unmarshal(body, &resourceList)
		obj = resourceList
	} else {
		err = xml.Unmarshal(body, &sansay)
		obj = sansay
	}
	if err != nil {
		level.Error(logger).Log("msg", "Error parsing XML", "path", path, "err", err)
		result <- err
		wg.Done()
		return
	}
	result <- obj
	wg.Done()
	return
}

func callRestAPI(c collector, path string) ([]byte, error) {
	username := c.username
	password := c.password
	logger := c.logger
	target := fmt.Sprintf("%s%s%s", c.target, c.targetPath, path)
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	_, err := url.Parse(target)
	if err != nil {
		level.Error(logger).Log("msg", "Could not parse target URL", "err", err)
		return nil, err
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", target, http.NoBody)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating HTTP request", "err", err)
		return nil, err
	}

	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)

	if err != nil {
		level.Error(logger).Log("msg", "Error for HTTP request", "err", err)
		return nil, err
	}
	level.Info(logger).Log("msg", "Received HTTP response", "status_code", resp.StatusCode)
	if resp.StatusCode == 404 {
		resp.Body.Close()
		return callSoapAPI(c, path)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 300 {
		err = fmt.Errorf("Invalid response from server: %d", resp.StatusCode)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		level.Info(logger).Log("msg", "Failed to read HTTP response body", "err", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

// callSoapAPI makes a SOAP call to the Sansay SBC -- used for older OS versions
func callSoapAPI(c collector, path string) ([]byte, error) {
	var err error
	var response []byte
	var statName string

	// Determine the stat name by splitting the path
	paths := strings.Split(path, "/")
	if len(paths) == 1 {
		statName = paths[0]
	} else {
		statName = paths[len(paths)-1]
	}

	target := fmt.Sprintf("%s%s", c.target, "/SSConfig/SansayWS")
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}
	client := soap.NewClient(target, soap.WithTLS(&tls.Config{InsecureSkipVerify: true}))
	service := NewSansayWS(client)
	if strings.HasSuffix(path, "download/resource") {
		params := &DownloadParams{
			Username: c.username,
			Password: c.password,
			Page:     0,
			Table:    "resource",
		}

		if reply, err := service.DoDownloadXmlFile(params); err == nil {
			response = []byte(reply.Xmlfile)
		}
	} else {
		params := &RealTimeStatsParams{
			Username: c.username,
			Password: c.password,
			StatName: statName,
		}
		if reply, err := service.DoRealTimeStats(params); err == nil {
			response = []byte(reply.Xmlfile)
		}
	}
	if err != nil {
		level.Error(c.logger).Log("msg", "Error calling SOAP API", "err", err)
		return nil, err
	}
	return response, nil
}

func addMetric(ch chan<- prometheus.Metric, name string, value string) error {
	metricName := fmt.Sprintf("sansay_%s", name)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", nil, nil),
		prometheus.GaugeValue,
		floatValue)
	return nil
}

func addLabeledMetric(ch chan<- prometheus.Metric, name string, value string, labels []string, labelValues []string) error {
	metricName := fmt.Sprintf("sansay_%s", name)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(metricName, "", labels, nil),
		prometheus.GaugeValue,
		floatValue, labelValues...)
	return nil
}

func addTrunkMetrics(ch chan<- prometheus.Metric, trunk Trunk, metricNames []string) error {
	for _, metric := range metricNames {
		baseName := strings.ToLower(metric)
		metricName := fmt.Sprintf("sansay_trunk_%s", baseName)

		value, err := getField(&trunk, metric)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
			continue
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
			continue
		}
		//fmt.Printf("New Metric: %s TG=%s Alias=%s\n", metricName, trunk.TrunkId, trunk.Alias)
		labels := []string{"trunkgroup", "alias"}
		labelValues := []string{trunk.TrunkId, trunk.Alias}
		if trunk.Direction != "" {
			labels = append(labels, "direction")
			labelValues = append(labelValues, trunk.Direction)
			fieldName := strings.Split(baseName, "_")
			if fieldName[1] == "calls" {
				labels = append(labels, "status")
				labelValues = append(labelValues, fieldName[2])
				metricName = metricName[0:strings.LastIndex(metricName, "_")]
			}
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(metricName, "", labels, nil),
			prometheus.GaugeValue,
			floatValue, labelValues...)
	}
	return nil
}

// setField sets field of v with given name to given value.
func setField(v interface{}, name string, value string) error {
	// v must be a pointer to a struct
	nme := []rune(name)
	nme[0] = unicode.ToUpper(nme[0])
	name = string(nme)
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return nil // fmt.Errorf("not a field name: %s", name)
	}

	// Field must be exported
	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}

	// We expect a string field
	if fv.Kind() != reflect.String {
		return fmt.Errorf("%s is not a string field", name)
	}

	// Set the value
	fv.SetString(value)
	return nil
}

// setField sets field of v with given name to given value.
func getField(v interface{}, name string) (string, error) {
	// v must be a pointer to a struct
	nme := []rune(name)
	nme[0] = unicode.ToUpper(nme[0])
	name = string(nme)
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return "", errors.New("v must be pointer to struct")
	}

	// Dereference pointer
	rv = rv.Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return "", fmt.Errorf("not a field name: %s", name)
	}

	// We expect a string field
	if fv.Kind() != reflect.String {
		return "", fmt.Errorf("%s is not a string field", name)
	}

	return fv.String(), nil
}
