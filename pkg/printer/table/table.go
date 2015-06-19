package table

import(
	"fmt"
	"io"
	"strconv"
	"github.com/olekukonko/tablewriter"
)

type columnFormater func (interface{}) string

type TablePrinter struct {
	writer io.Writer
	*tablewriter.Table
	columnFormaters []columnFormater
}

func NewTablePrinter(writer io.Writer) *TablePrinter {
	printer := new(TablePrinter)
	printer.writer = writer
	printer.initialize()
	return printer
}

func (tp *TablePrinter) initialize() error {
	tp.Table = tablewriter.NewWriter(tp.writer)
	tp.Table.SetColWidth(50)

	return nil
}

func (tp *TablePrinter) SetColumns(columns []string) {
	tp.Table.SetHeader(columns)
}

func (tp *TablePrinter) Print(values []interface{}) {
	if tp.columnFormaters == nil {
		tp.createFormaters(values)
	}

	strRow := make([]string, len(values))

	for i, v := range values {
		formater := tp.columnFormaters[i]
		strRow[i] = formater(v)
	}
	
	tp.Append(strRow)
}

func (tp *TablePrinter) Close() error {
	tp.Render()
	tp.initialize()
	return nil
}

func (tp *TablePrinter) createFormaters(row []interface{}) []columnFormater {
	if tp.columnFormaters == nil {
		tp.columnFormaters = make([]columnFormater, len(row))
		for i, v := range row {
			var formater columnFormater
			switch v.(type) {
			case int, int64:
				formater = func (value interface{}) string {
					return strconv.FormatInt(value.(int64), 10)
				}
			default:
				formater = func (value interface{}) string {
					return fmt.Sprintf("%s",value)
				}
			}

			tp.columnFormaters[i] = formater
		}
	}

	return tp.columnFormaters
}
