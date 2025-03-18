package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type StructuredPrinter struct {
	w io.Writer
}

func NewStructuredPrinter(w io.Writer) *StructuredPrinter {
	return &StructuredPrinter{
		w: w,
	}
}

func (s *StructuredPrinter) Println(a ...any) {
	fmt.Fprintln(s.w, a...)
}

func (s *StructuredPrinter) Printf(format string, a ...any) {
	fmt.Fprintf(s.w, format, a...)
}

type Tabular interface {
	Headers() string
	Rows() []string
}

func (s *StructuredPrinter) StructuredPrint(data Tabular) error {
	w := tabwriter.NewWriter(s.w, 0, 0, 4, ' ', 0)

	fmt.Fprintln(w, data.Headers())
	for _, row := range data.Rows() {
		fmt.Fprintln(w, row)
	}

	return w.Flush()
}
