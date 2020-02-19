package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/kramer-driver/via"
	"github.com/byuoitav/via-monitor/monitor"
)

/* global variable declaration */
// Changed: lowercase vars
var name string
var deviceList []structs.Device

func init() {

	if len(os.Getenv("ROOM_SYSTEM")) == 0 {
		log.L.Debugf("System is not tied to a specific room. Will not start via monitoring")
		return
	}

	name = os.Getenv("SYSTEM_ID")
	var err error
	fmt.Printf("Gathering information for %s from database\n", name)

	s := strings.Split(name, "-")
	sa := s[0:2]
	room := strings.Join(sa, "-")

	fmt.Printf("Waiting for database . . . .\n")
	for {
		// Pull room information from db
		state, err := db.GetDB().GetStatus()
		log.L.Debugf("%v\n", state)
		//+deploy not_requried
		if (err != nil || state != "completed") && !(len(os.Getenv("DEV_ROUTER")) > 0 || len(os.Getenv("STOP_REPLICATION")) > 0) {
			log.L.Debugf("Database replication in state %v. Retrying in 5 seconds.", state)
			time.Sleep(5 * time.Second)
			continue
		}
		log.L.Debugf("Database replication state: %v", state)

		devices, err := db.GetDB().GetDevicesByRoomAndRole(room, "EventRouter")
		if err != nil {
			log.L.Debugf("Connecting to the Configuration DB failed, retrying in 5 seconds.")
			time.Sleep(5 * time.Second)
			continue
		}

		if len(devices) == 0 {
			//there's a chance that there ARE routers in the room, but the initial database replication is occuring.
			//we're good, keep going
			state, err := db.GetDB().GetStatus()
			if (err != nil || state != "completed") && !(len(os.Getenv("STOP_REPLICATION")) > 0) {
				log.L.Debugf("Database replication in state %v. Retrying in 5 seconds.", state)
				time.Sleep(5 * time.Second)
				continue
			}
		}
		log.L.Debugf("Connection to the Configuration DB established.")
		break
	}
	deviceList, err = db.GetDB().GetDevicesByRoomAndType(room, "via-connect-pro")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func main() {

	port := ":8014"

	var re = regexp.MustCompile(`-CP1$`)
	test := re.MatchString(name)

	//start the VIA monitoring connection if the Controller is CP1
	if test == true && len(os.Getenv("ROOM_SYSTEM")) > 0 {
		for _, device := range deviceList {
			go via.StartMonitoring(device)
		}
	}

	// via functionality endpoints
	write.GET("/via/:address/reset", handlers.ResetVia)
	write.GET("/via/:address/reboot", handlers.RebootVia)

	// Set the volume
	write.GET("/via/:address/volume/set/:volvalue", handlers.SetViaVolume)

	// via informational endpoints
	read.GET("/via/:address/connected", handlers.GetViaConnectedStatus)
	read.GET("/via/:address/volume/level", handlers.GetViaVolume)
	read.GET("/via/:address/hardware", handlers.GetViaHardwareInfo)
	read.GET("/via/:address/active", handlers.GetViaActiveSignal)
	read.GET("/via/:address/roomcode", handlers.GetViaRoomCode)
	read.GET("/via/:address/users/status", handlers.GetConnectedUsers)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
