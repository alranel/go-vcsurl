# go-vcsurl

Go library for parsing and manipulating URLs of VCS services

## Overview [![GoDoc](https://godoc.org/github.com/alranel/go-vcsurl?status.svg)](https://godoc.org/github.com/alranel/go-vcsurl)

This package provides useful functions for parsing and manipulating URLs of VCS services such as GitHub, GitLab, BitBucket. It can be used to check whether a given URL points to an account, a repository, a file or a raw file. It also provides functions for converting file URLs to raw URLs and root raw URLs. Note that this library only performs syntactic checks and string manipulation, and it does not perform network calls to check that the given resources actually exist.

## Install

```bash
go get github.com/alranel/go-vcsurl
```

## Author

Alessandro Ranellucci

## License

MIT
