package main

import (
	"encoding/json"
	"fmt"
	"time"

	tsprofilerApi "github.com/cha87de/tsprofiler/api"
	"github.com/cha87de/tsprofiler/models"
	tsprofilerModels "github.com/cha87de/tsprofiler/models"
	tsprofiler "github.com/cha87de/tsprofiler/profiler"
	drModels "github.com/disresc/lib/models"
	"github.com/micro/go-micro/util/log"
)

var profilers map[string]tsprofilerApi.TSProfiler

func handle(event *drModels.Event) {
	profiler := getProfiler(event.GetSource())
	metrics := getMetricsFromEvent(event)
	tsdata := tsprofilerModels.TSInput{
		Metrics: metrics,
	}
	profiler.Put(tsdata)
}

func getProfiler(source string) tsprofilerApi.TSProfiler {
	if profilers == nil {
		profilers = make(map[string]tsprofilerApi.TSProfiler)
	}
	if _, exists := profilers[source]; !exists {
		// create new profiler
		profilers[source] = tsprofiler.NewProfiler(tsprofilerModels.Settings{
			Name:           source,
			BufferSize:     10, // default: 10, with default 1s frequency => every 10s
			States:         4,  // default: 4
			History:        1,  // default: 1
			FilterStdDevs:  0,
			FixBound:       true,
			OutputFreq:     time.Duration(6),
			OutputCallback: profileOutput,
			PeriodSize:     []int{},
		})
	}
	return profilers[source]
}

func getMetricsFromEvent(event *drModels.Event) []tsprofilerModels.TSInputMetric {
	metrics := make([]tsprofilerModels.TSInputMetric, 0)

	metric, found := getCPUData(event)
	if found {
		metrics = append(metrics, metric)
	}

	return metrics
}

func profileOutput(data models.TSProfile) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Failed to serialize profile: %s", err)
		return
	}
	item := drModels.EventItem{
		Transmitter: name,
		Metric:      "profile",
		Value:       string(json),
	}
	event := drModels.Event{
		Source:    data.Name,
		Timestamp: time.Now().Unix(),
		Items:     []*drModels.EventItem{&item},
	}
	interval := 60
	fmt.Printf("publishing profiler %s", string(json))
	transmitter.Publish(&event, interval)
}
