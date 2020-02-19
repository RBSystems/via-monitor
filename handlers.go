package monitor

import (
	"net/http"
	"strconv"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/kramer-driver/via"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func ResetVia(context echo.Context) error {
	address := context.Param("address")

	err := via.Reset(address)
	if err != nil {
		log.L.Debugf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Debugf("Success.")

	return context.JSON(http.StatusOK, "Success")
}

func RebootVia(context echo.Context) error {
	address := context.Param("address")

	err := via.Reboot(address)
	if err != nil {
		log.L.Debugf("There was a problem: %v", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.L.Debugf("Success.")

	return context.JSON(http.StatusOK, "Success")
}

func SetViaVolume(context echo.Context) error {
	address := context.Param("address")
	value := context.Param("volvalue")
	log.L.Debugf("Value passed by SetViaVolume is %v", value)

	volume, err := strconv.Atoi(value)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	} else if volume > 100 || volume < 1 {
		log.L.Debugf("Volume command error - volume value %s is outside the bounds of 1-100", value)
		return context.JSON(http.StatusBadRequest, "Error: volume must be a value from 1 to 100!")
	}

	volumec := strconv.Itoa(volume)
	log.L.Debugf("Setting volume for %s to %v...", address, volume)

	response, err := via.SetVolume(address, volumec)

	if err != nil {
		log.L.Debugf("An Error Occured: %s", err)
		return context.JSON(http.StatusBadRequest, "An error has occured while setting volume")
	}
	log.L.Debugf("Success: %s", response)

	return context.JSON(http.StatusOK, status.Volume{Volume: volume})
}

func GetViaConnectedStatus(context echo.Context) error {
	address := context.Param("address")

	connected := via.IsConnected(address)

	if connected {
		color.Set(color.FgGreen, color.Bold)
		log.L.Debugf("%s is connected", address)
	} else {
		color.Set(color.FgRed)
		log.L.Debugf("%s is not connected", address)
	}

	return context.JSON(http.StatusOK, connected)
}

func GetViaVolume(context echo.Context) error {
	address := context.Param("address")

	ViaVolume, err := via.GetVolume(address)

	if err != nil {
		color.Set(color.FgRed)
		log.L.Debugf("Failed to retreive VIA volume")
		return context.JSON(http.StatusBadRequest, "Failed to retreive VIA volume")
	} else {
		color.Set(color.FgGreen, color.Bold)
		log.L.Debugf("VIA volume is currently set to %v", strconv.Itoa(ViaVolume))
		return context.JSON(http.StatusOK, status.Volume{ViaVolume})
	}

}

// GetViaActiveSignal returns the status of users that are logged in to the VIA
func GetViaActiveSignal(context echo.Context) error {
	signal, err := via.GetActiveSignal(context.Param("address"))
	if err != nil {
		color.Set(color.FgRed)
		log.L.Errorf("Failed to retrieve VIA active signal: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, signal)
}

// GetViaRoomCode - Get the room code of a VIA and return it per request
func GetViaRoomCode(context echo.Context) error {
	code, err := via.GetRoomCode(context.Param("address"))
	if err != nil {
		color.Set(color.FgRed)
		log.L.Errorf("Failed to retrieve VIA room code: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, code)
}

// Get a list of all connected users to the VIA
func GetConnectedUsers(context echo.Context) error {
	userlist, err := via.GetStatusOfUsers(context.Param("address"))
	if err != nil {
		color.Set(color.FgRed)
		log.L.Errorf("Failed to retrieve current user list: %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, userlist)
}
