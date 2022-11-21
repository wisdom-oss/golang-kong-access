package golangkongaccess

import (
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
	// Check that the port number is in the supported range
	if gatewayAdminAPIPort < 1 || gatewayAdminAPIPort > 65535 {
		return ErrPortOutOfRange
	}
	if useSSL {
		gatewayAPIURL = fmt.Sprintf("https://%s:%d", gatewayHost, gatewayAdminAPIPort)
	} else {
		gatewayAPIURL = fmt.Sprintf("http://%s:%d", gatewayHost, gatewayAdminAPIPort)
	}
	return nil
}
