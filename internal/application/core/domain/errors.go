package domain

import (
	"fmt"
)

var (
	ErrKeyNotFound = fmt.Errorf("key not found")
	ErrKeyExists   = fmt.Errorf("key already exists")
)
