package models

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"os/exec"
	"sync"
	"time"

	"github.com/guerillagrow/jstorage"

	arrow "github.com/bmuller/arrow/lib"
	//"github.com/d2r2/go-i2c"

	"github.com/asdine/storm"
	sjson "github.com/asdine/storm/codec/json"
	"github.com/asdine/storm/q"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	//"gobot.io/x/gobot/drivers/i2c"
	"github.com/guerillagrow/gobox/models/common"

	"gobot.io/x/gobot/platforms/raspi"
)

// !NOTE: maybe add flag pkg to parse -c as config file flag/arg

var DB *storm.DB
var GoBox *Box
var BoxConfig *jstorage.Storage = jstorage.NewStorage()
var ARG_ConfigFile *string
var ARG_Debug *bool
var ARG_DBFile *string

//var Cron *cron.Cron

func Init() {
	cfgerr := BoxConfig.LoadFile(*ARG_ConfigFile)

	if cfgerr != nil {
		log.Fatalln(cfgerr)
	}

	var dbErr error
	DB, dbErr = storm.Open(*ARG_DBFile, storm.Codec(sjson.Codec))

	if dbErr != nil {
		log.Println("Could not initialize application database (maybe This model is used as library)!", dbErr, "DB file:", *ARG_DBFile)
		return
	}
	InitDatabase()

	devices := []gobot.Device{}
	GoBox = NewBox()
	GoBox.RPIAdaptor = raspi.NewAdaptor()
	//GoBox.LightRelay = gpio.NewLedDriver(GoBox.RPIAdaptor, "4")
	rl1Status, _ := BoxConfig.GetBool("devices/relay_l1/status")
	rl1GPIO, _ := BoxConfig.GetString("devices/relay_l1/gpio")
	if rl1Status == true {
		GoBox.LightRelay = gpio.NewGroveRelayDriver(GoBox.RPIAdaptor, rl1GPIO)
		devices = append(devices, GoBox.LightRelay)
	}
	//GoBox.SensorT1 = i2c.NewBMP180Driver(GoBox.RPIAdaptor)

	/*lcdState, _ := BoxConfig.GetBool("devices/lcd/status")
	 // !OBSOLET: LCD
	if lcdState == true {
		// !TODO: test lcd display functionality!
		var i2cerr error
		var lcderr error
		GoBox.LCDI2C, i2cerr = i2c.NewI2C(0x27, 1)
		checkError(i2cerr)
		GoBox.LCDDevice, lcderr = device.NewLcd(GoBox.LCDI2C, device.LCD_16x2)
		checkError(lcderr)

		GoBox.LCDDevice.BacklightOn()
		GoBox.LCDDevice.Clear()
		GoBox.LCDI2C.Debug = false
		//devices = append(devices, GoBox.LCDDevice)

	}*/
	GoBox.Robot = gobot.NewRobot("bot",
		[]gobot.Connection{GoBox.RPIAdaptor},
		devices,
		GoBox.RobotWork,
	)
	//GoBox.Start()
}

func InitDatabase() {
	//DB.Init(&User{})
	//DB.Init(&Humidity{})
	//DB.Init(&Temperature{})
	u := User{}
	DB.Find("Name", "root", &u)
	if u.ID < 1 {
		err := NewUser("root", "root@localhost", "toor", true)
		checkError(err)
		err = NewUser("grinspoon", "grinspoon@localhost", "toor", false)
		checkError(err)
	}
}

func NewBox() *Box {
	b := Box{}
	b.Init()
	return &b
}

type Box struct {
	LightSchedule string
	Robot         *gobot.Robot
	RPIAdaptor    *raspi.Adaptor
	//LCDI2C        *i2c.I2C // !OBSOLET
	//LCDDevice     *device.Lcd // !OBSOLET
	LightRelay *gpio.GroveRelayDriver
	//FanRelay    *gpio.GroveRelayDriver
	SensDCmd          *exec.Cmd
	SensDStdout       io.ReadCloser
	SensDStdin        io.WriteCloser
	t1Running         bool
	t2Running         bool
	lightRelayLastDay string
	mux               *sync.Mutex
}

