package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

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

// func
