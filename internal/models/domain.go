package models

import "gorm.io/gorm"

type Domain struct {
	gorm.Model
	FullDomain string `gorm:"unique"`
	Records    []Record
}
