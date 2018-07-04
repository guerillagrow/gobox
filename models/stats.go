package models

import "time"

import "errors"

import "log"
import "fmt"

//import "github.com/asdine/storm"
import "github.com/asdine/storm/q"

// Statistic functions to aggregate data

var ERR_NOENTRIES error = errors.New("No entries found to process!")

func GenerateTemperatureMedian(sensor string, timeSteps time.Duration, limit int, minSetSize int) (int, error) {

	// !TODO: optimize
	log.Println(fmt.Sprintf("GenerateTemperatureMedian(%s, %s, %d, %d)", sensor, timeSteps, limit, minSetSize))

	if limit < 2 {
		return 0, errors.New(fmt.Sprintf("limit (%s) can't be lower than 2 for median function!", limit))
	}

	node := DB.From("sensd", "temperature", sensor)
	nodeStats := DB.From("stats", "sensd", "temperature", sensor, timeSteps.String()).WithBatch(true)
	lastProcessed := time.Time{}
	err := DB.Get(fmt.Sprintf("stats/sensd/temperature/%s/%s/__meta__", sensor, timeSteps), "last_processed", &lastProcessed)
	if err != nil {
		// Create last processed date and start from first (oldest) entry
		firstEntryList := TemperatureSlice{}
		q1 := node.Select()
		q1.OrderBy("Created") //.Reverse()
		q1.Limit(1)
		q1.Find(&firstEntryList)
		if len(firstEntryList) < 1 {
			log.Println(fmt.Sprintf("ERR GenerateTemperatureMedian(%s, %s, %d, %d) ===> #1 %s", sensor, timeSteps, limit, minSetSize, ERR_NOENTRIES))
			return 0, ERR_NOENTRIES
		}
		firstEntry := firstEntryList[0]
		lastProcessed = firstEntry.Created
	}
	rawStack := TemperatureSlice{}
	q2 := node.Select(q.And(
		q.Gt("Created", lastProcessed),
	))
	q2.Limit(limit)
	q2.OrderBy("Created")
	q2.Find(&rawStack)

	if len(rawStack) < 1 || (len(rawStack) < minSetSize && minSetSize > 0) {
		log.Println(fmt.Sprintf("ERR GenerateTemperatureMedian(%s, %s, %d, %d) ===> #2 %s", sensor, timeSteps, limit, minSetSize, ERR_NOENTRIES))
		return 0, ERR_NOENTRIES
	}
	groupedStack := rawStack.GroupByCreated(timeSteps)

	if len(groupedStack) < 1 {
		log.Println(fmt.Sprintf("ERR GenerateTemperatureMedian(%s, %s, %d, %d) ===> #3 %s", sensor, timeSteps, limit, minSetSize, ERR_NOENTRIES))
		return 0, ERR_NOENTRIES
	}
	ti := 0
	lastElem := Temperature{}
	for gi, gs := range groupedStack {
		if len(gs) < minSetSize && minSetSize > 0 {
			log.Println(fmt.Sprintf("RET GenerateTemperatureMedian(%s, %s, %d, %d) %d ===> #4 %s", sensor, timeSteps, limit, minSetSize, ti, ERR_NOENTRIES))

			return ti, nil
		} else if minSetSize < 1 && len(groupedStack)-1 == gi {
			log.Println(fmt.Sprintf("RET GenerateTemperatureMedian(%s, %s, %d, %d) %d ===> #5 %s", sensor, timeSteps, limit, minSetSize, ti, ERR_NOENTRIES))

			return ti, nil
		}
		if len(gs) < 1 {
			continue
		}
		tmedian := float64(0)
		ttemp := Temperature{}
		for i2, e := range gs {
			if i2 == len(gs)-1 {
				lastElem = e
			}
			tmedian = tmedian + e.Value
		}
		ttemp.Created = gs[0].Created.Round(timeSteps)
		ttemp.Value = tmedian / float64(len(gs))
		ttemp.Sensor = gs[0].Sensor
		nodeStats.Save(&ttemp)
		ti++
	}
	log.Println(fmt.Sprintf("GenerateTemperatureMedian(%s, %s, %d, %d) ===> %d", sensor, timeSteps, limit, minSetSize, ti))

	if !lastElem.Created.IsZero() {
		DB.Set(fmt.Sprintf("stats/sensd/temperature/%s/%s/__meta__", sensor, timeSteps), "last_processed", lastElem.Created)
	}

	nodeStats.Commit()
	return ti, nil
}

