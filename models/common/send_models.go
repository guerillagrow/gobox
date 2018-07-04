package common

import (
	"time"
)

type Response struct {
	Sensor  string    `json:"sensor"`
	Type    string    `json:"type"`
	Created time.Time `json:"created"`
	Value   float64   `json:"value"`
}
