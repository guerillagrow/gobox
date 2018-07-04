package models

import (
	"time"
)

type Plant struct {
	ID          int64  `storm:"id,increment"`
	Parent      int64  `storm:"index"`
	Name        string `storm:"index"`
	Description string
	Specs       string
	Notes       string
	Yields      []PlantYield
}

func (self *Plant) Delete() {
	DB.DeleteStruct(self)
}

func (self *Plant) Save() {
	DB.Save(self)
}

type PlantYield struct {
	ID      int64     `storm:"id,increment"`
	PlantID int64     `storm:"index"`
	Created time.Time `storm:"index"`
	Yield   float64
}

func (self *PlantYield) Delete() {
	DB.DeleteStruct(self)
}

func (self *PlantYield) Save() {
	DB.Save(self)
}