func GenerateHumidityMedian(sensor string, timeSteps time.Duration, limit int, minSetSize int) (int, error) {

	log.Println(fmt.Sprintf("GenerateHumidityMedian(%s, %s, %d, %d)", sensor, timeSteps, limit, minSetSize))
	if limit < 2 {
		return 0, errors.New(fmt.Sprintf("limit (%s) can't be lower than 2 for median function!", limit))
	}

	node := DB.From("sensd", "humidity", sensor)
	nodeStats := DB.From("stats", "sensd", "humidity", sensor, timeSteps.String()).WithBatch(true)
	lastProcessed := time.Time{}
	err := DB.Get(fmt.Sprintf("stats/sensd/humidity/%s/%s/__meta__", sensor, timeSteps), "last_processed", &lastProcessed)
	if err != nil || lastProcessed.IsZero() {
		// Create last processed date and start from first (oldest) entry
		firstEntryList := HumiditySlice{}
		q1 := node.Select()
		q1.OrderBy("Created") //.Reverse()
		q1.Limit(1)
		q1.Find(&firstEntryList)
		if len(firstEntryList) < 1 {
			log.Println(fmt.Sprintf("ERR GenerateHumidityMedian(%s, %s, %d, %d) ===> #1 %s", sensor, timeSteps, limit, minSetSize, ERR_NOENTRIES))
			return 0, ERR_NOENTRIES
		}
		firstEntry := firstEntryList[0]
		lastProcessed = firstEntry.Created
	}
	rawStack := HumiditySlice{}
	q2 := node.Select(q.And(
		q.Gt("Created", lastProcessed),
	))
	q2.Limit(limit)
	q2.OrderBy("Created")
	q2.Find(&rawStack)

	if len(rawStack) < 1 || (len(rawStack) < minSetSize && minSetSize > 0) {
		log.Println(fmt.Sprintf("ERR GenerateHumidityMedian(%s, %s, %d, %d) ===> #2 %s", sensor, timeSteps, limit, minSetSize, ERR_NOENTRIES))
		return 0, ERR_NOENTRIES
	}

	groupedStack := rawStack.GroupByCreated(timeSteps)

	if len(groupedStack) < 1 {
		log.Println(fmt.Sprintf("ERR GenerateHumidityMedian(%s, %s, %d, %d) ===> #3 %s", sensor, timeSteps, limit, minSetSize, ERR_NOENTRIES))
		return 0, ERR_NOENTRIES
	}
	ti := 0
	lastElem := Humidity{}
	for gi, gs := range groupedStack {
		if len(gs) < minSetSize && minSetSize > 0 {
			log.Println(fmt.Sprintf("RET GenerateHumidityMedian(%s, %s, %d, %d) %d ===> 4 %s", sensor, timeSteps, limit, minSetSize, ti, ERR_NOENTRIES))

			return ti, nil
		} else if minSetSize < 1 && len(groupedStack)-1 == gi {
			log.Println(fmt.Sprintf("RET GenerateHumidityMedian(%s, %s, %d, %d) %d ===> #5 %s", sensor, timeSteps, limit, minSetSize, ti, ERR_NOENTRIES))

			return ti, nil
		}
		if len(gs) < 1 {
			continue
		}
		tmedian := float64(0)
		ttemp := Humidity{}
		for i2, e := range gs {
			if i2 == len(gs)-1 {
				lastElem = e
			}
			tmedian = tmedian + e.Value
		}
		ttemp.Created = gs[0].Created.Round(timeSteps)
		ttemp.Value = tmedian / float64(len(gs))
		ttemp.Sensor = gs[0].Sensor
		nodeStats.Save(&ttemp)
		ti++
	}
	log.Println(fmt.Sprintf("GenerateHumidityMedian(%s, %s, %d, %d) ===> %d", sensor, timeSteps, limit, minSetSize, ti))

	if !lastElem.Created.IsZero() {
		DB.Set(fmt.Sprintf("stats/sensd/humidity/%s/%s/__meta__", sensor, timeSteps), "last_processed", lastElem.Created)
	}
	nodeStats.Commit()
	return ti, nil
}
