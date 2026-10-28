package resolver

import "errors"

var ErrNotFile = errors.New("path is not a regular file!")

var ErrNotDir = errors.New("path is not a directory!")

var ErrNotProjectRoot = errors.New("path is not a project root folder!")
