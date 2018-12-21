package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/just-add-power-hdip-ms/helpers"
	"github.com/labstack/echo"
)

//JustAddPowerChannelResult type for result
type JustAddPowerChannelResult struct {
	Data string `json:"data"`
}

//JustAddPowerChannelIntResult type for result
type JustAddPowerChannelIntResult struct {
	Data int `json:"data"`
}

//SetReceiverToTransmissionChannel change inputs
func SetReceiverToTransmissionChannel(context echo.Context) error {
	log.L.Debugf("Setting receiver to transmitter")

	transmitter := context.Param("transmitter")
	receiver := context.Param("receiver")

	go CheckTransmitterChannel(transmitter)

	log.L.Debugf("Routing %v to %v", receiver, transmitter)

	ipAddress, err := net.ResolveIPAddr("ip", transmitter)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when resolving IP Address ["+transmitter+"]"))
	}

	channel := fmt.Sprintf("%v", ipAddress.IP[3])

	log.L.Debugf("Channel %v", channel)

	result, errrr := helpers.JustAddPowerRequest("http://"+receiver+"/cgi-bin/api/command/channel", channel, "POST")

	if errrr != nil {
		return context.JSON(http.StatusInternalServerError, errrr)
	}

	var jsonResult JustAddPowerChannelResult
	err = json.Unmarshal(result, &jsonResult)

	log.L.Debugf("Result %v", jsonResult)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when unpacking json"))
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmitter})
	}
}

func CheckTransmitterChannel(address string) {
	channel, err := GetTransmissionChannelforAddress(address)

	ipAddress, err2 := net.ResolveIPAddr("ip", address)

	if err == nil && err2 == nil {
		if string(ipAddress.IP[3]) == channel {
			//we're good
			return
		}
	}

	helpers.SetTransmitterChannelForAddress(address)
}

//GetTransmissionChannel retrieves the transmission channel for a just add power device
func GetTransmissionChannel(context echo.Context) error {
	log.L.Debugf("Getting trasnmission channel")

	address := context.Param("address")

	transmissionChannel, err := GetTransmissionChannelforAddress(address)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when unpacking json"))
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmissionChannel})
	}

}

func GetTransmissionChannelforAddress(address string) (string, *nerr.E) {
	log.SetLevel("debug")

	log.L.Debugf("Getting transmitter channel for address %v", address)

	ipAddress, err := net.ResolveIPAddr("ip", address)

	log.L.Debugf("%+v", ipAddress.IP)

	if err != nil {
		return "", nerr.Translate(err).Addf("Error when resolving IP Address [" + address + "]")
	}

	result, errrrrr := helpers.JustAddPowerRequest("http://"+address+"/cgi-bin/api/details/channel", "", "GET")

	if errrrrr != nil {
		log.L.Debugf("%v", err)
		return "", nerr.Translate(errrrrr)
	}

	var jsonResult JustAddPowerChannelIntResult
	err = json.Unmarshal(result, &jsonResult)

	log.L.Debugf("Result %s %v", result, jsonResult)
	log.L.Debugf("len of IP %v", len(ipAddress.IP))

	transmissionChannel := fmt.Sprintf("%v.%v.%v.%v",
		ipAddress.IP[0], ipAddress.IP[1], ipAddress.IP[2], jsonResult.Data)

	return transmissionChannel, nil
}

//SetTransmitterChannel sets the transmission channel for a just add power device
func SetTransmitterChannel(context echo.Context) error {
	log.L.Debugf("Setting transmitter channel")

	transmitter := context.Param("transmitter")

	_, err := helpers.SetTransmitterChannelForAddress(transmitter)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmitter})
	}

}

// JustGetDetailsDevice gets the device details which is used in the hardware stuff.
func JustGetDetailsDevice(context echo.Context) error {
	log.L.Infof("In justGetDeviceDetails")

	address := context.Param("address")

	result, err := helpers.GetDeviceDetails(address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, result)
}

// JustGetSignal gets the signal data from the device
func JustGetSignal(context echo.Context) error {
	log.L.Infof("In justGetDetails")

	address := context.Param("address")

	result, err := helpers.GetDeviceSignal(address)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, result)
}
