package client

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/v1/operationalinsights"
	"github.com/Azure/go-autorest/autorest"
	azureauth "github.com/Azure/go-autorest/autorest/azure/auth"
)

type Client struct {
	QueryClient *operationalinsights.QueryClient
}

func NewClientWithCLIAuthorizer() (*Client, error) {
	authorizer, err := azureauth.NewAuthorizerFromCLI()
	if err != nil {
		return nil, fmt.Errorf("fail to create authorizer from Azure CLI: %v", err)
	}
	c := operationalinsights.NewQueryClient()
	c.Authorizer = authorizer
	return &Client{QueryClient: &c}, nil
}

func (c *Client) Query(ctx context.Context, qReq *QueryRequest) (*QueryResponse, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", qReq.SubscriptionId),
		"resourceGroup":  autorest.Encode("path", qReq.ResourceGroupName),
		"cluster":        autorest.Encode("path", qReq.ClusterName),
	}

	qReq.Default()

	var sb strings.Builder

	const datetimeFormat = "2006-01-02 15:04:05.999999999"
	locUtc, _ := time.LoadLocation("UTC")

	sb.WriteString(fmt.Sprintf("set query_take_max_records=%d;", qReq.MaxRecords))
	sb.WriteString(fmt.Sprintf("set truncationmaxsize=67108864;\n"))
	sb.WriteString(fmt.Sprintf("let ContainerIdList = KubePodInventory\n"))
	sb.WriteString(fmt.Sprintf("| where ClusterName == '%s'\n", qReq.ClusterName))
	sb.WriteString(fmt.Sprintf("| where Namespace == '%s'\n", qReq.Namespace))
	sb.WriteString(fmt.Sprintf("| where Name startswith '%s'\n", qReq.WorkloadName))
	sb.WriteString(fmt.Sprintf("| where ContainerName endswith '%s'\n", qReq.ContainerName))
	sb.WriteString(fmt.Sprintf("| distinct ContainerID;\n"))
	sb.WriteString(fmt.Sprintf("ContainerLog\n"))
	sb.WriteString(fmt.Sprintf("| where ContainerID in (ContainerIdList)\n"))
	for _, v := range qReq.FilterContains {
		sb.WriteString(fmt.Sprintf("| where LogEntry contains '%s'\n", v))
	}
	sb.WriteString(fmt.Sprintf("| where TimeGenerated between (datetime(%s)..datetime(%s))\n", qReq.StartTime.In(locUtc).Format(datetimeFormat), qReq.EndTime.In(locUtc).Format(datetimeFormat)))

	if qReq.ShowDescending {
		sb.WriteString(fmt.Sprintf("| order by TimeGenerated desc"))
	} else {
		sb.WriteString(fmt.Sprintf("| order by TimeGenerated asc"))
	}

	if qReq.ShowQuery {
		fmt.Println()
		fmt.Println(sb.String())
		fmt.Println()
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{cluster}/providers/microsoft.insights/logs?api-version=2018-03-01-preview", pathParameters),
		autorest.WithJSON(map[string]string{
			"query": sb.String(),
		}))

	req, err := preparer.Prepare((&http.Request{}).WithContext(ctx))
	if err != nil {
		return nil, err
	}

	resp, err := c.QueryClient.ExecuteSender(req)
	if err != nil {
		return nil, err
	}

	results, err := c.QueryClient.ExecuteResponder(resp)
	if err != nil {
		return nil, err
	}

	if len(*results.Tables) == 0 {
		return nil, fmt.Errorf("query result does not contain any tables")
	}

	pTable := (*(results.Tables))[0]

	qResp := &QueryResponse{}

	for _, row := range *pTable.Rows {
		l := LogEntry{}
		err := parseLogEntry(&l, *pTable.Columns, row)
		if err != nil {
			return nil, fmt.Errorf("cannot parse log entry: %v", err)
		}
		qResp.LogEntries = append(qResp.LogEntries, l)
	}
	return qResp, nil
}

func parseLogEntry(logEntry *LogEntry, columns []operationalinsights.Column, row []interface{}) error {
	v := reflect.ValueOf(logEntry).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("i must be a pointer")
	}

	tagFieldMapping := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		colName := typeField.Tag.Get("col")
		tagFieldMapping[colName] = i
	}

	for i, col := range columns {
		colName := *col.Name
		colType := *col.Type
		fieldIdx, ok := tagFieldMapping[colName]
		if !ok {
			return fmt.Errorf("cannot find filed for column %q", colName)
		}
		fieldVal := v.Field(fieldIdx)
		if colType == "string" {
			fieldVal.Set(reflect.ValueOf(row[i].(string)))
		} else if colType == "datetime" {
			t, err := time.Parse(time.RFC3339Nano, row[i].(string))
			if err != nil {
				return fmt.Errorf("cannot parse datetime: %v", err)
			}
			fieldVal.Set(reflect.ValueOf(t))
		} else {
			return fmt.Errorf("unknown column type %q", colType)
		}
	}
	return nil
}
