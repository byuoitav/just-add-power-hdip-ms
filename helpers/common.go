package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

//JustAddPowerRequest .
func JustAddPowerRequest(url string, body string, method string) ([]byte, *nerr.E) {

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

//SetTransmitterChannelForAddress .
func SetTransmitterChannelForAddress(transmitter string) (string, *nerr.E) {
	ipAddress, err := net.ResolveIPAddr("ip", transmitter)
	ipAddress.IP = ipAddress.IP.To4()

	if err != nil {
		return "", nerr.Translate(err).Addf("Error when resolving IP Address [" + transmitter + "]")
	}

	log.L.Debugf("Setting transmitter ipaddr %v", ipAddress)

	channel := fmt.Sprintf("%v", ipAddress.IP[3])

	log.L.Debugf("Setting transmitter channel %+v", channel)

	result, er := JustAddPowerRequest("http://"+transmitter+"/cgi-bin/api/command/channel", channel, "POST")

	log.L.Debugf("Result %v", result)

	if er != nil {
		return "", nerr.Translate(er)
	}

	return "ok", nil
}
