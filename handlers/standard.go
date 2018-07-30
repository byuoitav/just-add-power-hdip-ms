package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/labstack/echo"
)

//JustAddPowerChannelResult type for result
type JustAddPowerChannelResult struct {
	Data string `json:"data"`
}

//SetReceiverToTransmissionChannel change inputs
func SetReceiverToTransmissionChannel(context echo.Context) error {
	log.L.Debugf("Setting receiver to transmitter")

	transmitter := context.Param("transmitter")
	receiver := context.Param("receiver")

	log.L.Debugf("Routing %v to %v", receiver, transmitter)

	channel := strings.Split(transmitter, ".")[3]

	response, err := http.Post(
		"http://"+receiver+"/cgi-bin/api/command/channel",
		"",
		bytes.NewReader([]byte(channel)))

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	var result JustAddPowerChannelResult

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if result.Data != "OK" {
		return context.JSON(http.StatusInternalServerError, "Just add power device did not return OK")
	}

	if response.StatusCode/100 != 2 {
		return context.JSON(http.StatusInternalServerError, "Just add power device did not return HTTP OK")
	}

	return context.JSON(http.StatusOK, status.Input{Input: transmitter})
}

func GetRecieverTrasmissionChannel(context echo.Context) error {

	return context.JSON(http.StatusInternalServerError, "Not implemented yet")
}

func SetTransmitterChannel(context echo.Context) error {

	log.L.Debugf("Setting transmitter channel")

	transmitter := context.Param("transmitter")

	log.L.Debugf("Setting transmitter channel %v", transmitter)

	channel := strings.Split(transmitter, ".")[3]

	response, err := http.Post(
		"http://"+transmitter+"/cgi-bin/api/command/channel",
		"",
		bytes.NewReader([]byte(channel)))

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	var result JustAddPowerChannelResult

	err = json.Unmarshal(bytes, &result)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	if result.Data != "OK" {
		return context.JSON(http.StatusInternalServerError, "Just add power device did not return OK")
	}

	if response.StatusCode/100 != 2 {
		return context.JSON(http.StatusInternalServerError, "Just add power device did not return HTTP OK")
	}

	return context.JSON(http.StatusOK, status.Input{Input: transmitter})
}
