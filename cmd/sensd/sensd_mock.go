package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guerillagrow/jstorage"

	"github.com/guerillagrow/gobox/models/common"
)

var Config *jstorage.Storage = jstorage.NewStorage()
var mainQueue chan common.Response

var ARG_CONFIG_FILE *string

func main() {

	sigs := make(chan os.Signal, 1)
	//done = make(chan bool, 1)

	go func() {
		sig := <-sigs
		if sig == os.Interrupt || sig == os.Kill || sig == syscall.SIGTERM {
			os.Exit(1)
		}
		//done <- true
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ARG_CONFIG_FILE = flag.String("c", "", "JSON config file")
	flag.Parse()

	mainQueue = make(chan common.Response)
	err := Config.LoadFile(*ARG_CONFIG_FILE)
	if err != nil {
		fmt.Println("JSTorage Failure!", err)
		return
	}
	go sensorT1Work()
	go sensorT2Work()
	go func() {
		for {
			time.Sleep(15 * time.Second)
			Config.LoadFile(Config.File)
		}
	}()
	readRoutine()
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func readRoutine() {
	for {
		e := <-mainQueue
		eb, err := json.Marshal(&e)
		if err == nil {
			fmt.Println(string(eb))
		}
	}
}

func sensorT1Work() {
	for {
		var res common.Response
		tmp := random(15, 40)
		hum := random(45, 100)
		var err error

		if err != nil {
			//log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		res = common.Response{}
		res.Created = time.Now()
		res.Sensor = "T1"
		res.Type = "t"
		res.Value = float64(tmp)

		select {
		case mainQueue <- res:
		}
		res = common.Response{}
		res.Created = time.Now()
		res.Sensor = "T1"
		res.Type = "h"
		res.Value = float64(hum)

		select {
		case mainQueue <- res:
		}
		// !DEBUG
		//log.Println("Save sensor data T1.", "Pin:", sensorPin, "Temp:", t1.Value, "Humidity:", h1.Value)

		t1Sleep, _ := Config.GetInt("devices/t1/read_every")
		if t1Sleep < 1 {
			t1Sleep = 30
		}
		time.Sleep(time.Duration(t1Sleep) * time.Second)
	}
}

func sensorT2Work() {
	for {

		var res common.Response
		tmp := random(15, 40)
		hum := random(45, 100)
		var err error

		if err != nil {
			//log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		res = common.Response{}
		res.Created = time.Now()
		res.Sensor = "T2"
		res.Type = "t"
		res.Value = float64(tmp)

		select {
		case mainQueue <- res:
		}

		res = common.Response{}
		res.Created = time.Now()
		res.Sensor = "T2"
		res.Type = "h"
		res.Value = float64(hum)

		select {
		case mainQueue <- res:
		}
		// !DEBUG
		//log.Println("Save sensor data T2.", "Pin:", sensorPin, "Temp:", t1.Value, "Humidity:", h1.Value)

		t2Sleep, _ := Config.GetInt("devices/t2/read_every")
		if t2Sleep < 1 {
			t2Sleep = 30
		}
		time.Sleep(time.Duration(t2Sleep) * time.Second)
	}
}
