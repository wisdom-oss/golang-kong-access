package golang_kong_access

import (
	"errors"
	"net/http"
)

/*
IsUpstreamSetup

Check if an upstream with the name is configured in the Kong API Gateway
*/
func IsUpstreamSetup(upstreamName string) (bool, error) {
	// check if the gateway connection was set up
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	logger.WithField("upstream", upstreamName).Debug("Checking if the upstream is configured on the gateway")
	// Make a http request to the gateway
	response, err := http.Get(gatewayAPIURL + "/upstreams/" + upstreamName)
	if err != nil {
		logger.WithError(err).Error("Unable to check if the upstream exists")
		return false, err
	}
	// Check the status code of the response
	switch response.StatusCode {
	case 200:
		logger.WithField("upstream", upstreamName).Debug(
			"The gateway responded with 200 OK -> the upstream is configured")
		return true, nil
	case 404:
		logger.WithField("upstream", upstreamName).Warning("The supplied upstream is not configured on the gateway")
		return false, nil
	default:
		logger.WithField("upstream", upstreamName).WithField("httpCode",
			response.StatusCode).Error("The gateway responded with an unexpected status code")
		return false, errors.New("unexpected status code")
	}
}
