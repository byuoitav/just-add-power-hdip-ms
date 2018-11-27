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
	"github.com/byuoitav/common/v2/auth"
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
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	log.L.Debugf("Setting receiver to transmitter")

	transmitter := context.Param("transmitter")
	receiver := context.Param("receiver")

	go CheckTransmitterChannel(transmitter)

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

func CheckTransmitterChannel(address string) {
	channel, err := GetTransmissionChannelforAddress(address)

	ipAddress, err2 := net.ResolveIPAddr("ip", address)

	if err == nil && err2 == nil {
		if string(ipAddress.IP[15]) == channel {
			//we're good
			return
		}
	}

	SetTransmitterChannelForAddress(address)
}

//GetTransmissionChannel retrieves the transmission channel for a just add power device
func GetTransmissionChannel(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "read-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

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

	log.L.Debugf("Getting transmitter channel for address %v", address)

	ipAddress, err := net.ResolveIPAddr("ip", address)

	log.L.Debugf("%+v", ipAddress.IP)

	if err != nil {
		return "", nerr.Translate(err).Addf("Error when resolving IP Address [" + address + "]")
	}

	result, errrrrr := justAddPowerRequest("http://"+address+"/cgi-bin/api/details/channel", "", "GET")

	if errrrrr != nil {
		log.L.Debugf("%v", err)
		return "", nerr.Translate(errrrrr)
	}

	var jsonResult JustAddPowerChannelIntResult
	err = json.Unmarshal(result, &jsonResult)

	log.L.Debugf("Result %s %v", result, jsonResult)

	transmissionChannel := fmt.Sprintf("%v.%v.%v.%v",
		ipAddress.IP[12], ipAddress.IP[13], ipAddress.IP[14], jsonResult.Data)

	return transmissionChannel, nil
}

//SetTransmitterChannel sets the transmission channel for a just add power device
func SetTransmitterChannel(context echo.Context) error {
	if ok, err := auth.CheckAuthForLocalEndpoints(context, "write-state"); !ok {
		if err != nil {
			log.L.Warnf("Problem getting auth: %v", err.Error())
		}
		return context.String(http.StatusUnauthorized, "unauthorized")
	}

	log.L.Debugf("Setting transmitter channel")

	transmitter := context.Param("transmitter")

	_, err := SetTransmitterChannelForAddress(transmitter)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	} else {
		return context.JSON(http.StatusOK, status.Input{Input: transmitter})
	}

}

func SetTransmitterChannelForAddress(transmitter string) (string, *nerr.E) {
	ipAddress, err := net.ResolveIPAddr("ip", transmitter)

	if err != nil {
		return "", nerr.Translate(err).Addf("Error when resolving IP Address [" + transmitter + "]")
	}

	log.L.Debugf("Setting transmitter channel %v", transmitter)

	channel := string(ipAddress.IP[15])

	result, err := justAddPowerRequest("http://"+transmitter+"/cgi-bin/api/command/channel", channel, "POST")

	log.L.Debugf("Result %v", result)

	if err != nil {
		return "", nerr.Translate(err)
	} else {
		return "ok", nil
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
