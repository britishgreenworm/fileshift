package models

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Name    string
	UserId  int64
	Status  string
	Size    int64
	Changed bool
	OutBox  bool
}
