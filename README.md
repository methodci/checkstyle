# checkstyle / chksutil

[![CI](https://github.com/methodci/checkstyle/actions/workflows/ci.yml/badge.svg)](https://github.com/methodci/checkstyle/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/methodci/checkstyle?status.svg)](https://godoc.org/github.com/methodci/checkstyle)
[![Go Report Card](https://goreportcard.com/badge/github.com/methodci/checkstyle)](https://goreportcard.com/report/github.com/methodci/checkstyle)

`checkstyle` is a library for working with checkstyle files.

The included `chksutil` is a utility for inspecting and diffing checkstyle files.

```console
$ chksutil diff old.xml new.xml
Fixed info on GoalSetting.php:44 - DocblockTypeContradiction: Cannot resolve types for $value - docblock-defined type int does not contain null
Fixed info on GoalSetting.php:44 - RedundantConditionGivenDocblockType: Found a redundant condition when evaluating docblock-defined type $value and trying to reconcile type 'int' to !null
Created info on GoalSetting.php:72 - MissingParamType: Parameter $value has no provided 
```

## Install chksutil

binaries are attached to the releases

https://github.com/methodci/checkstyle/releases

otherwise you can install from source via

```console
$ go install github.com/methodci/checkstyle/cmd/chksutil@latest
```
