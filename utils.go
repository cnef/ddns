package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	publicIPURL     = "http://members.3322.org/dyndns/getip"
	recordIPURL     = "http://119.29.29.29/d?dn=%s.%s"
	recordUpdateURL = "https://dnsapi.cn/Record.Modify"
	recordListURL   = "https://dnsapi.cn/Record.List"
)

func trim(str string) string {
	return strings.TrimSpace(str)
}

func getCurrentIP() (string, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(publicIPURL)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return trim(string(body)), nil
}

func getRecordIP(record, domain string) (string, error) {
	url := fmt.Sprintf(recordIPURL, record, domain)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return trim(string(body)), nil
}

func updateRecordIP(token, record, domain, recordID, publicIP string) error {

	vals := url.Values{}
	vals.Add("login_token", token)
	vals.Add("sub_domain", record)
	vals.Add("domain", domain)
	vals.Add("record_id", recordID)
	vals.Add("value", publicIP)
	vals.Add("record_type", "A")
	vals.Add("format", "json")
	vals.Add("record_line_id", "0")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.PostForm(recordUpdateURL, vals)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if strings.Contains(string(body), "successful") {
		return nil
	}

	return errors.New(string(body))
}

func listRecords(token, domain string) (string, error) {

	vals := url.Values{}
	vals.Add("login_token", token)
	vals.Add("domain", domain)
	vals.Add("format", "json")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.PostForm(recordListURL, vals)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return string(body), nil
}

func getRecordID(token, record, domain string) (recordID string, err error) {

	data, err := listRecords(token, domain)
	if err != nil {
		return "", fmt.Errorf("get domain records error: %v", err)
	}

	records := &RecordList{}
	err = json.Unmarshal([]byte(data), records)
	if err != nil {
		return "", fmt.Errorf("json unmarshal domain records error: %v", err)
	}

	recordID, err = func() (string, error) {
		for _, r := range records.Records {
			if r.Name == record {
				if r.Type == "A" {
					return r.ID, nil
				}
				return "", fmt.Errorf("record %s type is %s, type must be A", r.Name, r.Type)
			}
		}
		return "", fmt.Errorf("don't found record %s on domain %s, you need create it manually", record, domain)
	}()

	if err != nil {
		return "", err
	}

	return recordID, nil
}
