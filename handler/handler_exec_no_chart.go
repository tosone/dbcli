//go:build no_chart

package handler

import (
	"context"
	"io"

	"github.com/xo/usql/metacmd"
)

// doExecChart executes a single query against the database, displaying its output as a chart.
func (h *Handler) doExecChart(_ context.Context, _ io.Writer, _ metacmd.Option, _, _ string, _ bool, _ []interface{}) error {
	return nil
}
