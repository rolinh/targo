# targo - A Go package to create and extract tar archives

![License](https://img.shields.io/badge/license-BSD-blue.svg)
[![Build Status](https://travis-ci.org/Rolinh/targo.png?branch=master)](https://travis-ci.org/Rolinh/targo)
[![GoDoc](http://godoc.org/github.com/Rolinh/targo?status.svg)](http://godoc.org/github.com/Rolinh/targo)
[![GoWalker](http://img.shields.io/badge/doc-gowalker-blue.svg?style=flat)](https://gowalker.org/github.com/Rolinh/targo)

`targo` provides functions to create or extract tar archives. This package has
no dependencies and relies only on the Go standard library.

## Usage ([full documentation here](http://godoc.org/github.com/Rolinh/targo))

- `func Create(destPath, dirPath string) error`: create a tar archive from
  `dirPath` into `destPath`.
- `func CreateInPlace(dirPath string) error`: create a tar archive from
  `dirPath` "in-place", ie `dirPath` is removed once the archive has been
  created and a `dirPath.tar` file is created.
- `func Extract(destPath, archivePath string) error`: extract a tar archive from
  `archivePath` into `destPath`.
- `func ExtractInPlace(archivePath string) error`: extract a tar archive from
  `archivePath` "in-place", ie `archivePath` is removed after the archive has
  been extracted (note: it expects `archivePath` to have a file extension).

## Notes

- As pointed out in the documentation of `Create` and `CreateInPlace` (see [#1](https://github.com/Rolinh/targo/issues/1)), the use of [filepath.Dir](https://golang.org/pkg/path/filepath/#Dir) introduce different behavior depending on the way you define your path:
  - With __'/foo/bar'__, `filepath.Dir` will consider __'bar'__ as the last token and return __'/foo'__. This will include the '__bar__' directory to be part of the tar as the root directory.
  - With __'/foo/bar/'__, `filepath.Dir` will consider __'/.'__ as the last token and return __'/foo/bar'__. This will ignore the '__bar__' directory and there will be not root directory in the resulting tar.
