package models

import (
	"time"
)

func SaveException(etype string, esource string, emsg error) error {
	e := Exception{}
	e.Created = time.Now()
	e.Type = etype
	e.Source = esource
	e.Message = emsg
	DB.Save(&e)
	return nil
}

type Exception struct {
	ID      int64     `storm:"id,increment"`
	Created time.Time `storm:"index"`
	Source  string    `storm:"index"`
	Type    string    `storm:"index"`
	Message error
}
