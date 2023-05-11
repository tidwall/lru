**This project has been discontinued. Please use [tidwall/tinylru](https://github.com/tidwall/tinylru) instead.**

# lru

[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tidwall/lru)

A simple and efficient lru cache package for Go.

Under the hood it uses a hashmap combined with a doubly linked list.

# Getting Started

## Installing

To start using lru, install Go and run `go get`:

```sh
$ go get -u github.com/tidwall/lru
```

This will retrieve the library.

## Usage

The `Cache` type works similar to a standard Go map, and includes four methods:
`Set`, `Get`, `Delete`, `Len`.

```go
package main

import "github.com/tidwall/lru"

func main() {
    // Create a cache that holds no more than three entries.
    cache := lru.New(3, nil)
    
    // Add four entries, which is one more than what the cache allows.
    cache.Set("43084", `{"id":43084,"name":"Anne Bancroft"}`)
    cache.Set("99103", `{"id":99103,"name":"Duane Jones"}`)
    cache.Set("19520", `{"id":19520,"name":"Bette Davis"}`)
    cache.Set("67531", `{"id":67531,"name":"Joan Crawford"}`)


    fmt.Println(cache.Get("43084")) // got evicted when 67531 was added
    fmt.Println(cache.Get("99103"))
    fmt.Println(cache.Get("19520"))
    fmt.Println(cache.Get("67531"))

    // output: 
    // <nil>
    // {"id":99103,"name":"Duane Jones"}
    // {"id":19520,"name":"Bette Davis"}
    // {"id":67531,"name":"Joan Crawford"}
}
```

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

`lru` source code is available under the MIT [License](/LICENSE).
