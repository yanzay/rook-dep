# rook

Dependency management tool for Go

## Install

`rook` has no dependencies, so install it with simple
```bash
go get github.com/yanzay/rook
```

## Usage
Prepare the `Rookfile` with required packages and their versions:

Example `Rookfile`:
```
google/go-github/github: latest
golang/oauth2: latest
hashicorp/vault/api: v0.3.0
gorilla/mux: latest
```

Note that package `golang/oauth2` will be internally resolved to path `golang.org/x/oauth2`.

Than simply run:
```
rook
```

It will install all your dependencies. That's it!

`rook` saves current dependency versions in file `Rookfile.lock`. If you want your build to be reproducable, commit `Rookfile.lock` and `rook` will install dependencies from this file in future. If you want to upgrade packages labeled as `latest`, just delete `Rookfile.lock` and launch `rook` again.

## Search

To see which versions of package are available now, run:
```
rook search owner/package
```

Example:
```
$ rook search hashicorp/vault/api
latest
v0.3.1
v0.3.0
v0.3.0-rc
v0.2.0
v0.2.0.rc1
v0.1.2
v0.1.1
v0.1.0
---
Rook success!
```

