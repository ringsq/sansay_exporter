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

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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
	target   string
	username string
	password string
	logger   log.Logger
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
	paths := []string{"realtime", "resource"}
	var wg sync.WaitGroup
	var err error
	var sansay Sansay
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
		case error:
			err = obj
		case Sansay:
			err = nil
			sansay = obj
		default:
			err = errors.New("Invalid type returned from target")
		}

		if err != nil {
			level.Info(c.logger).Log("msg", "Error scraping target", "err", err)
			ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("sansay_error", "Error scraping target", nil, nil), err)
		}
		c.processCollection(ch, sansay)
	}
	wg.Wait()
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc("sansay_scrape_duration_seconds", "Total sansay time scrape took (walk and processing).", nil, nil),
		prometheus.GaugeValue,
		time.Since(start).Seconds())

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
	target := fmt.Sprintf("%s%s", c.target, path)
	username := c.username
	password := c.password
	logger := c.logger
	var sansay Sansay

	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	_, err := url.Parse(target)
	if err != nil {
		level.Error(logger).Log("msg", "Could not parse target URL", "err", err)
		result <- err
		wg.Done()
		return
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", target, http.NoBody)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating HTTP request", "err", err)
		result <- err
		wg.Done()
		return
	}

	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)

	if err != nil {
		level.Error(logger).Log("msg", "Error for HTTP request", "err", err)
		result <- err
		wg.Done()
		return
	}
	level.Info(logger).Log("msg", "Received HTTP response", "status_code", resp.StatusCode)
	if resp.StatusCode > 300 {
		result <- fmt.Errorf("Invalid response from server: %d", resp.StatusCode)
		wg.Done()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		level.Info(logger).Log("msg", "Failed to read HTTP response body", "err", err)
		result <- err
		wg.Done()
		return
	}
	err = xml.Unmarshal(body, &sansay)
	if err != nil {
		level.Error(logger).Log("msg", "Error parsing XML", "err", err)
		result <- err
		wg.Done()
		return
	}
	result <- sansay
	wg.Done()
	return
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
