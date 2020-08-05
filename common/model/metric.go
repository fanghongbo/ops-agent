package model

import "fmt"

type MetricValue struct {
	Endpoint    string      `json:"endpoint"`
	Metric      string      `json:"metric"`
	Value       interface{} `json:"value"`
	Step        int64       `json:"step"`
	CounterType string      `json:"counterType"`
	Tags        string      `json:"tags"`
	Timestamp   int64       `json:"timestamp"`
}

func (u *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%v>",
		u.Endpoint,
		u.Metric,
		u.CounterType,
		u.Tags,
		u.Step,
		u.Timestamp,
		u.Value,
	)
}
