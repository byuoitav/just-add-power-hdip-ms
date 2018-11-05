package main

import (
	"net/http"
	"os"

	"github.com/byuoitav/common/health"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/just-add-power-hdip-ms/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	si "github.com/byuoitav/device-monitoring-microservice/statusinfrastructure"
)

func main() {
	port := ":8022"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	log.L.Debugf("Tied to a room system: %v", os.Getenv("ROOM_SYSTEM"))

	// Use the `router` routing group to require authentication
	//router := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	router.GET("/health", health.HealthCheck)
	router.GET("/mstatus", GetStatus)
	router.GET("/status", GetStatus)

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

// GetStatus returns the status and version number of this instance of the API.
func GetStatus(context echo.Context) error {
	var s si.Status
	var err error

	s.Version, err = si.GetVersion("version.txt")
	if err != nil {
		return context.JSON(http.StatusOK, "Failed to open version.txt")
	}

	s.Status = si.StatusOK
	s.StatusInfo = ""

	return context.JSON(http.StatusOK, s)
}
