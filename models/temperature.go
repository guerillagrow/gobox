package models

import (
	//"fmt"
	//"errors"
	//"log"
	"sort"
	"time"

	//"github.com/guerillagrow/tconv"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type ISensorData interface {
	GetCreated() time.Time
	GetSensor() string
	GetValue() float64
}

type Temperature struct {
	ID      int64     `storm:"id,increment"`
	Sensor  string    `storm:"index"`
	Created time.Time `storm:"index"`
	Value   float64
}

func (self Temperature) GetNode() storm.Node {
	return DB.From("sensd", "temperature", self.Sensor)
}
func (self Temperature) GetValue() float64 {
	return self.Value
}

func (self Temperature) GetSensor() string {
	return self.Sensor
}

func (self Temperature) GetCreated() time.Time {
	return self.Created
}

func (self *Temperature) Delete() error {
	if self.Sensor == "" {
		return nil
	}
	node := self.GetNode()
	return node.DeleteStruct(self)
	//DB.DeleteStruct(self)
}

func (self *Temperature) Save() error {
	if self.Sensor == "" {
		return nil
	}
	node := self.GetNode()
	return node.Save(self)
	//DB.Save(self)
}

type TemperatureSlice []Temperature

func (p *TemperatureSlice) DeleteAll() error {
	for _, e := range *p {
		e.Delete()
	}
	return nil
}

func (p TemperatureSlice) Len() int {
	return len(p)
}

func (p TemperatureSlice) Less(i, j int) bool {
	return p[i].Created.Before(p[j].Created)
}

func (p TemperatureSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *TemperatureSlice) GroupByCreated(t time.Duration) [][]Temperature {
	// !TODO: fix index out of range bug! -> Same for HumiditySlice
	// !HOT
	sort.Sort(TemperatureSlice(*p))
	res := [][]Temperature{}
	if len(*p) < 1 {
		return res
	}

	ti := 0
	lt := time.Time{}
	for i, e := range *p {
		if i == 0 {
			lt = e.Created.Round(t)
		} else if e.Created.After(lt.Add(t)) {
			lt = e.Created.Round(t)
			res = append(res, []Temperature{})
			ti++
		}
		if e.Created.After(lt) && e.Created.Before(lt.Add(t)) {
			if len(res)-1 < ti {
				res = append(res, []Temperature{})
			}
			res[ti] = append(res[ti], e)
		}
	}

	return res
}

func QueryTemperatureData(
	sensor string,
	fromDate time.Time,
	toDate time.Time,
	limit int,
	orderAsc bool,
	statsBucket bool,
	statsTimeGroup time.Duration) (TemperatureSlice, error) {

	//log.Println("QueryTemperatureData() ->", tconv.T2Str(statsTimeGroup))
	res := TemperatureSlice{}

	var node storm.Node

	if statsBucket {
		node = DB.From("stats", "sensd", "temperature", sensor, statsTimeGroup.String())
	} else {
		node = DB.From("sensd", "temperature", sensor)
	}

	query := node.Select(q.And(
		//q.Eq("Sensor", sensor),
		q.Gte("Created", fromDate),
		q.Lte("Created", toDate),
	))

	query.OrderBy("Created")
	if !orderAsc {
		query.Reverse()
	}
	query.Limit(limit)

	query.Find(&res)

	return res, nil
}

type Humidity struct {
	ID      int64     `storm:"id,increment"`
	Sensor  string    `storm:"index"`
	Created time.Time `storm:"index"`
	Value   float64
}

func (self Humidity) GetNode() storm.Node {
	return DB.From("sensd", "humidity", self.Sensor)
}

func (self Humidity) GetCreated() time.Time {
	return self.Created
}

func (self Humidity) GetValue() float64 {
	return self.Value
}

func (self Humidity) GetSensor() string {
	return self.Sensor
}

func (self *Humidity) Delete() error {
	if self.Sensor == "" {
		return nil
	}
	node := self.GetNode()
	//node := DB.From("sensd/humidity")
	return node.DeleteStruct(self)
	//DB.DeleteStruct(self)
}

func (self *Humidity) Save() error {
	if self.Sensor == "" {
		return nil
	}
	node := self.GetNode()
	//node := DB.From("sensd/humidity")
	return node.Save(self)
	//DB.Save(self)
}

type HumiditySlice []Humidity

func (p *HumiditySlice) DeleteAll() error {
	for _, e := range *p {
		e.Delete()
	}
	return nil
}
func (p HumiditySlice) Len() int {
	return len(p)
}

func (p HumiditySlice) Less(i, j int) bool {
	return p[i].Created.Before(p[j].Created)
}

func (p HumiditySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *HumiditySlice) GroupByCreated(t time.Duration) [][]Humidity {
	sort.Sort(HumiditySlice(*p))
	res := [][]Humidity{}
	if len(*p) < 1 {
		return res
	}

	ti := 0
	lt := time.Time{}
	for i, e := range *p {
		if i == 0 {
			lt = e.Created.Round(t)
		} else if e.Created.After(lt.Add(t)) {
			lt = e.Created.Round(t)
			res = append(res, []Humidity{})
			ti++
		}
		if e.Created.After(lt) && e.Created.Before(lt.Add(t)) {
			if len(res)-1 < ti {
				res = append(res, []Humidity{})
			}
			res[ti] = append(res[ti], e)
		}
	}

	return res
}

func QueryHumidityData(
	sensor string,
	fromDate time.Time,
	toDate time.Time,
	limit int,
	orderAsc bool,
	statsBucket bool,
	statsTimeGroup time.Duration) (HumiditySlice, error) {

	//log.Println("QueryHumidityData() ->", tconv.T2Str(statsTimeGroup))

	res := HumiditySlice{}

	var node storm.Node

	if statsBucket {
		node = DB.From("stats", "sensd", "humidity", sensor, statsTimeGroup.String())
	} else {
		node = DB.From("sensd", "humidity", sensor)
	}

	query := node.Select(q.And(
		//q.Eq("Sensor", sensor),
		q.Gte("Created", fromDate),
		q.Lte("Created", toDate),
	))

	query.OrderBy("Created")
	if !orderAsc {
		query.Reverse()
	}
	query.Limit(limit)

	query.Find(&res)

	return res, nil
}
