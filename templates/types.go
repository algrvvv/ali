package templates

import "errors"

var (
	ErrFailedToInitNewTempl = errors.New("failed to init new templ")
	ErrTemplAlreadyExists   = errors.New("template already exists\nuse --force for force rewrite")
)
