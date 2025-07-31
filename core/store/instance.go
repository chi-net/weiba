package store

import (
	"gorm.io/gorm"
	"sync"
)

var (
	mu  sync.RWMutex
	sql *gorm.DB
)
