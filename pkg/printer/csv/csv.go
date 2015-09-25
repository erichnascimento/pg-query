package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type columnFormater func(interface{}) string

type RowPrinter interface {
	SetColumns([]string)
	Print([]interface{})
}

type CSVPrinter struct {
	writer          *csv.Writer
	columnFormaters []columnFormater
}

func NewCSVPrinter(writer io.Writer) *CSVPrinter {
	printer := new(CSVPrinter)
	printer.writer = csv.NewWriter(writer)
	//printer.Table = tablewriter.NewWriter(writer)
	//printer.Table.SetColWidth(50)

	return printer
}

func (tp *CSVPrinter) SetColumns(columns []string) {
	//tp.Table.SetHeader(columns)
	return
}

func (tp *CSVPrinter) Print(values []interface{}) {
	if tp.columnFormaters == nil {
		tp.createFormaters(values)
	}

	strRow := make([]string, len(values))

	for i, v := range values {
		formater := tp.columnFormaters[i]
		strRow[i] = formater(v)
	}

	//tp.Append(strRow)
	//fmt.Printf("%s", strRow)
	tp.writer.Write(strRow)
}

func (tp *CSVPrinter) Close() error {
	tp.writer.Flush()
	return nil
}

func (tp *CSVPrinter) createFormaters(row []interface{}) []columnFormater {
	if tp.columnFormaters == nil {
		tp.columnFormaters = make([]columnFormater, len(row))
		for i, v := range row {
			var formater columnFormater
			switch v.(type) {
			case int, int64:
				formater = func(value interface{}) string {
					return strconv.FormatInt(value.(int64), 10)
				}
			default:
				formater = func(value interface{}) string {
					if value == nil {
						return ""
					}

					return fmt.Sprintf("%s", value)
				}
			}

			tp.columnFormaters[i] = formater
		}
	}

	return tp.columnFormaters
}
