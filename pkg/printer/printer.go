package printer

type RowPrinter interface {
	SetColumns([]string)
	Print([]interface{})
	Close() error
}
