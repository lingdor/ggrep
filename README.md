# ggrep
a command-line tool for text search, sepcially for logid searching. you can use logid group the muliti-lines merge to a one record.
# demo
seek multi condition for log scanning:

```shell
ggrep --group 'LogId\[[:logid:]\]' --grep 'LOG1' --grep 'LOG2' --match-orderly  demo/demo.log
```

log-content:
```text
20240119 LogId[111] LOG1
20240119 LogId[111] LOG2
20240119 LogId[111] LOG3
50240119 LogId[222] LOG1

```
output
```text
20240119 LogId[111] LOG120240119 LogId[111] LOG2
50240119 LogId[222] LOG1
```

merge multi-lines to a line for logid
```shell
cat demo/ ggrep --group 'logId\[\w+\]' --grep 'log1' --grep 'log2' --merge-lines demo/demo.log
```

log-content:
```text
20240119 LogId[111] LOG1
20240119 LogId[111] LOG2
20240119 LogId[111] LOG3
50240119 LogId[222] LOG1

```
output
```text
20240119 LogId[111] LOG120240119 LogId[111] LOG2
50240119 LogId[222] LOG1
```

# install
```shell

curl /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

```

or
```shell
go install github.com/lingdor/ggrep@latest
```
or download your release versions: [Releases](https://github.com/lingdor/ggrep/releases).