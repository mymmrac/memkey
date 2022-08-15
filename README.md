# :package: MemKey

[![Go Reference](https://pkg.go.dev/badge/github.com/mymmrac/memkey#section-readme.svg)](https://pkg.go.dev/github.com/mymmrac/memkey)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mymmrac/memkey?logo=go)](go.mod)

[![CI Status](https://github.com/mymmrac/memkey/actions/workflows/ci.yml/badge.svg)](https://github.com/mymmrac/memkey/actions/workflows/ci.yml)
[![Race Testing](https://github.com/mymmrac/memkey/actions/workflows/race-tests.yml/badge.svg)](https://github.com/mymmrac/memkey/actions/workflows/race-tests.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_memkey&metric=alert_status)](https://sonarcloud.io/dashboard?id=mymmrac_memkey)
[![Go Report](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/mymmrac/memkey)
<br>
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_memkey&metric=coverage)](https://sonarcloud.io/dashboard?id=mymmrac_memkey)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_memkey&metric=code_smells)](https://sonarcloud.io/dashboard?id=mymmrac_memkey)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_memkey&metric=ncloc)](https://sonarcloud.io/dashboard?id=mymmrac_memkey)

Very simple type-safe, thread-safe in memory key-value store with zero dependencies.

## :zap: Getting Started

How to get the library:

```shell
go get -u github.com/mymmrac/memkey
```

Create new store:

```go
s := &memkey.Store[comparable]{}
```

> As a key for store you can use any `comparable` type, values can be of any type.

## :kite: Features

All functions are type-safe and thread-safe, there is also type-unsafe variant of all functions (just use same methods
of [`Store`](https://pkg.go.dev/github.com/mymmrac/memkey#Store) struct).

| Method                                                            | Description                             |
|-------------------------------------------------------------------|-----------------------------------------|
| [`Get`](https://pkg.go.dev/github.com/mymmrac/memkey#Get)         | Get value                               |
| [`Set`](https://pkg.go.dev/github.com/mymmrac/memkey#Set)         | Set value                               |
| [`Has`](https://pkg.go.dev/github.com/mymmrac/memkey#Has)         | Check if value exists                   |
| [`Delete`](https://pkg.go.dev/github.com/mymmrac/memkey#Delete)   | Delete value and return true if deleted |
| [`Len`](https://pkg.go.dev/github.com/mymmrac/memkey#Len)         | Number of elements stored               |
| [`Keys`](https://pkg.go.dev/github.com/mymmrac/memkey#Keys)       | Get keys                                |
| [`Values`](https://pkg.go.dev/github.com/mymmrac/memkey#Values)   | Get values                              |
| [`Entries`](https://pkg.go.dev/github.com/mymmrac/memkey#Entries) | Get all key-value pairs                 |
| [`ForEach`](https://pkg.go.dev/github.com/mymmrac/memkey#ForEach) | Iterate over key-value pairs            |

## :jigsaw: Usage

```go
s := &memkey.Store[int]{}

Set(s, 1, "mem")
Set(s, 2, "key")
Set(s, 3, 42.0)

m, ok := Get[string](s, 1)
k, ok := Get[string](s, 2)
// Here `m` & `k` will be of type string and have zero value if not found, 
// `ok` will indicate if value was found

n, ok := Get[float64](s, 3)
// Here `n` will be float64

found := Has[uint](s, 2)
// Here `found` will be `false`, since value with key `2` is `string` and not an `uint`

keys := KeysRaw(s)
// Here `keys` will be a slice of `1`, `2` and `3` in non deterministic order, 
// `..Raw` methods will ignore type of values
```

## :closed_lock_with_key: License

MemKey is distributed under [MIT](LICENSE).
