package db

import (
	"os"
	"strings"
)

type log struct {
	dir           string
	activeFile    *os.File
	currentID     uint32
	currentOffset uint32
}
