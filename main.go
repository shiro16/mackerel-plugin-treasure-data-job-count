package main

import (
	"flag"
	"fmt"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	td_client "github.com/treasure-data/td-client-go"
)

type TreasureDataJobCountPlugin struct {
	Prefix string
	ApiKey string
	Status string
}

func (t TreasureDataJobCountPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(t.Prefix)
	return map[string](mp.Graphs){
		t.Prefix: mp.Graphs{
			Label: labelPrefix,
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "seconds", Label: "Seconds"},
			},
		},
	}
}

func (t TreasureDataJobCountPlugin) FetchMetrics() (map[string]interface{}, error) {
	c, err := td_client.NewTDClient(td_client.Settings{
		ApiKey: t.ApiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch uptime metrics: %s", err)
	}

	jobs, err := c.ListJobs()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch uptime metrics: %s", err)
	}

	count := 0
	for _, job := range jobs.ListJobsResultElements {
		if job.Status == t.Status {
			count++
		}
	}
	return map[string]interface{}{"seconds": count}, nil
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "treasure-data-job-count", "Metric key prefix")
	optApiKey := flag.String("treasure-data-api-key", "", "Treasure Data Api Key")
	optStatus := flag.String("job-status", "running", "Count Job Status")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var treasureDataJobCount TreasureDataJobCountPlugin

	treasureDataJobCount.Prefix = *optPrefix
	treasureDataJobCount.ApiKey = *optApiKey
	treasureDataJobCount.Status = *optStatus

	helper := mp.NewMackerelPlugin(treasureDataJobCount)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
