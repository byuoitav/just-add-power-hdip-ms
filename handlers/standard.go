package handlers

import (
	"bytes"
	"encoding/json"
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

	channel := string(ipAddress.IP[3])

	log.L.Debugf("Channel %v", channel)

	result, err := justAddPowerRequest("http://"+receiver+"/cgi-bin/api/command/channel", channel, "POST")

	log.L.Debugf("Result %v", result)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
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

	if err != nil {
		return context.JSON(http.StatusInternalServerError, nerr.Translate(err).Addf("Error when resolving IP Address ["+address+"]"))
	}

	result, err := justAddPowerRequest("http://"+address+"/cgi-bin/api/details/channel", "", "GET")

	log.L.Debugf("Result %v", result)

	transmissionChannel := string(ipAddress.IP[0]) + "." + string(ipAddress.IP[1]) + "." + string(ipAddress.IP[2]) + "." + result.Data

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
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

	channel := string(ipAddress.IP[3])

	result, err := justAddPowerRequest("http://"+transmitter+"/cgi-bin/api/command/channel", channel, "POST")

	log.L.Debugf("Result %v", result)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmitter})
	}
}

func justAddPowerRequest(url string, body string, method string) (JustAddPowerChannelResult, *nerr.E) {
	var retValue JustAddPowerChannelResult

	var netRequest, err = http.NewRequest(method, url, bytes.NewReader([]byte(body)))

	if err != nil {
		return retValue, nerr.Translate(err).Addf("Error when creating new just add power netrequest")
	}

	var netClient = http.Client{
		Timeout: time.Second * 10,
	}

	response, err := netClient.Do(netRequest)

	if err != nil {
		return retValue, nerr.Translate(err).Addf("Error when posting to Just add power device")
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return retValue, nerr.Translate(err).Addf("Error when reading Just add power device response body")
	}

	err = json.Unmarshal(bytes, &retValue)

	if err != nil {
		return retValue, nerr.Translate(err).Addf("Unable to unpackage Just add power device response body")
	}

	if response.StatusCode/100 != 2 {
		return retValue, nerr.Create("Just add power device did not return HTTP OK", "BadResponse")
	}

	return retValue, nil
}
