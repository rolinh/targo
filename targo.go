// Copyright 2015 Robin Hahling. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package targo provides functions to create or extract tar archives.
package targo

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Create creates a tar archive from the directory specified by dirPath.
// The resulting tar archive format is in POSIX.1 format.
// Calling this function with dirPath having a trailing slash results in the
// content of dirPath to be used at the root level of the archive rather than
// the directory pointed out by dirPath itself.
func Create(destPath, dirPath string) error {
	fi, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("given path is not a directory: " + dirPath)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tw := tar.NewWriter(file)
	defer tw.Close()

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		var link string
		mode := info.Mode()
		switch {
		// symlinks need special treatment
		case mode&os.ModeSymlink != 0:
			linkDest, _ := os.Readlink(path)
			if link, err = filepath.EvalSymlinks(path); err == nil {
				if rel, err := filepath.Rel(filepath.Dir(path), link); err == nil {
					link = rel
				} else {
					link = linkDest
				}
			} else {
				// it may be a broken symlink, simply attempt to read it
				link = linkDest
			}

		// we don't want to tar these sort of files
		case mode&(os.ModeNamedPipe|os.ModeSocket|os.ModeDevice) != 0:
			return nil
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		// Name is usually only the basename when created with FileInfoHeader()
		hdr.Name, err = filepath.Rel(filepath.Dir(dirPath), path)
		if err != nil {
			return err
		}

		// Some pre-POSIX.1-1988 tar implementations indicated a directory by
		// having a trailing slash in the name. Honor that here.
		if info.IsDir() {
			hdr.Name += "/"
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// no content to write if it is a directory or symlink
		if !info.Mode().IsRegular() {
			return nil
		}

		return func() error {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err = io.Copy(tw, f); err != nil {
				return err
			}
			return nil
		}()
	})

	return err
}

// CreateInPlace behaves just as Create but it creates archive in place. This
// means that the original directory specified by dirPath is removed after the
// tar archive is created. The .tar suffix is automatically added to dirPath
// and is used as the name of the newly created archive.
func CreateInPlace(dirPath string) error {
	if err := Create(dirPath+".tar", dirPath); err != nil {
		return err
	}
	return os.RemoveAll(dirPath)
}

// Extract extracts a tar archive given its path.
func Extract(destPath, archivePath string) error {
	fi, err := os.Stat(archivePath)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return errors.New("given path is a directory: " + archivePath)
	}

	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return err
	}

	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	tr := tar.NewReader(archiveFile)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		mode := hdr.FileInfo().Mode()
		switch {
		case mode&os.ModeDir != 0:
			if err := os.Mkdir(filepath.Join(destPath, hdr.Name), mode); err != nil {
				return err
			}
		case mode&os.ModeSymlink != 0:
			os.Symlink(hdr.Linkname, filepath.Join(destPath, hdr.Name))
		default: // consider it a regular file
			createFile := func() error {
				f, err := os.Create(filepath.Join(destPath, hdr.Name))
				if err != nil {
					return err
				}
				defer f.Close()

				if _, err := io.Copy(f, tr); err != nil {
					return err
				}
				return nil
			}

			if err = createFile(); err != nil {
				return err
			}
		}
	}

	return nil
}

// ExtractInPlace extracts a tar archive, in place, given its path. The
// original tar archive is removed after extraction and only its content
// remains.
// Note that archivePath is expected to contain a file extension.
func ExtractInPlace(archivePath string) error {
	ext := filepath.Ext(archivePath)
	if ext == "" {
		return errors.New("expected a file extension (" + archivePath + ")")
	}
	destPath := filepath.Dir(strings.TrimSuffix(archivePath, ext))
	if err := Extract(destPath, archivePath); err != nil {
		return err
	}
	return os.Remove(archivePath)
}
