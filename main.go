package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/imroc/req/v3"
)

var (
	gl2WhileAPI string
)

func init() {
	gl2WhileAPI = os.Getenv("GL2_WHILE_API")
}

func readText() []string {

	data, err := os.ReadFile("cli.txt")

	if err != nil {
		panic(err)
	}

	println(string(data))

	items := strings.Split(string(data), "\n")

	return items

}

const (
	NC_CLI_IM   = "nc.cli.im"
	CLI_IM      = "cli.im"
	USER_CLI_IM = "user.cli.im"
)

func main() {

	urlPathMap := make(map[string][]string)

	urlPathMap[NC_CLI_IM] = []string{}
	urlPathMap[CLI_IM] = []string{}
	urlPathMap[USER_CLI_IM] = []string{}

	urls := readText()
	uniqueMap := make(map[string]struct{})

	for _, ele := range urls {

		urlObject, err := url.Parse(ele)
		if err != nil {
			panic(err)
		}

		switch urlObject.Host {
		case NC_CLI_IM:
			if isUnique(uniqueMap, urlObject.Host+urlObject.Path) {
				urlPathMap[NC_CLI_IM] = append(urlPathMap[NC_CLI_IM], urlObject.Path)
			}
		case CLI_IM:
			if isUnique(uniqueMap, urlObject.Host+urlObject.Path) {
				urlPathMap[CLI_IM] = append(urlPathMap[CLI_IM], urlObject.Path)
			}
		case USER_CLI_IM:

			if isUnique(uniqueMap, urlObject.Host+urlObject.Path) {
				urlPathMap[USER_CLI_IM] = append(urlPathMap[USER_CLI_IM], urlObject.Path)
			}
		}
	}

	SyncGL2Api(urlPathMap)

	fmt.Println(len(urlPathMap[NC_CLI_IM]))
	fmt.Println(len(urlPathMap[CLI_IM]))
	fmt.Println(len(urlPathMap[USER_CLI_IM]))

}

func isUnique(uniqueMap map[string]struct{}, key string) bool {

	if _, ok := uniqueMap[key]; ok {
		return false
	} else {
		uniqueMap[key] = struct{}{}
		return true
	}

}

func SyncGL2Api(urlPathMap map[string][]string) {
	for app, uris := range urlPathMap {
		for _, uri := range uris {
			err := DealWithAppPair(uri, app)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

		}
	}

}

func DealWithAppPair(uri string, app string) error {

	client := req.NewClient()

	resp, err := client.NewRequest().SetQueryParams(map[string]string{
		"app": app,
		"uri": uri,
	}).Get(gl2WhileAPI + "/api/uri/app_uri/verify")

	if err != nil {
		return err
	}

	if resp.GetStatusCode() != 200 {
		return errors.New("status code is not 200")
	}

	bizResp := &BizResp{}

	err = json.Unmarshal(resp.Bytes(), resp)

	if err != nil {
		return err
	}

	if !bizResp.Success {
		resp, err := client.NewRequest().SetBody(AddUriPair{
			App: app,
			Uri: uri,
		}).Post(gl2WhileAPI + "/api/uri/app_uri")

		if err != nil {
			return err
		}

		if resp.GetStatusCode() != 200 {
			return errors.New("status code is not 200")
		}

		bizResp = &BizResp{}

		err = json.Unmarshal(resp.Bytes(), resp)

		if err != nil {
			return err
		}

		if !bizResp.Success {
			return fmt.Errorf("add app url error,app: %s , url: %s", app, uri)
		}
	}

	return nil
}

type BizResp struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

type AddUriPair struct {
	App string `json:"app"`
	Uri string `json:"uri"`
}
