# fiximports

`fiximports` formats and adjusts imports for go source files. It improves on `goimports`
by auto-detecting and grouping local go module imports and by fixing disjointed groups.

## Usage

```
go install github.com/corverroos/fiximports
fiximports -verbose path/to/foo.go path/with/bar.go 
```