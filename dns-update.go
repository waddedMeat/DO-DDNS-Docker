package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/digitalocean/godo"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	changed, currentIP, err := getIPChange(os.Getenv("LOG_PATH") + "/log.txt")

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}

	if changed != false {
		fmt.Println("no change")
		os.Exit(0)
	}

	if err := updateDNS(os.Getenv("TOKEN"), os.Getenv("DOMAIN"), strings.Split(os.Getenv("NAMES"), ","), currentIP); err != nil {
		fmt.Println(err)
		os.Exit(255)
	}
}

func updateDNS(token string, domain string, names []string, currentIP *string) error {

	if token == "" || domain == "" || len(names) == 0 {
		return errors.New("missing token, domain, and/or names")
	}

	client := godo.NewFromToken(token)

	for _, name := range names {
		ctx := context.Background()

		records, _, err := client.Domains.RecordsByName(ctx, domain, name, &godo.ListOptions{})
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, record := range records {
			request := createDomainRecordEditRequest(record)
			record.Data = *currentIP

			_, _, err := client.Domains.EditRecord(ctx, domain, record.ID, request)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func getIPChange(logPath string) (bool, *string, error) {

	if logPath == "" {
		return false, nil, errors.New("missing log path")
	}

	data, _ := ioutil.ReadFile(logPath)
	lastIP := string(data)

	currentIP := getCurrentIP()

	if lastIP != currentIP {
		err := ioutil.WriteFile(logPath, []byte(currentIP), 0644)
		if err != nil {
			fmt.Println(err)
		}
		return true, &currentIP, nil
	}

	return false, nil, nil
}

func createDomainRecordEditRequest(record godo.DomainRecord) *godo.DomainRecordEditRequest {
	return &godo.DomainRecordEditRequest{
		Type:     record.Type,
		Name:     record.Name,
		Data:     record.Data,
		Priority: record.Priority,
		Port:     record.Port,
		TTL:      record.TTL,
		Weight:   record.Weight,
		Flags:    record.Flags,
		Tag:      record.Tag,
	}
}

//plucked from https://gist.github.com/ankanch/8c8ec5aaf374039504946e7e2b2cdf7f
func getCurrentIP() string {
	url := "https://api.ipify.org?format=text"
	// we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api

	fmt.Printf("Getting IP address from  ipify ...\n")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My IP is:%s\n", ip)
	return string(ip)
}
