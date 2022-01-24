# fiximports

`fiximports` formats and adjusts imports for go source files. It improves on `goimports`
by auto-detecting and grouping local go module imports and by fixing disjointed groups.

## Usage

```
go install github.com/corverroos/fiximports
fiximports -verbose path/to/foo.go path/with/bar.go 
```

## pre-commit

Usage as part of [pre-commit](https://pre-commit.com/) githook framework is supported via:
```
- repo: https://github.com/corverroos/fiximports
  hooks:
  - id: fiximports
```