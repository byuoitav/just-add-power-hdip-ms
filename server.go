package main

import (
	"net/http"
	"os"

	"github.com/byuoitav/authmiddleware"
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

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	secure.GET("/health", health.HealthCheck)
	secure.GET("/mstatus", GetStatus)
	secure.GET("/status", GetStatus)

	//Functionality endpoints
	secure.GET("/input/:transmitter/:receiver", handlers.SetReceiverToTransmissionChannel)

	//Status endpoints
	secure.GET("/input/get/:address", handlers.GetTransmissionChannel)

	//Configuration endpoints
	secure.PUT("/configure/:transmitter", handlers.SetTransmitterChannel)

	secure.PUT("/log-level/:level", log.SetLogLevel)
	secure.GET("/log-level", log.GetLogLevel)

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
