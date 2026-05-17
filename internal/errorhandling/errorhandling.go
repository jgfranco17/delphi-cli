package errorhandling

import (
	"io"

	"github.com/fatih/color"
)

type Handler struct {
	writer io.Writer
}

func New(w io.Writer) *Handler {
	return &Handler{writer: w}
}

func (h *Handler) RenderError(err error) {
	redFmt := color.New(color.FgRed, color.Bold)
	redFmt.Fprintf(h.writer, "[ERROR]: %v\n", err)
}
