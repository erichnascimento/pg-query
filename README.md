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
    pg-query --config <config> [--hosts <hosts>] [--databases <db>] [--format <fmt>] <sql>
    pg-query -h | -- help
    pg-query -v | --version

  Options:
    -c, --config config    configuration file path
    -d, --databases db     databases to run
    -H, --hosts hosts      hosts to run
    -F, --format fmt       output format: table | csv
    -h, --help             output help information
    -v, --version          output version

```

## Examples

### Run a query over all configured databases in ./examples/config.yml
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

# License

MIT