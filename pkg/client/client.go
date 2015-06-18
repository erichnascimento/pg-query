package	client

import (
	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"fmt"
	"database/sql"
	"log"

	"github.com/erichnascimento/pg-query/pkg/printer"
)

type Client struct {
	*Config
	*sql.DB
}

func New(config *Config) (*Client, error) {
	c := &Client{
		Config: config,
	}

	if true {
		log.Printf(`Connecting to "%s" (%s:%d)`, config.Database, config.Host, config.Port)
	}
	
	err := c.connect()
	return c, err
}

func (c *Client) connect() error {

	connectionStr := fmt.Sprintf(`host=%s port=%d user=%s dbname=%s sslmode=disable`, c.hostIP, c.Port, c.User, c.Database)

	if c.Password != "" {
		connectionStr += fmt.Sprintf(` password=%s`, c.Password)
	}

	db, err := sql.Open("postgres", connectionStr)
	if err != nil  {
		return err
	}
	c.DB = db

	return err
}

func (c *Client) Close() error {
	return c.DB.Close()
}

func (c *Client) Query(sql string, printer printer.RowPrinter) error {
	rows, err := c.DB.Query(sql)
	if err != nil {
		return err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	numColumns := len(columns)
	printer.SetColumns(columns)

	values := make([]interface{}, numColumns)
	pValues := rowToPointer(values)

	for rows.Next() {
		
		err := rows.Scan(pValues...)
		if err != nil {
			return err
		}

		printer.Print(values)
	}

	return nil
}

func rowToPointer(row []interface{}) []interface{} {
	rowPointer := make([]interface{}, len(row))
	for i, _ := range row {
		rowPointer[i] = &row[i]
	}

	return rowPointer
}