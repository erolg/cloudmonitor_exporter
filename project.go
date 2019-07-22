package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// datapoint represents the member of Datapoints field from QueryMetricLastResponse
type datapoint struct {
	Average    float64 `json:"Average"`
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Value      float64 `json:"Value"`
	InstanceId string  `json:"instanceId"`
	Timestamp  int64   `json:"timestamp"`
	UserId     string  `json:"userId"`
	Port       string  `json:port`
	Vip        string  `json:vip`
}

// GetResponseFunc returns a function to retrieve queryMetricLast
type GetResponseFunc func(client *cms.Client, request *cms.DescribeMetricLastRequest) string

// Project represents the dashborad from which metrics collected
type Project struct {
	client      *cms.Client
	getResponse GetResponseFunc
	Namespace   string
}

func defaultGetResponseFunc(client *cms.Client, request *cms.DescribeMetricLastRequest) (result string) {
Loop:
	response, err := client.DescribeMetricLast(request)
	if err != nil {
		log.Println("Encounter response error from Aliyun:", err)
		time.Sleep(time.Duration(1) * time.Minute)
		responseError.Inc()
		goto Loop
	}
	result = response.Datapoints
	return
}

func retrieve(metric string, p Project) []datapoint {
	request := cms.CreateDescribeMetricLastRequest()
	request.Namespace = p.Namespace
	request.MetricName = metric

	requestsStats.Inc()
	var source string
	if p.getResponse == nil {
		source = defaultGetResponseFunc(p.client, request)
	} else {
		source = p.getResponse(p.client, request)
	}

	datapoints := make([]datapoint, 10)
	if err := json.Unmarshal([]byte(source), &datapoints); err != nil {
		panic(err)
	}

	return datapoints
}
