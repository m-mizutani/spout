package server

import (
	"net/http"

	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/usecase"
)

func getLogs(ctx *model.Context, uc *usecase.Usecase, r *http.Request) (*httpResponse, error) {
	var options []usecase.ExportLogsOption

	if v := r.URL.Query().Get("limit"); v != "" {
		options = append(options, usecase.WithLimit(v))
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		options = append(options, usecase.WithOffset(v))
	}
	if v := r.URL.Query().Get("query"); v != "" {
		options = append(options, usecase.WithQuery(v))
	}

	logs, err := uc.ExportLogs(ctx, options...)
	if err != nil {
		return nil, err
	}

	return &httpResponse{
		code: http.StatusOK,
		data: logs,
	}, nil
}
