# ggrep
a command-line tool for text search, sepcially for logid searching. you can use logid group the muliti-lines merge to a one record.
# demo
merge multi-lines to a line for logid:
```shell
ggrep --group 'logid[\w+]' 'LOG1' 'LOG2' --buff-line 200 --merge-lines
```

log-content:
```text
20240119 LogId[111] LOG1
20240119 LogId[111] LOG2
20240119 LogId[111] LOG3
20240119 LogId[111] LOG4
20240119 LogId[111] LOG5
50240119 LogId[222] LOG1
20240119 LogId[222] LOG2
20240119 LogId[222] LOG3
20240119 LogId[222] LOG4
20240119 LogId[222] LOG5
```
output
```text
20240119 LogId[111] LOG220240119 LogId[111] LOG5
20240119 LogId[222] LOG220240119 LogId[222] LOG5
```

seek multi condition for log scanning
```shell
ggrep --group 'logid[\w+]' 'LOG1' 'LOG2' --buff-line 200
```

log-content:
```text
20240119 LogId[111] LOG1
20240119 LogId[111] LOG2
20240119 LogId[111] LOG3
20240119 LogId[111] LOG4
20240119 LogId[111] LOG5
50240119 LogId[222] LOG1
20240119 LogId[222] LOG2
20240119 LogId[222] LOG3
20240119 LogId[222] LOG4
20240119 LogId[222] LOG5
```
output
```text
20240119 LogId[111] LOG2
20240119 LogId[111] LOG5
20240119 LogId[222] LOG2
20240119 LogId[222] LOG5
```

# install

```shell
go install github.com/lingdor/ggrep@latest
```
or download your release versions: [Releases](https://github.com/lingdor/ggrep/releases).