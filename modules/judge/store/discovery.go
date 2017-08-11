package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

// Discovery  根据 metric tags 等信息判断是否需要调用自动发现API
func Discovery(item *model.JudgeItem) {

	discovery := false

	url := g.Config().Discovery.API
	timeout := time.Duration(g.Config().Discovery.Timeout) * time.Millisecond

	endpoint := item.Endpoint
	metric := item.Metric
	tags := item.Tags
	// value := item.Value
	// timestamp := item.Timestamp

	tag := ""

	//enabled := g.Config().Discovery.Enabled
	debug := g.Config().Debug
	endpointList := g.Config().Discovery.Endpoints
	metricList := g.Config().Discovery.Metrics
	tagList := g.Config().Discovery.Tags

	// 如果 endpoints 列表为空，则标识自动发现所有的 endpoint
	if len(endpointList) == 0 {
		discovery = true
	}

	// 如果 endpoints 列表为空，则标识自动发现所有的 metric
	if len(metricList) == 0 {
		discovery = true
	}

	// 如果 endpoints 列表为空，则标识自动发现所有的 tag
	if len(tagList) == 0 {
		discovery = true
	}

	for _, e := range endpointList {
		if strings.Contains(endpoint, e) {
			if debug {
				fmt.Println("discovery endpoint: " + item.Endpoint + " match " + e)
			}
			discovery = true
		}

	}

	for _, m := range metricList {
		if strings.Contains(metric, m) {
			if debug {
				fmt.Println("discovery metric: " + item.Metric + " match " + m)
			}
			discovery = true
		}

	}

	for _, t := range tagList {
		for tagKey := range tags {
			if strings.Contains(tagKey, t) {
				tag = tagKey
				if debug {
					fmt.Println("discovery tag : " + tagKey + " match " + t)
				}
				discovery = true
			}
		}

	}

	if discovery {
		if debug {
			fmt.Printf("discovery summary  : endpint: %s  metric: %s tag: %s \n", endpoint, metric, tag)

		}

		resp, _ := CallDiscoveryAPI(url, timeout, item)

		if debug && resp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("discovery response : %s \n", string(body))
		}

	}

}

// CallDiscoveryAPI 调用自动发现API
func CallDiscoveryAPI(url string, timeout time.Duration, item *model.JudgeItem) (resp *http.Response, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("cannot connect to discovery url: %s \n", url)
			//fmt.Println(err)
		}
	}()

	req := new(bytes.Buffer)
	json.NewEncoder(req).Encode(item)

	bodyType := "application/json;charset=utf-8"

	client := http.Client{
		Timeout: timeout,
	}
	resp, err = client.Post(url, bodyType, req)

	defer resp.Body.Close()

	return resp, err

}
