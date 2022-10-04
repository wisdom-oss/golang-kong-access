package golang_kong_access

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

/* The URL on which the Kong API Gateway's admin api is running on */
var gatewayAPIURL = ""

/* The logger for this package */
var logger = log.New()

/* A http client which is available for the whole package */
var httpClient = &http.Client{}

/*
SetUpGatewayConnection

Call this function to set up the connection to the API Gateway
*/
func SetUpGatewayConnection(gatewayHost string, gatewayAdminAPIPort int, useSSL bool) error {
	// Default to not using SSL to connect to the API Gateway
	useSSL = false
	protocol := "http"
	// Check that the port number is in the supported range
	if gatewayAdminAPIPort < 1 || gatewayAdminAPIPort > 65535 {
		return errors.New("the given port is outside of the supported range")
	}
	if useSSL == true {
		protocol = "https"
		gatewayAPIURL = fmt.Sprintf("%s://%s:%d", protocol, gatewayHost, gatewayAdminAPIPort)
	} else {
		gatewayAPIURL = fmt.Sprintf("%s://%s:%d", protocol, gatewayHost, gatewayAdminAPIPort)
	}
	return nil
}