func (box *Box) Init() {
	box.mux = &sync.Mutex{}

}

func (box *Box) Start() {
	rl1Status, _ := BoxConfig.GetBool("devices/relay_l1/status")
	if rl1Status == true {
		if box.LightRelay == nil {
			rl1GPIO, _ := BoxConfig.GetString("devices/relay_l1/gpio")
			GoBox.LightRelay = gpio.NewGroveRelayDriver(GoBox.RPIAdaptor, rl1GPIO)
		}
		go box.relayL1Work()
	}
	go box.Robot.Start()
	var cerr error
	sensdBin, _ := BoxConfig.GetString("sensd_bin")
	box.mux.Lock()
	box.SensDCmd = exec.Command(sensdBin, "-c", BoxConfig.File)
	box.mux.Unlock()
	if sensdBin != "" {
		box.SensDStdout, cerr = box.SensDCmd.StdoutPipe()
		checkError(cerr)
		box.SensDStdin, cerr = box.SensDCmd.StdinPipe()
		checkError(cerr)
		serr := box.SensDCmd.Start()
		if serr != nil {
			log.Println("Couldn't start sensD process!", serr)
		}
		go box.ReadSensDPipe()
	}
}

func (box *Box) relayL1Work() {
	var tD string // curent day (yyyy-mm-dd)
	var t string  // current time (hh:ii)
	for {
		tOn, _ := BoxConfig.GetString("devices/relay_l1/settings/on")
		tOff, _ := BoxConfig.GetString("devices/relay_l1/settings/off")
		if tOn == "" {
			log.Fatalln("Missing Time-On / Time-Off parameters in raspberrypi.json config file!")
		}
		if tOn == tOff {
			log.Fatalln("Time-On and Time-Off parameters are the same in raspberrypi.json config file!")
		}

		if box.LightRelay == nil {
			rl1GPIO, _ := BoxConfig.GetString("devices/relay_l1/gpio")
			box.LightRelay = gpio.NewGroveRelayDriver(box.RPIAdaptor, rl1GPIO)
		}

		tD = arrow.Now().CFormat("%Y-%m-%d")
		t = arrow.Now().CFormat("%H:%M")

		if tOn > tOff && (box.lightRelayLastDay == "" || box.lightRelayLastDay < tD) {
			if t >= tOff && box.lightRelayLastDay != "" && box.LightState() == true {
				box.LightOff()
				box.lightRelayLastDay = tD
			} else if t >= tOn && box.LightState() == false {
				box.LightOn()
				box.lightRelayLastDay = tD
			}
		} else if tOn < tOff {
			if t >= tOff || t < tOn {
				box.LightOff()
				box.lightRelayLastDay = tD
			} else if t >= tOn && box.LightState() == false {
				box.LightOn()
				box.lightRelayLastDay = tD
			}
		}
		time.Sleep(10 * time.Second)

	}
}

