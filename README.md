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
