package main

import (
	"fmt"
	"github.com/yuuki0xff/clustertest/models"
	"io"
)

type resultRender interface {
	Render(w io.Writer, d models.TaskDetail)
}

type singleResultRender struct{}

func (render singleResultRender) Render(w io.Writer, d models.TaskDetail) {
	fmt.Fprintf(w, "Status: %s\n", d.State())
	tr := d.Result()
	if tr == nil {
		// Result is not available.
		return
	}

	if err := tr.Error(); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
	}
	if r := tr.BeforeResult(); r != nil {
		render.renderHeader(w, "Before")
		render.renderResult(w, r)
	}
	if r := tr.ScriptResult(); r != nil {
		render.renderHeader(w, "Main")
		render.renderResult(w, r)
	}
	if r := tr.AfterResult(); r != nil {
		render.renderHeader(w, "After")
		render.renderResult(w, r)
	}
}
func (singleResultRender) renderHeader(w io.Writer, name string) {
	fmt.Fprintf(w, "-------------------- %s --------------------\n", name)
}
func (singleResultRender) renderResult(w io.Writer, r models.ScriptResult) {
	fmt.Fprintf(w, "ExitCode: %d\n", r.ExitCode())
	if host := r.Host(); host != "" {
		fmt.Fprintf(w, "Host: %s\n", host)
	}
	if start := r.StartTime(); !start.IsZero() {
		fmt.Fprintf(w, "Start: %s\n", start.String())
		if end := r.EndTime(); !end.IsZero() {
			fmt.Fprintf(w, "End: %s\n", end.String())
		}
	}
	io.WriteString(w, "Output:\n")
	w.Write(r.Output())
}

type multipleResultRender struct {
	notFirst bool
}

func (render *multipleResultRender) Render(w io.Writer, d models.TaskDetail) {
	if !render.notFirst {
		render.notFirst = true
	} else {
		// Write gap between results.
		io.WriteString(w, "\n\n")
	}

	fmt.Fprintf(w, "==================>> TaskID: %s <<==================\n", d.TaskID().String())
	r := singleResultRender{}
	r.Render(w, d)
}
