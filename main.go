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
}

func (t TreasureDataJobCountPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(t.Prefix)

	return map[string]mp.Graphs{
		t.Prefix: mp.Graphs{
			Label: labelPrefix,
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "error", Label: "Error"},
				mp.Metrics{Name: "killed", Label: "Killed"},
				mp.Metrics{Name: "queued", Label: "Queued"},
				mp.Metrics{Name: "running", Label: "Running"},
				mp.Metrics{Name: "success", Label: "Success"},
			},
		},
	}
}

func (t TreasureDataJobCountPlugin) FetchMetrics() (map[string]interface{}, error) {
	c, err := td_client.NewTDClient(td_client.Settings{
		ApiKey: t.ApiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("Faild to Treasure Data Connection: %s", err)
	}

	jobs, err := c.ListJobs()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch jobs: %s", err)
	}

	stat := make(map[string]uint32)
	for _, job := range jobs.ListJobsResultElements {
		stat[job.Status]++
	}

	statRet := make(map[string]interface{})
	for key, value := range stat {
		statRet[key] = value
	}

	return statRet, nil
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "treasure-data-job-count", "Metric key prefix")
	optApiKey := flag.String("treasure-data-api-key", "", "Treasure Data Api Key")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var treasureDataJobCount TreasureDataJobCountPlugin

	treasureDataJobCount.Prefix = *optPrefix
	treasureDataJobCount.ApiKey = *optApiKey

	helper := mp.NewMackerelPlugin(treasureDataJobCount)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
