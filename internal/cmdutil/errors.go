package cmdutil

import (
	"fmt"
	"io"
	"os"
)

type CmdErrorHandler struct {
	w io.Writer
}

func NewCmdErrorHandler(w io.Writer) *CmdErrorHandler {
	return &CmdErrorHandler{w: w}
}

func (h *CmdErrorHandler) HandleAndExit(err error, code int) {
	if err == nil {
		return
	}

	fmt.Fprintf(h.w, "error: %v\n", err)
	os.Exit(code)
}
