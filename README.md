# immuparser

Compile with:
```bash
go build main.go
```
## Usage
```bash
./main --help
Parse a file with following format:
2020-05-21 20:02:00 73.115.58.89:50800 Wine

Usage:
  immuparser [flags]

Examples:

main --ledger myledger
main --ledger myledger --flushSize 1000
main --ledger myledger --flushSize 1000 --source /tmp/log.txt


Flags:
      --config string   config file (default is $HOME/.immuparser.yaml)
  -f, --flushSize int   flush size (default 1000)
  -h, --help            help for immuparser
  -l, --ledger string   ledger name (default "default")
  -s, --source string   source file path (default "/tmp/log.txt")

```


Launch command with
```bash
./main --ledger {ledgername} --flushSize {flush size} --source {file path}
```
Example:
```bash
./main --ledger myledger444 --flushSize 1000 --source ./log.txt
```
