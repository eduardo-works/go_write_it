package project

import "errors"

var ErrNotFile = errors.New("not a regular file")
var ErrNotDir = errors.New("not a directory")
var ErrNotProjectRoot = errors.New("not a project root folder")
var ErrNotBelongToProject = errors.New("does not belong to project")
var ErrNotManuscriptFile = errors.New("not a manuscript file")
var ErrInvalidValue = errors.New("invalid value")
var ErrEmptyValue = errors.New("empty value")
