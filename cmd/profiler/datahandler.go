package main

import (
	"strconv"

	tsprofilerModels "github.com/cha87de/tsprofiler/models"
	drModels "github.com/disresc/lib/models"
	drUtil "github.com/disresc/lib/util"
)

func getCPUData(event *drModels.Event) (tsprofilerModels.TSInputMetric, bool) {
	cpuTotalItem, found := drUtil.FindEventItem(event, "kvmtop-cpu", "cpu_total")
	if !found {
		return tsprofilerModels.TSInputMetric{}, false
	}
	cpuTotal, err := strconv.Atoi(cpuTotalItem.GetValue())
	if err != nil {
		return tsprofilerModels.TSInputMetric{}, false
	}
	cpuStealItem, found := drUtil.FindEventItem(event, "kvmtop-cpu", "cpu_steal")
	if !found {
		return tsprofilerModels.TSInputMetric{}, false
	}
	cpuSteal, err := strconv.Atoi(cpuStealItem.GetValue())
	if err != nil {
		return tsprofilerModels.TSInputMetric{}, false
	}

	util := cpuTotal + cpuSteal
	min := 0
	max := 100

	return tsprofilerModels.TSInputMetric{
		Name:     name,
		Value:    float64(util),
		FixedMin: float64(min),
		FixedMax: float64(max),
	}, true
}