func (box *Box) ReadSensDPipe() {

	r := bufio.NewReader(box.SensDStdout)

	for {
		line, _, err := r.ReadLine()

		if err == io.EOF {
			// exit goroutine if sensD was shutdown
			return
		}

		if err != nil {
			checkError(err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		res := common.Response{}
		jerr := json.Unmarshal(line, &res)
		if jerr != nil {
			checkError(jerr)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if res.Sensor != "" && res.Type != "" {
			//log.Println("// !DEBUG", "Got sensor data:", res)
			if res.Type == "h" {
				t := Humidity{}
				t.Created = res.Created
				t.Sensor = res.Sensor
				t.Value = float64(res.Value)
				t.Save()
				//log.Println("// !DEBUG", "Saved sensor data!", res)
			} else if res.Type == "t" {
				t := Temperature{}
				t.Created = res.Created
				t.Sensor = res.Sensor
				t.Value = float64(res.Value)
				t.Save()
				//log.Println("// !DEBUG", "Saved sensor data!", res)
			}
		}

	}
}

func (box *Box) RobotWork() {
	gobot.Every(5*time.Second, func() {
		BoxConfig.LoadFile(BoxConfig.File)
	})

	// !TODO: make it more intelligent
	// !DEBUG // !DEV
	// Stats routines
	gobot.Every((60+20)*time.Second, func() { // 1m
		GenerateHumidityMedian("T1", 1*time.Minute, 1000, 0)
		GenerateHumidityMedian("T2", 1*time.Minute, 1000, 0)
		GenerateTemperatureMedian("T1", 1*time.Minute, 1000, 0)
		GenerateTemperatureMedian("T2", 1*time.Minute, 1000, 0)
	})
	gobot.Every((10+2)*time.Minute, func() { // 10m
		GenerateHumidityMedian("T1", 10*time.Minute, 1000, 0)
		GenerateHumidityMedian("T2", 10*time.Minute, 1000, 0)
		GenerateTemperatureMedian("T1", 10*time.Minute, 1000, 0)
		GenerateTemperatureMedian("T2", 10*time.Minute, 1000, 0)
	})
	gobot.Every(62*time.Minute, func() { //1h
		GenerateHumidityMedian("T1", 1*time.Hour, 3000, 0)
		GenerateHumidityMedian("T2", 1*time.Hour, 3000, 0)
		GenerateTemperatureMedian("T1", 1*time.Hour, 3000, 0)
		GenerateTemperatureMedian("T2", 1*time.Hour, 3000, 0)
	})

	// delete old sensor metrics
	gobot.Every(1*time.Hour, func() { //1h
		qlimit := 2000
		dolderthan, _ := BoxConfig.GetInt64("stats/delete_older_than") // value in hours
		if dolderthan < 1 {
			dolderthan = 24 * 3 // default 3 days
		}
		qtimel := time.Now().Add(-(time.Duration(dolderthan) * time.Hour))

		t := Temperature{Sensor: "T1"}
		query := t.GetNode().Select(q.And(
			q.Lt("Created", qtimel),
		)).Limit(qlimit)
		query.Each(&Temperature{}, func(v interface{}) error {
			e, nok := v.(*Temperature)
			if !nok {
				return nil
			}
			e.Delete()
			return nil
		})

		t = Temperature{Sensor: "T2"}
		query = t.GetNode().Select(q.And(
			q.Lt("Created", qtimel),
		)).Limit(qlimit)
		query.Each(&Temperature{}, func(v interface{}) error {
			e, nok := v.(*Temperature)
			if !nok {
				return nil
			}
			e.Delete()
			return nil
		})

		h := Humidity{Sensor: "T1"}
		query = h.GetNode().Select(q.And(
			q.Lt("Created", qtimel),
		)).Limit(qlimit)
		query.Each(&Humidity{}, func(v interface{}) error {
			e, nok := v.(*Humidity)
			if !nok {
				return nil
			}
			e.Delete()
			return nil
		})

		h = Humidity{Sensor: "T2"}
		query = h.GetNode().Select(q.And(
			q.Lt("Created", qtimel),
		)).Limit(qlimit)
		query.Each(&Humidity{}, func(v interface{}) error {
			e, nok := v.(*Humidity)
			if !nok {
				return nil
			}
			e.Delete()
			return nil
		})

	})

	/*gobot.Every(5*time.Minute, func() {
		debug.FreeOSMemory()
	})*/

}

func (box *Box) RobotState() bool {
	return box.Robot.Running()
}

func (box *Box) LightState() bool {
	return box.LightRelay.State()
}

func (box *Box) LightOn() error {
	box.mux.Lock()
	err := box.LightRelay.On()
	box.mux.Unlock()
	return err
}

func (box *Box) LightOff() error {
	box.mux.Lock()
	err := box.LightRelay.Off()
	box.mux.Unlock()
	return err
}

func (box *Box) LightToggle() error {
	box.mux.Lock()
	err := box.LightRelay.Toggle()
	box.mux.Unlock()
	return err
}

func (box *Box) Stop() error {
	box.Robot.Stop()
	if box.SensDCmd != nil && box.SensDCmd.Process != nil {
		box.SensDCmd.Process.Kill()
	}
	time.Sleep(1 * time.Second)
	//DB.Commit()
	DB.Close()
	return nil
}

func checkError(e error) {
	if e == nil {
		return
	}
	log.Println(e)
}
