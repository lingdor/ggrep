# ggrep
a command-line tool for text search, sepcially for logid searching. you can use logid group the muliti-lines merge to a one record.
# demo
seek multi condition for log scanning:

```shell
ggrep --group 'LogId\[[:logid:]\]' --grep 'LOG1' --grep 'LOG2' --orderly-match  demo/demo.log
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
cat demo/demo.log|ggrep  --group 'LogId\[[:logid:]\]' --grep 'LOG1' --grep 'LOG2' --merge-lines  --color=always
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
test group match

```shell
cat demo/demo.log|ggrep  --group 'LogId\[[:logid:]\]' --color=always|head
```


# install
```shell

/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/lingdor/ggrep/HEAD/install.sh)"

```

or
```shell
go install github.com/lingdor/ggrep@latest
```
or download your release versions: [Releases](https://github.com/lingdor/ggrep/releases).