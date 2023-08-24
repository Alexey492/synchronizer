package logic

import "time"

type Operation int

const (
	Create Operation = iota
	Update
	Delete
)

type FileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
	IsDir   bool
}
