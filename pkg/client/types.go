package client

import "time"

type QueryRequest struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string

	Namespace     string
	WorkloadName  string
	ContainerName string

	StartTime time.Time
	EndTime   time.Time

	FilterContains []string

	MaxRecords     int
	ShowQuery      bool
	ShowDescending bool
}

type QueryResponse struct {
	LogEntries []LogEntry
}

type LogEntry struct {
	TenantId       string    `col:"TenantId"`
	SourceSystem   string    `col:"SourceSystem"`
	TimeGenerated  time.Time `col:"TimeGenerated"`
	Computer       string    `col:"Computer"`
	TimeOfCommand  time.Time `col:"TimeOfCommand"`
	ContainerID    string    `col:"ContainerID"`
	Image          string    `col:"Image"`
	ImageTag       string    `col:"ImageTag"`
	Repository     string    `col:"Repository"`
	Name           string    `col:"Name"`
	LogEntry       string    `col:"LogEntry"`
	LogEntrySource string    `col:"LogEntrySource"`
	Type           string    `col:"Type"`
	ResourceId     string    `col:"_ResourceId"`
}

func (r *QueryRequest) Default() {
	now := time.Now()

	if r.EndTime.IsZero() {
		r.EndTime = now
	}

	if r.StartTime.IsZero() {
		r.StartTime = r.EndTime.Add(-time.Hour * 24)
	}

	if r.MaxRecords == 0 {
		r.MaxRecords = 1000
	}
}
