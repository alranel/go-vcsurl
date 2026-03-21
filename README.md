# go-vcsurl

Go library for parsing and manipulating URLs of VCS services

## Overview [![GoDoc](https://godoc.org/github.com/alranel/go-vcsurl?status.svg)](https://godoc.org/github.com/alranel/go-vcsurl)

This package provides useful functions for parsing and manipulating URLs of VCS
services. It can be used to check whether a given URL points to an account, a
repository, a file or a raw file, and to convert file URLs to raw URLs and root
raw URLs.

Supported services:

* GitHub
* GitLab (including instances on custom domains, detected via `_gitlab_session` cookie)
* BitBucket
* Forgejo (including instances on custom domains, detected via `/swagger.v1.json`)
* Gitea (including instances on custom domains, detected via `/swagger.v1.json`)
* Generic HTTP git repositories

## Install

```bash
go get github.com/alranel/go-vcsurl/v2
```

## Author

Alessandro Ranellucci

## License

MIT
