package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/guerillagrow/jstorage"

	"github.com/guerillagrow/gobox/models/common"

	"syscall"
	"time"

	dht "github.com/d2r2/go-dht"
	//i2c "github.com/d2r2/go-i2c"
	//si7021 "github.com/d2r2/go-si7021"
)

// !TODO: Add d1+d2 DHT22 sensors

var Config *jstorage.Storage = jstorage.NewStorage()
var mainQueue chan common.Response

var ARG_CONFIG_FILE *string
var ARG_Debug *bool

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
	ARG_Debug = flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	mainQueue = make(chan common.Response)
	err := Config.LoadFile(*ARG_CONFIG_FILE)
	if err != nil {
		fmt.Println("JSTorage Failure!", err)
		return
	}
	go sensorWorkT("T1", dht.DHT11) // DHT11
	go sensorWorkT("T2", dht.DHT11) // DHT11

	go sensorWorkT("D1", dht.DHT22) // DHT22
	go sensorWorkT("D2", dht.DHT22) // DHT22

	statusS1, _ := Config.GetBool(fmt.Sprintf("devices/%s/status", "s1"))
	statusS2, _ := Config.GetBool(fmt.Sprintf("devices/%s/status", "s2"))

	if statusS1 || statusS2 {
		initSi7021()
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			Config.LoadFile(Config.File)
		}
	}()
	readRoutine()
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

func sensorWorkT(sensorName string, dhtType dht.SensorType) {
	for {
		status, cerr := Config.GetBool(fmt.Sprintf("devices/%s/status", strings.ToLower(sensorName)))
		if !status {
			time.Sleep(10 * time.Second)
			continue
		}
		var res common.Response
		sensorPin, cerr := Config.GetInt(fmt.Sprintf("devices/%s/gpio", strings.ToLower(sensorName)))
		if cerr != nil {
			DebugLogError(cerr)
			time.Sleep(2 * time.Second)
			continue
		}
		tmp, hum, _, err := dht.ReadDHTxxWithRetry(dhtType, int(sensorPin), false, 4)

		if err != nil {
			DebugLogError(err)
			//log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		tc := time.Now()
		res = common.Response{}
		res.Created = tc
		res.Sensor = strings.ToUpper(sensorName)
		res.Type = "t"
		res.Value = float64(tmp)

		select {
		case mainQueue <- res:
		}
		res = common.Response{}
		res.Created = tc
		res.Sensor = strings.ToUpper(sensorName)
		res.Type = "h"
		res.Value = float64(hum)

		select {
		case mainQueue <- res:
		}
		// !DEBUG
		//log.Println("Save sensor data T1.", "Pin:", sensorPin, "Temp:", t1.Value, "Humidity:", h1.Value)

		t1Sleep, _ := Config.GetInt(fmt.Sprintf("devices/%s/read_every", strings.ToLower(sensorName)))
		if t1Sleep < 1 {
			t1Sleep = 30
		}
		time.Sleep(time.Duration(t1Sleep) * time.Second)
	}
}

func initSi7021() {

	// !TODO

	//go sensorWorkS("S1") // Si7021
	//go sensorWorkS("S2") // Si7021
}

func sensorWorkS(sensorName string) {
	for {
		status, cerr := Config.GetBool(fmt.Sprintf("devices/%s/status", strings.ToLower(sensorName)))
		if !status {
			time.Sleep(5 * time.Second)
			continue
		}
		var res common.Response
		sensorPin, cerr := Config.GetInt(fmt.Sprintf("devices/%s/gpio", strings.ToLower(sensorName)))
		if cerr != nil {
			DebugLogError(cerr)
			time.Sleep(2 * time.Second)
			continue
		}
		tmp, hum, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, int(sensorPin), false, 4)

		if err != nil {
			//log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		tc := time.Now()
		res = common.Response{}
		res.Created = tc
		res.Sensor = strings.ToUpper(sensorName)
		res.Type = "t"
		res.Value = float64(tmp)

		select {
		case mainQueue <- res:
		}
		res = common.Response{}
		res.Created = tc
		res.Sensor = strings.ToUpper(sensorName)
		res.Type = "h"
		res.Value = float64(hum)

		select {
		case mainQueue <- res:
		}
		// !DEBUG
		//log.Println("Save sensor data T1.", "Pin:", sensorPin, "Temp:", t1.Value, "Humidity:", h1.Value)

		t1Sleep, _ := Config.GetInt(fmt.Sprintf("devices/%s/read_every", strings.ToLower(sensorName)))
		if t1Sleep < 1 {
			t1Sleep = 30
		}
		time.Sleep(time.Duration(t1Sleep) * time.Second)
	}
}

func DebugLogError(e error) {
	if *ARG_Debug {
		log.Println("[DEBUG]", e)
	}
}

func DebugLog(s string) {
	if *ARG_Debug {
		log.Println("[DEBUG]", s)
	}
}
