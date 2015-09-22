package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/tj/docopt"

	"github.com/erichnascimento/pg-query/pkg/client"
	"github.com/erichnascimento/pg-query/pkg/config"
	"github.com/erichnascimento/pg-query/pkg/printer"
	"github.com/erichnascimento/pg-query/pkg/printer/csv"
	"github.com/erichnascimento/pg-query/pkg/printer/table"
)

// Version of program
const VERSION = "0.1.0"

// Usage description
const Usage = `
  Usage:
    pg-query [options] --config <config> <sql>
    pg-query [options] --config <config> --input-file <sql-file>
    pg-query -h | -- help
    pg-query -v | --version

  Options:
    -c, --config ./config.yml                configuration file path
    -d, --databases db1,db2,dbN              databases filter
    -f, --input-file                         SQL file to execute
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
	if err != nil {
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
			h.ApplyDatabaseFilter(dbsArr)
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

	sql := ""

	// SQL String
	if val := args["<sql>"]; val != nil {
		sql = val.(string)
	}

	// SQL from file
	if val := args["<sql-file>"]; val != nil {
		filename := val.(string)
		fileContent, err := ioutil.ReadFile(filename)

		if err != nil {
			log.Fatalf("Error on open sql file: %s", err)
		}

		sql = string(fileContent)
	}

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
