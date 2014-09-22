package db

import (
	"strings"

	"github.com/go-distributed/xtree/third-party/github.com/google/btree"
)

// Path is a collection of infomation on specific path.
type Path struct {
	p     string
	v     Value
	level int
}

func newPath(pathname string) Path {
	pathname = strings.TrimRight(pathname, "/")
	res := Path{
		p: pathname,
	}

	names := strings.Split(
		strings.TrimLeft(pathname, "/"),
		"/")

	res.level = len(names)
	for _, name := range names {
		if name == "" {
			res.level--
		}
	}

	return res
}

func newPathForLs(pathname string) Path {
	res := newPath(pathname)
	res.level += 1
	res.p += "/"
	return res
}

func (a Path) Less(treeItem btree.Item) bool {
	b := treeItem.(Path)

	if a.level != b.level {
		return a.level < b.level
	}

	return a.p < b.p
}
