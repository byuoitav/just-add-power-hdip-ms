package helpers

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

// JustAddPowerDetailsResult type for the hardware details stuff
type JustAddPowerDetailsResult struct {
	Data struct {
		Firmware struct {
			Date   string `json:"date"`
			Update struct {
				Eta      bool   `json:"eta"`
				Message  string `json:"message"`
				Progress bool   `json:"progress"`
				Result   bool   `json:"result"`
				Status   bool   `json:"status"`
			} `json:"update"`
			Version string `json:"version"`
		} `json:"firmware"`
		Model   string `json:"model"`
		Network struct {
			Ipaddress string `json:"ipaddress"`
			Mac       string `json:"mac"`
			Mtu       int    `json:"mtu"`
			Netmask   string `json:"netmask"`
			Speed     string `json:"speed"`
		} `json:"network"`
		Status string `json:"status"`
		Time   string `json:"time"`
		Uptime string `json:"uptime"`
	} `json:"data"`
}

// GetDeviceDetails sends a request to get the Just Add Power device to get IP, and other info
func GetDeviceDetails(address string) (structs.HardwareInfo, *nerr.E) {
	var details structs.HardwareInfo

	addr, e := net.LookupAddr(address)
	if e != nil {
		details.Hostname = address
	} else {
		details.Hostname = strings.Trim(addr[0], ".")
	}

	//Send the request to the Just Add Power API
	result, err := JustAddPowerRequest("http://"+address+"/cgi-bin/api/details/device", "", "GET")
	if err != nil {
		return details, err
	}

	//jsonResult is the response back from the API, it has all the information
	var jsonResult JustAddPowerDetailsResult
	gerr := json.Unmarshal(result, &jsonResult)
	if gerr != nil {
		return details, err
	}

	details.ModelName = jsonResult.Data.Model                  //Model of the device
	details.FirmwareVersion = jsonResult.Data.Firmware.Version //Version of firmware on the device
	details.BuildDate = jsonResult.Data.Firmware.Date          //The Date the firmware was released
	details.PowerStatus = jsonResult.Data.Uptime               //Reports on how long the device has been booted

	// Get the Network info stuff
	details.NetworkInfo.IPAddress = jsonResult.Data.Network.Ipaddress
	details.NetworkInfo.MACAddress = jsonResult.Data.Network.Mac

	return details, nil
}

func GetDeviceSignal(address string) (structs.ActiveSignal, *nerr.E) {
	var signalStatus structs.ActiveSignal

	//Send the request to the Just Add Power API
	result, err := JustAddPowerRequest("http://"+address+"/cgi-bin/api/details/device", "", "GET")
	if err != nil {
		return signalStatus, err
	}

	//jsonResult is the response back from the API, it has all the information
	var jsonResult JustAddPowerDetailsResult
	gerr := json.Unmarshal(result, &jsonResult)
	if gerr != nil {
		return signalStatus, err
	}

	if jsonResult.Data.Status != "Starting Services" || jsonResult.Data.Status != "Streaming Video" || jsonResult.Data.Status != "Decoding Video" {
		signalStatus.Active = false
	} else {
		signalStatus.Active = true
	}

	return signalStatus, nil
}
