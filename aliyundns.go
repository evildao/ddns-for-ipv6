package main

import (
	"errors"
	"log"
	"slices"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/spf13/viper"
)

const v6Type = "AAAA"

type alidnsStruct struct {
	client *alidns.Client
}

func NewAlidnsClient() (*alidnsStruct, error) {
	client, err := alidns.NewClientWithAccessKey(
		viper.GetString("RegionId"),
		viper.GetString("AccessKey"),
		viper.GetString("AccessSecret"),
	)
	if err != nil {
		return nil, err
	}
	return &alidnsStruct{
		client: client,
	}, nil
}

func (c *alidnsStruct) SetNewIPV6(ips []string) error {
	if len(ips) == 0 {
		return nil
	}
	oldMap, err := c.GetOldRRID()
	if err != nil {
		return err
	}
	existsIps := make([]string, 0)
	for id, oip := range oldMap {
		if !slices.Contains(ips, oip) {
			if err := c.DeleteRR(id); err != nil {
				log.Printf("记录删除失败 %v", err)
			}
		} else {
			existsIps = append(existsIps, oip)
		}
	}
	for _, ip := range slices.DeleteFunc(ips, func(ip string) bool {
		return slices.Contains(existsIps, ip)
	}) {
		request := alidns.CreateAddDomainRecordRequest()
		request.DomainName = viper.GetString("DomainName")
		request.RR = viper.GetString("RR")
		request.Type = v6Type
		request.Value = ip
		if _, err := c.client.AddDomainRecord(request); err != nil {
			log.Printf("记录添加失败-%v", err)
		}
	}
	return nil
}

func (c *alidnsStruct) GetOldRRID() (map[string]string, error) {
	if viper.GetString("RR") == "" {
		return nil, errors.New("RR Empty")
	}
	listRequest := alidns.CreateDescribeDomainRecordsRequest()
	listRequest.Scheme = "https"
	listRequest.DomainName = viper.GetString("DomainName")
	listRequest.PageSize = "500"
	listRequest.RRKeyWord = viper.GetString("RR")
	listResponse, err := c.client.DescribeDomainRecords(listRequest)
	if err != nil {
		return nil, err
	}
	oldMap := make(map[string]string)
	for _, v := range listResponse.DomainRecords.Record {
		oldMap[v.RecordId] = v.Value
	}
	return oldMap, nil
}

func (c *alidnsStruct) DeleteRR(id string) error {
	request := alidns.CreateDeleteDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = id
	response, err := c.client.DeleteDomainRecord(request)
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return errors.New("删除失败")
	}
	return nil
}
