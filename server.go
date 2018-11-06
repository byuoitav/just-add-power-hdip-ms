package main

import (
	"net/http"
	"os"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/just-add-power-hdip-ms/handlers"
)

func main() {
	port := ":8022"
	router := common.NewRouter()

	log.L.Debugf("Tied to a room system: %v", os.Getenv("ROOM_SYSTEM"))

	// Use the `router` routing group to require authentication
	//router := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	//Functionality endpoints
	router.GET("/input/:transmitter/:receiver", handlers.SetReceiverToTransmissionChannel)

	//Status endpoints
	router.GET("/input/get/:address", handlers.GetTransmissionChannel)

	//Configuration endpoints
	router.PUT("/configure/:transmitter", handlers.SetTransmitterChannel)

	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
