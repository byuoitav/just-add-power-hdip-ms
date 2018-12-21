package helpers

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

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

func SetTransmitterChannelForAddress(transmitter string) (string, *nerr.E) {
	ipAddress, err := net.ResolveIPAddr("ip", transmitter)

	if err != nil {
		return "", nerr.Translate(err).Addf("Error when resolving IP Address [" + transmitter + "]")
	}

	log.L.Debugf("Setting transmitter channel %v", transmitter)

	channel := string(ipAddress.IP[3])

	result, err := JustAddPowerRequest("http://"+transmitter+"/cgi-bin/api/command/channel", channel, "POST")

	log.L.Debugf("Result %v", result)

	if err != nil {
		return "", nerr.Translate(err)
	} else {
		return "ok", nil
	}
}
