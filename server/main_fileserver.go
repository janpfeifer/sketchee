package main

// http.FileSystem implementation to serve the static files using HTTP. Copied from
// https://golang.org/pkg/net/http/#FileServer
// It's only needed because the default `http.Dir` doesn't seem to work in WSL.
// More in https://github.com/golang/go/issues/27128

import (
	"github.com/golang/glog"
	"net/http"
	"os"
	"strings"
)

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []os.FileInfo, err error) {
	glog.V(1).Infof("Readdir(n=%d)", n)
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fs dotFileHidingFileSystem) Open(name string) (http.File, error) {
	glog.V(1).Infof("Open(name=%q)", name)
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}
