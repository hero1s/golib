package ip_geo

import (
	"errors"
	"git.moumentei.com/plat_go/golib/log"
	geoIp2 "github.com/oschwald/geoip2-golang"
	"net"
)

type IpCountry struct {
	IsoCode string            `json:"isoCode"`
	Names   map[string]string `json:"names"`
}
type IpCity struct {
	GeoNameID uint              `json:"geoNameId"`
	Names     map[string]string `json:"names"`
}
type IpInfo struct {
	Ip        string    `json:"ip"`
	Country   IpCountry `json:"country"`
	City      IpCity    `json:"city"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	TimeZone  string    `json:"timeZone"`
}
type GeoClient struct {
	IpInfoReader *geoIp2.Reader
}

var (
	geoClient *GeoClient
)

// 导入GeoLite2-City.mmdb
func InitGeoIp2Client(geoFilePath string) error {
	geoClient = &GeoClient{}
	var err error
	geoClient.IpInfoReader, err = geoIp2.Open(geoFilePath)
	if err != nil {
		log.Errorf("初始化IP数据库错误:%v", err)
	}
	return err
}

func GeoIp2GetIpInfo(ipStr string) (ipInfo *IpInfo, err error) {
	ip := net.ParseIP(ipStr)
	ipInfo = &IpInfo{
		Ip: ipStr,
	}
	if geoClient == nil {
		err = errors.New("客户端未初始化")
		return
	}
	var record *geoIp2.City
	record, err = geoClient.IpInfoReader.City(ip)
	if err != nil {
		return
	}
	ipInfo.Ip = ipStr
	ipInfo.Country = IpCountry{IsoCode: record.Country.IsoCode, Names: record.Country.Names}
	ipInfo.City = IpCity{GeoNameID: record.City.GeoNameID, Names: record.City.Names}
	ipInfo.Latitude = record.Location.Latitude
	ipInfo.Longitude = record.Location.Longitude
	ipInfo.TimeZone = record.Location.TimeZone
	return
}
