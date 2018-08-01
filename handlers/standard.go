package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/status"
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

	log.L.Debugf("Routing %v to %v", receiver, transmitter)

	ipAddress, err := net.ResolveIPAddr("ip", transmitter)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when resolving IP Address ["+transmitter+"]"))
	}

	channel := fmt.Sprintf("%v", ipAddress.IP[15])

	log.L.Debugf("Channel %v", channel)

	result, errrr := justAddPowerRequest("http://"+receiver+"/cgi-bin/api/command/channel", channel, "POST")

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

//GetTransmissionChannel retrieves the transmission channel for a just add power device
func GetTransmissionChannel(context echo.Context) error {

	log.L.Debugf("Getting trasnmission channel")

	address := context.Param("address")

	log.L.Debugf("Getting transmitter channel for address %v", address)

	ipAddress, err := net.ResolveIPAddr("ip", address)

	log.L.Debugf("%+v", ipAddress.IP)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when resolving IP Address ["+address+"]"))
	}

	result, errrrrr := justAddPowerRequest("http://"+address+"/cgi-bin/api/details/channel", "", "GET")

	if errrrrr != nil {
		log.L.Debugf("%v", err)
		return context.JSON(http.StatusInternalServerError, errrrrr.Error())
	}

	var jsonResult JustAddPowerChannelIntResult
	err = json.Unmarshal(result, &jsonResult)

	log.L.Debugf("Result %s %v", result, jsonResult)

	transmissionChannel := fmt.Sprintf("%v.%v.%v.%v",
		ipAddress.IP[12], ipAddress.IP[13], ipAddress.IP[14], jsonResult.Data)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when unpacking json"))
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmissionChannel})
	}

}

//SetTransmitterChannel sets the transmission channel for a just add power device
func SetTransmitterChannel(context echo.Context) error {

	log.L.Debugf("Setting transmitter channel")

	transmitter := context.Param("transmitter")

	ipAddress, err := net.ResolveIPAddr("ip", transmitter)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when resolving IP Address ["+transmitter+"]"))
	}

	log.L.Debugf("Setting transmitter channel %v", transmitter)

	channel := string(ipAddress.IP[15])

	result, err := justAddPowerRequest("http://"+transmitter+"/cgi-bin/api/command/channel", channel, "POST")

	log.L.Debugf("Result %v", result)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmitter})
	}
}

func justAddPowerRequest(url string, body string, method string) ([]byte, *nerr.E) {

	var netRequest, err = http.NewRequest(method, url, bytes.NewReader([]byte(body)))

	if err != nil {
		return nil, nerr.Translate(err).Addf("Error when creating new just add power netrequest")
	}

	var netClient = http.Client{
		Timeout: time.Second * 10,
	}

	response, err := netClient.Do(netRequest)

	if err != nil {
		return nil, nerr.Translate(err).Addf("Error when posting to Just add power device")
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, nerr.Translate(err).Addf("Error when reading Just add power device response body")
	}

	if response.StatusCode/100 != 2 {
		return bytes, nerr.Create("Just add power device did not return HTTP OK", "BadResponse")
	}

	return bytes, nil
}
