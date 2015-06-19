package main

import (
	"log"
	"strings"
	"os"

	"github.com/tj/docopt"

	"github.com/erichnascimento/pg-query/pkg/config"
	"github.com/erichnascimento/pg-query/pkg/client"
	"github.com/erichnascimento/pg-query/pkg/printer/csv"
	"github.com/erichnascimento/pg-query/pkg/printer/table"
	"github.com/erichnascimento/pg-query/pkg/printer"
)

const VERSION = "0.0.2"

const Usage = `
  Usage:
    pg-query --config <config> [--hosts <hosts>] [--databases <db>] [--format <fmt>] <sql>
    pg-query -h | -- help
    pg-query -v | --version

  Options:
    -c, --config ./config.yml                configuration file path
    -d, --databases db1,db2,dbN              databases filter
    -H, --hosts host1,host2,hostN            hosts filter
    -F, --format (table | csv)               output format: table | csv
    -h, --help                               output help information
    -v, --version                            output version
`

func explodeStr(str string, trim bool) (parts []string) {
	parts = strings.Split(str, ",")

	if !trim {
		return
	}

	for i, s := range parts {
		parts[i] = strings.TrimSpace(s)
	}

	return
}

func runSQLQuery(sql string, conf *client.Config, rowPrinter printer.RowPrinter) error {
	// Validate connection config
	if err := conf.Validate(); err != nil {
		log.Fatalf(`configuration error: "%s"`, err)
	}

	// Open connection
	db, err := client.New(conf)
	if err != nil {
		log.Fatalf(`connection error: "%s"`, err)
	}
	defer db.Close()

	// Execute SQL
	err = db.Query(sql, rowPrinter)
	if err != nil  {
		return err
	}

	rowPrinter.Close()
	return nil
}

func main() {
	// Parse args
	args, err := docopt.Parse(Usage, nil, true, VERSION, false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	// Load config from file
	configFile := args["--config"].(string)
	err, conf := config.CreateFromFile(configFile)
	if err != nil {
		log.Fatalf("Error on load config from file: %s", err)
	}

	// Filter hosts
	if hosts := args["--hosts"]; hosts != nil {
		if err := conf.ApplyHostFilter(explodeStr(hosts.(string), true)); err != nil {
			log.Fatalf("Error on filter hosts: %s", err)
		}
	}

	// Filter databases for each host
	if dbs := args["--databases"]; dbs != nil {
		dbsArr := explodeStr(dbs.(string), true)
		for _, h := range conf.Hosts {
			if err := h.ApplyDatabaseFilter(dbsArr); err != nil {
				log.Fatalf("Error on filter dbs: %s", err)
			}
		}
	}

	// Output format
	format := "table"
	if fmt := args["--format"]; fmt != nil {
		format = fmt.(string)
	}

	var rowPrinter printer.RowPrinter
	switch format {
		case "csv":
			rowPrinter = csv.NewCSVPrinter(os.Stdout)
		default:
			rowPrinter = table.NewTablePrinter(os.Stdout)
	}

	sql := args["<sql>"].(string)

	// Run SQL
	for _, h := range conf.Hosts {
		for _, db := range h.Databases {
			// Create a connection config
			conf := client.NewConfig(h.Host, h.Port, h.User, h.Password, db)
			err := runSQLQuery(sql, conf, rowPrinter)
			if err != nil {
				log.Fatalf("Error on execute query: %s", err)
			}
		}
	}

	//log.Printf("Bye :)")
}

