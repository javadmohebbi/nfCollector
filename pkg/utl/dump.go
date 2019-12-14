package utl

import (
	"github.com/landoop/tableprinter"
	"os"
)

// Dump - Print it as table
func Dump(metrics []Metric) {
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.Print(metrics)
}

