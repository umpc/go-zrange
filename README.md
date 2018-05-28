# go-zrange

[![GoDoc](https://godoc.org/github.com/umpc/go-zrange?status.svg)](https://godoc.org/github.com/umpc/go-zrange)
[![Go Report Card](https://goreportcard.com/badge/github.com/umpc/go-zrange)](https://goreportcard.com/report/github.com/umpc/go-zrange)

```sh
go get -u github.com/umpc/go-zrange
```

This package implements an efficient algorithm for performing spatial range queries with Geohash-encoded keys and a search radius.

The `RadialRange` method appears to be sufficient for range queries of around 5,000km or less. Changes that efficiently add support for larger query ranges are welcome.

## Example usage

```go
...

rangeParams := zrange.RadialRangeParams{
  Radius:    32.18688,
  Latitude:  37.334722,
  Longitude: -122.008889,
}

keyRanges := rangeParams.RadialRange()

...
```

**Note:** Geohash range searches are imprecise. Results should be filtered using the input radius where precision is desired.

## References

* The `RadialRange` method was inspired by the algorithm in the "Search" section of [this page](https://web.archive.org/web/20180526044934/https://github.com/yinqiwen/ardb/wiki/Spatial-Index#search).
