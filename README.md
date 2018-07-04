GoBox - Growbox automation
---------------------------

1. Features
2. System requirements  
	2.1. Required hardware
3. Installation
4. Documentation  
	4.1. Configuration  
	4.2. Building from source  
	4.2.1	Linux  
	4.2.2	Windows
5. Side notes
6. TODOs


GoBox is a growbox automation app for a Raspberry Pi. It allows you to connect two
DHT11 (Temperature & Humidity) sensors and 2 relays (one for light and one for air).
You can adjust the relay on/off times via web interface and also see the latest
sensor data.


### 1.) Features

* Standalone app no mysql, web server etc. required
* Controlling of 2 Relays and 2 DHT11 sensors
* Switch relays based on boolean expressions
* Standalone app no mysql, web server etc. required
* Password protected web interface with multible users
* Logging of temperature and humidity


![GoBOx screenshot](https://raw.githubusercontent.com/guerillagrow/gobox/master/gobox_screen.png)


## 2.) System requirements

* min. 512 MB RAM, 1024 recommended
* min. 1 GB free storage (SD-Card), > 2 GB free storage recommended
* Linux Raspberry pi OS (debian), but any other linux or unix system should be find too.
* Shell & root access
	
	
## 2.1.) Required hardware

* Raspberry Pi Model 3 (arm6)
* Jumper wires
* 2 DHT11 sensors
* 2 10 A relays (one for light & one for the exhaust fan)
* 2 multiway connectors for your relays
* For convenience a GPIO board
* Of course a exhaust fan and maybe a active coal filter
* Some fans for air circulation and plant movement
* Some lights


## 3.) Installation

1. Extract the gobox_v*.zip / tar file
2. Navigate into the folder containing the extracted files
3. Activate the remote root shell access and create following folders on your raspberry pi:

	`/usr/local/gobox/`   

4. Change the ./conf/app.conf and ./conf/raspberrypi.json file according to your used GPIO pins  
	4.1.	 Set the system time of your raspberrypi properly so the timestamps are correct
5. Change the $RPI_IP variable inside the upload_gobox.sh file to theIP of your raspberry pi then run: 		

	`$ ./upload_gobox.sh`
	
6.	Start gobox with following command:

	```
	$ sudo service gobox start
	```
	
	On older distributions you may run:		
	
	```
	$ /etc/init.d/gobox start
	```
	
7.	If you want gobox as auto start service run following command:	

	```
	$ sudo update-rc.d gobox defaults	
	$ sudo update-rc.d gobox enable
	```
	
Thats it. Happy growing! :-) You can now access the web interface on:
http://[raspberrypi]:8080
		
		
## 4.) Documentation

### Configuration

<b>Note: at the moment the multi user support and user profile features are missing!
The default login is:</b>
```
	user:     root@localhost    (email is used for authentication)
	password: toor
```
<b style="color:red;">Please run gobox only in your trusted LAN due the fact of not yet fully implemented user management capabilities!</b>


There are 2 Configuration files one for the web server and on for the sensor stuff
and relay time configuration.

The `./conf/app.conf` file contains you http and app configuration, you could f.e.
change the http port or the "runmode"

The `./conf/raspberrypi.json` file contains your relay config and 
settings for your sensors. You can set the GPIO pin according to your setting.

Description of config variables in rapspberrypi.json:
```
devices.t1                  => DHT11 Sensor #1
devices.t1.status           => Is this sensor plugged in / used
devices.t1.gpio             => GPIO Data Pin 
devices.t1.read_every       => Value in seconds / Interval to read from sensor
# Same for devices.t2.*

devices.relay_l1                       => Relay config
devices.relay_l1.status                => Is this relay plugged in / used
devices.relay_l1.settings              => Relay settings
devices.relay_l1.settings.condition    => Contition used to turn relay on/off, if not blank 
	                                      the on/off time is ignored and only the condition 
                                          is used to switch the relay
devices.relay_l1.settings.force        => -1: None; 0: Force off; 1: Force on
devices.relay_l1.settings.on           => Time when relay goes on, like: 08:30
devices.relay_l1.settings.off          => Time when relay goes off, like: 20:30
   ...
```

The `sensd` Daemon executable reads the sensor data according to your raspberrypi.json
config file. You can use any kind of script or process to read the sensor data maybe a
custom python script or similar. You just have to enter it in the config file as "sensd_bin".
You custom sensd script must encode the temperature and humidity data as JSON object including following tags:

```
	Sensor  string    `json:"sensor"`
	Type    string    `json:"type"`
	Created time.Time `json:"created"` -> encoded as string
	Value   float64   `json:"value"`
```

### 4.2.) Building from source

#### 4.2.1) Linux

Requirements:

* Go / golang >= 1.10 & configured GOPATH etc.
* gcc compiler (arm-linux-gnueabihf-gcc) only for sensd
	you can replace it with your own script

Run following command to build gobox:

```
$ ./build.sh
```

#### 4.2.2) Windows

Requirements:

* Go / golang >= 1.10 & configured GOPATH etc.
* gcc compiler f.e. MinGW (arm-linux-gnueabihf-gcc) only for sensd
	you can replace it with your own script

	See:
	https://sourceforge.net/projects/mingw-gcc-arm-eabi/files/

Run following command to build gobox:

```
$ build.bat
```


## 5.) Side notes

Keep in mind that the lower as you set the "read_every" (number of seconds) value of 
your sensors the more storage will be consumed by the logged sensor data.
This might also affect the query times for loading graphs etc.

Usually all logged sensor data that is older than a month will be deleted every 24 hours.

Whats comming next? Well I think of better graphs and query options.
So in the next release I think we will get much better performance for the web forntend loading and graph stuff.
There also might come a userfriendly command line setup to make the installation process easier!


## 6.) TODOs

* Add tests
* Add csrf/xsrf protection
* Clean up code base. Make it more idiomatic.
* Maybe replace Beego with gin or echo which are more idiomatic go.
* Finish stats generation (use new config vars) // !HOT
* Extend documentation about raspberrypi.json config file
* Add intelligent relay switch by condition expression evaluation of the config file
	Like:	
	
	```
	($temp_t1 >= 30 && $ton <= $tcurrent)
	// Or:
	($temp_t1 >= 30 && $temp_t2 >= 30 && $relay_l1_status == true)
	
	// Available variables inside an expression:
	// $temp_t1         = Temperature value of sensor T1
	// $temp_t2         = Temperature value of sensor T2
	// $humi_t1         = Humidity value of sensor T1
	// $humi_t2         = Humidity value of sensor T2
	// $tcurrent        = Current time
	// $ton             = Relay time on setting
	// $toff            = Relay time off setting
	// $tlastswitch     = Last time the relay was toggled / switched
	// $relay_l1_status = Current status of relay L1
	// $relay_l2_status = Current status of relay L2
	```


----------------------------------------------------------------------------------

```
 ---------------------------------------------------------------------------------
 -          W                                                                    -
 -         WWW                                                                   -
 -         WWW                                GoBox                              -
 -        WWWWW              the open source growbox automation system           -
 -  W     WWWWW     W                                                            -
 -  WWW   WWWWW   WWW                     Happy growing!                         -
 -   WWW  WWWWW  WWW                                                             -
 -    WWW  WWW  WWW                                                              -
 -     WWW WWW WWW                                                               -
 -       WWWWWWW                                                                 -
 -    WWWW  |  WWWW             http://github.com/guerillagrow/gobox             -
 -          |                                                                    -
 -          |                                                                    -
 ---------------------------------------------------------------------------------
```