package model

import (
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/genproto/googleapis/logging/v2"
)

type (
	GoogleProjectID string
)

type CloudLoggingLog struct {
	Severity string `json:"severity"`
	Payload  any    `json:"payload"`
	InsertID string `json:"insert_id"`
	LogName  string `json:"log_name"`
	Trace    string `json:"trace"`
	SpanID   string `json:"span_id"`

	Resource       *monitoredres.MonitoredResource `json:"resource"`
	SourceLocation *logging.LogEntrySourceLocation `json:"source_location"`
}
