package golang_kong_access

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

/*
ReadUpstreamInformation

Read the information about the upstream and return them as an object
*/
func ReadUpstreamInformation(upstreamName string) (*UpstreamConfiguration, error) {
	if gatewayAPIURL == "" {
		return nil, errors.New("the connection to the api gateway was not set up")
	}
	if upstreamName == "" || strings.TrimSpace(upstreamName) == "" {
		return nil, errors.New("empty upstream upstreamName supplied")
	}
	// Create a new upstream information object
	upstreamConfiguration := &UpstreamConfiguration{}
	// Make a http request for information about the upstream
	response, err := http.Get(gatewayAPIURL + "/upstreams/" + upstreamName)
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the gateway")
		return nil, err
	}
	switch response.StatusCode {
	case 200:
		logger.Debug("Got response from the api gateway.")
		jsonParseError := json.NewDecoder(response.Body).Decode(upstreamConfiguration)
		if jsonParseError != nil {
			logger.WithError(jsonParseError).Error("An error occurred while parsing the gateways response")
			return nil, jsonParseError
		}
		return upstreamConfiguration, nil
	case 404:
		logger.WithField("upstream", upstreamName).Warning("The supplied upstream is not configured on the gateway")
		return nil, errors.New("upstream not found")
	default:
		logger.WithFields(log.Fields{
			"upstream":   upstreamName,
			"httpCode":   response.StatusCode,
			"httpStatus": response.Status,
		}).Error("The gateway responded with an unexpected status code")
		return nil, errors.New("unexpected http status")
	}
}
