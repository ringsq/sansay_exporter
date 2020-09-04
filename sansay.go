package main

import (
	"encoding/xml"
	"github.com/hooklift/gowsdl/soap"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type UploadParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com uploadParams"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type UploadResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com uploadResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Version string `xml:"version,omitempty"`
}

type ReplaceLargeParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com replaceLargeParams"`

	Binfile []byte `xml:"binfile,omitempty"`

	Table string `xml:"table,omitempty"`

	TableID int32 `xml:"tableID,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type ReplaceResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com replaceResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Version string `xml:"version,omitempty"`
}

type UpdateParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com updateParams"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type UpdateLargeParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com updateLargeParams"`

	Binfile []byte `xml:"binfile,omitempty"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type UpdateResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com updateResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DeleteParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com deleteParams"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DeleteLargeParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com deleteLargeParams"`

	Binfile []byte `xml:"binfile,omitempty"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DeleteResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com deleteResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DownloadParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com downloadParams"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Page int32 `xml:"page,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DownloadLargeParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com downloadLargeParams"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DownloadResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com downloadResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	HasMore int32 `xml:"hasMore,omitempty"`

	Version string `xml:"version,omitempty"`
}

type DownloadLargeResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com downloadLargeResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Binfile []byte `xml:"binfile,omitempty"`

	Version string `xml:"version,omitempty"`
}

type QueryParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com queryParams"`

	Table string `xml:"table,omitempty"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	Page int32 `xml:"page,omitempty"`

	QueryString string `xml:"queryString,omitempty"`

	Version string `xml:"version,omitempty"`
}

type QueryResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com queryResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	HasMore int32 `xml:"hasMore,omitempty"`

	Version string `xml:"version,omitempty"`
}

type RoutelookupParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com routelookupParams"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	QueryString string `xml:"queryString,omitempty"`

	Version string `xml:"version,omitempty"`
}

type RoutelookupResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com routelookupResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	Version string `xml:"version,omitempty"`
}

type RealTimeStatsParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com realTimeStatsParams"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	StatName string `xml:"statName,omitempty"`

	Version string `xml:"version,omitempty"`
}

type RealTimeStatsResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com realTimeStatsResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	Version string `xml:"version,omitempty"`
}

type SystemStatsParams struct {
	XMLName xml.Name `xml:"http://ws.sansay.com SystemStatsParams"`

	Username string `xml:"username,omitempty"`

	Password string `xml:"password,omitempty"`

	SysStatName string `xml:"sysStatName,omitempty"`

	Version string `xml:"version,omitempty"`
}

type SystemStatsResult struct {
	XMLName xml.Name `xml:"http://ws.sansay.com SystemStatsResult"`

	RetCode int32 `xml:"retCode,omitempty"`

	Msg string `xml:"msg,omitempty"`

	Xmlfile string `xml:"xmlfile,omitempty"`

	Version string `xml:"version,omitempty"`
}

type SansayWS interface {
	DoUploadXmlFile(request *UploadParams) (*UploadResult, error)

	DoReplaceLarge(request *ReplaceLargeParams) (*ReplaceResult, error)

	DoDelete(request *DeleteParams) (*DeleteResult, error)

	DoDeleteLarge(request *DeleteLargeParams) (*DeleteResult, error)

	DoUpdate(request *UpdateParams) (*UpdateResult, error)

	DoUpdateLarge(request *UpdateLargeParams) (*UpdateResult, error)

	DoDownloadXmlFile(request *DownloadParams) (*DownloadResult, error)

	DoDownloadLargeXmlFile(request *DownloadLargeParams) (*DownloadLargeResult, error)

	DoQueryXmlFile(request *QueryParams) (*QueryResult, error)

	DoRouteLookup(request *RoutelookupParams) (*RoutelookupResult, error)

	DoRealTimeStats(request *RealTimeStatsParams) (*RealTimeStatsResult, error)

	DoSystemStats(request *SystemStatsParams) (*SystemStatsResult, error)
}

type sansayWS struct {
	client *soap.Client
}

func NewSansayWS(client *soap.Client) SansayWS {
	return &sansayWS{
		client: client,
	}
}

func (service *sansayWS) DoUploadXmlFile(request *UploadParams) (*UploadResult, error) {
	response := new(UploadResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoReplaceLarge(request *ReplaceLargeParams) (*ReplaceResult, error) {
	response := new(ReplaceResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoDelete(request *DeleteParams) (*DeleteResult, error) {
	response := new(DeleteResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoDeleteLarge(request *DeleteLargeParams) (*DeleteResult, error) {
	response := new(DeleteResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoUpdate(request *UpdateParams) (*UpdateResult, error) {
	response := new(UpdateResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoUpdateLarge(request *UpdateLargeParams) (*UpdateResult, error) {
	response := new(UpdateResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoDownloadXmlFile(request *DownloadParams) (*DownloadResult, error) {
	response := new(DownloadResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoDownloadLargeXmlFile(request *DownloadLargeParams) (*DownloadLargeResult, error) {
	response := new(DownloadLargeResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoQueryXmlFile(request *QueryParams) (*QueryResult, error) {
	response := new(QueryResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoRouteLookup(request *RoutelookupParams) (*RoutelookupResult, error) {
	response := new(RoutelookupResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoRealTimeStats(request *RealTimeStatsParams) (*RealTimeStatsResult, error) {
	response := new(RealTimeStatsResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *sansayWS) DoSystemStats(request *SystemStatsParams) (*SystemStatsResult, error) {
	response := new(SystemStatsResult)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
