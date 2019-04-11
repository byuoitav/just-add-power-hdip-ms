package helpers

import (
	"fmt"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

// JustAddPowerVideoWallParameters type setting the video wall parameters
type JustAddPowerVideoWallParameters struct {
	ColumnPosition int `json:"columnPosition"`
	RowPosition    int `json:"rowPosition"`
	TotalColumns   int `json:"totalColumns"`
	TotalRows      int `json:"totalRows"`
}

//SetVideoWall .
func SetVideoWall(address string, wallParams JustAddPowerVideoWallParameters) (string, *nerr.E) {
	bodyToSend := fmt.Sprintf("[%v,%v,%v,%v]", wallParams.TotalRows, wallParams.TotalColumns, wallParams.RowPosition, wallParams.ColumnPosition)

	log.L.Debugf("Updating video wall for %s with body %v", address, bodyToSend)

	result, er := JustAddPowerRequest("http://"+address+"/cgi-bin/api/command/videowall/layout", bodyToSend, "POST")

	if er != nil {
		return "", nerr.Translate(er)
	}

	log.L.Debugf("Result %s", result)

	return "ok", nil
}
