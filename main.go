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
	From   int
	To     int
	Status string
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

	options := &td_client.ListJobsOptions{}
	if t.From >= 0 {
		options.WithFrom(t.From)
	}
	if t.To >= 0 {
		options.WithTo(t.To)
	}
	if t.Status != "" {
		options.WithStatus(t.Status)
	}
	jobs, err := c.ListJobsWithOptions(options)
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch jobs: %s", err)
	}

	stat := map[string]uint32{"error": 0, "killed": 0, "queued": 0, "running": 0, "success": 0}
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
	optFrom := flag.Int("from", -1, "Treasure Data job from the nth index in the list")
	optTo := flag.Int("to", -1, "Treasure Data job up to the nth index in the list")
	optStatus := flag.String("status", "", "Treasure Data job status")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var treasureDataJobCount TreasureDataJobCountPlugin

	treasureDataJobCount.Prefix = *optPrefix
	treasureDataJobCount.ApiKey = *optApiKey
	treasureDataJobCount.From = *optFrom
	treasureDataJobCount.To = *optTo
	treasureDataJobCount.Status = *optStatus

	helper := mp.NewMackerelPlugin(treasureDataJobCount)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
