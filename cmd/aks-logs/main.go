package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mirage20/aks-logs/pkg/client"
)

var (
	subscriptionId    string
	resourceGroupName string
	clusterName       string
	namespace         string
	workloadName      string
	containerName     string
	contains          string
	startTime         string
	endTime           string

	maxRecords int

	showDescending bool
	showQuery      bool
	showRaw        bool
	showVersion    bool

	versionString string
)

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println("aks-logs", versionString)
		return
	}

	check(validateFlags())

	var err error
	var start, end time.Time
	if len(startTime) > 0 {
		start, err = time.Parse(time.RFC3339Nano, startTime)
		check(err)
	}
	if len(endTime) > 0 {
		end, err = time.Parse(time.RFC3339Nano, endTime)
		check(err)
	}

	c, err := client.NewClientWithCLIAuthorizer()
	check(err)

	req := client.QueryRequest{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		Namespace:         namespace,
		WorkloadName:      workloadName,
		ContainerName:     containerName,
		StartTime:         start,
		EndTime:           end,
		FilterContains: func() []string {
			if len(strings.TrimSpace(contains)) > 0 {
				return strings.Split(contains, ",")
			}
			return nil
		}(),
		MaxRecords:     maxRecords,
		ShowQuery:      showQuery,
		ShowDescending: showDescending,
	}
	ctx := context.Background()
	resp, err := c.Query(ctx, &req)
	check(err)
	for _, logEntry := range resp.LogEntries {
		if showRaw {
			b, _ := json.Marshal(logEntry)
			fmt.Println(string(b))
		} else {
			fmt.Println(logEntry.LogEntry)
		}
	}
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func validateFlags() error {
	if len(subscriptionId) == 0 {
		return fmt.Errorf("subscriptionId cannot be empty")
	}
	if len(resourceGroupName) == 0 {
		return fmt.Errorf("resourceGroupName cannot be empty")
	}
	if len(clusterName) == 0 {
		return fmt.Errorf("clusterName cannot be empty")
	}
	if len(namespace) == 0 {
		return fmt.Errorf("namespace cannot be empty")
	}
	if len(workloadName) == 0 {
		return fmt.Errorf("workloadName cannot be empty")
	}
	if len(containerName) == 0 {
		return fmt.Errorf("containerName cannot be empty")
	}
	return nil
}

func init() {
	flag.StringVar(&subscriptionId, "subscriptionId", "", "Subscription id of the AKS Cluster. Required")
	flag.StringVar(&resourceGroupName, "resourceGroupName", "", "Resource group name of the AKS Cluster. Required")
	flag.StringVar(&clusterName, "clusterName", "", "AKS Cluster name. Required")
	flag.StringVar(&namespace, "namespace", "", "Kubernetes namespace of the workload. Required")
	flag.StringVar(&workloadName, "workloadName", "", "Kubernetes workload name. Required")
	flag.StringVar(&containerName, "containerName", "", "Container name of the Kubernetes workload. Required")
	flag.StringVar(&startTime, "startTime", "", "Start time of the log entries. Default to (current time - 24h) if not specified. Format: 2006-01-02T15:04:05.999999999+07:00")
	flag.StringVar(&endTime, "endTime", "", "End time of the log entries. Default to (current time) if not specified. Format: 2006-01-02T15:04:05.999999999+07:00")
	flag.StringVar(&contains, "contains", "", "Filter log entries by list of contents if specified. Example: -contains=val1,val2")
	flag.IntVar(&maxRecords, "maxRecords", 0, "Maximum number of log entries to query. Default to 1000")
	flag.BoolVar(&showDescending, "showDescending", false, "Output the logs in descending order based on generated time")
	flag.BoolVar(&showQuery, "showQuery", false, "Output the Kusto query")
	flag.BoolVar(&showRaw, "showRaw", false, "Output the raw results")
	flag.BoolVar(&showVersion, "version", false, "Output version information")
}
