# pg-query

Use this tool for run queries on multiple Postgres hosts and databases pre-configured.

## Instalation

```
$ go get github.com/erichnascimento/pg-query
```

## Usage
```
$ pg-query -h

Usage:
  pg-query [options] --config <config> <sql>
  pg-query [options] --config <config> --input-file <sql-file>
  pg-query -h | -- help
  pg-query -v | --version

Options:
  -c, --config ./config.yml                configuration file path
  -d, --databases db1,db2,dbN              databases filter
  -f, --input-file                         SQL file to execute. For reading from stdin, use "-"
  -H, --hosts host1,host2,hostN            hosts filter
  -F, --format (table | csv)               output format: table | csv
  -h, --help                               output help information
  -v, --version                            output version

```

## Examples

### Config file example
```yaml
hosts:
    - # My dev Postgres
        name: dev.local
        user: postgres
        password: postgres
        host: localhost
        port: 5432
        databases: [myapp]

    - # My Production Postgres
        name: production
        user: postgres
        password: postgres
        host: myapp.io
        port: 5432
        databases: [myapp1,myapp2]
```

### Simple way
Execute a query over all configured databases in ./config.yml
```
$ pg-query -c config.yml "SELECT NOW()"
2015/06/18 15:53:32 Connecting to "myapp" (localhost:5432)
+--------------------------------------+
|                 NOW                  |
+--------------------------------------+
| 2015-06-18 18:53:32.198413 +0000 UTC |
+--------------------------------------+

2015/06/18 15:53:32 Connecting to "myapp1" (myapp.io:5432)
+--------------------------------------+
|                 NOW                  |
+--------------------------------------+
| 2015-06-18 18:53:32.219316 +0000 UTC |
+--------------------------------------+

2015/06/18 15:53:32 Connecting to "myapp2" (myapp.io:5432)
+--------------------------------------+
|                 NOW                  |
+--------------------------------------+
| 2015-06-18 18:53:32.220903 +0000 UTC |
+--------------------------------------+
```

### Filter by host
Execute a query over all configured databases in ./config.yml that host match with `localhost` or `server`
```
$ pg-query -c config.yml -H localhost,server "SELECT NOW()"
```


### Filter by database name

Execute a query over all configured databases in ./config.yml that database name match with `myapp1` or `myapp2`

```
$ pg-query -c config.yml -d myapp1,myapp2 "SELECT NOW()"
```

### Filter by host and database name

Execute a query over all configured databases in ./config.yml that database name match with `myapp` and host match witch `localhost`

```
$ pg-query -c config.yml -d myapp -H localhost "SELECT NOW()"
```

### Output format

Show data as table (default)

```
$ pg-query -c config.yml -F table -H localhost "SELECT NOW()"
```

Show data as csv

```
$ pg-query -c config.yml -F csv -H localhost "SELECT NOW()"
```


# License

MIT
