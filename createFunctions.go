package golang_kong_access

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

/*
CreateNewUpstream

Create a new upstream object in the api gateway. If the function returns true the upstream has been created.
If errors are returned by other used functions the error will be returned
*/
func CreateNewUpstream(upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if upstreamName == "" || strings.TrimSpace(upstreamName) == "" {
		return false, errors.New("empty upstream upstreamName supplied")
	}
	// Build the request body for the upstream creation request
	requestBody := url.Values{}
	requestBody.Set("upstreamName", upstreamName)
	// Make the request to the gateway
	response, err := http.PostForm(gatewayAPIURL+"/upstreams", requestBody)
	if err != nil {
		logger.WithError(err).Error("the upstream could not be created")
		return false, err
	}
	switch response.StatusCode {
	case 201:
		logger.WithField("httpCode", response.StatusCode).Debug("successfully created the upstream in the gateway")
		return true, nil
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error(
			"a upstream with this name already exists in the gateway")
		return false, errors.New("upstream already exists")
	default:
		logger.WithFields(log.Fields{
			"upstream":   upstreamName,
			"httpCode":   response.StatusCode,
			"httpStatus": response.Status,
		}).Error("The gateway responded with an unexpected status code")
		return false, errors.New("unexpected http status")
	}
}

/*
CreateTargetInUpstream creates a new target with the supplied targetAddress on the supplied upstream
*/
func CreateTargetInUpstream(targetAddress string, upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if targetAddress == "" || strings.TrimSpace(targetAddress) == "" {
		return false, errors.New("empty targetAddress supplied")
	}
	if upstreamName == "" || strings.TrimSpace(upstreamName) == "" {
		return false, errors.New("empty upstreamName supplied")
	}
	// Build the request body
	requestBody := url.Values{}
	requestBody.Set("targetAddress", targetAddress)

	// Send the request body to the gateway
	response, err := http.PostForm(gatewayAPIURL+"/upstreams/"+upstreamName+"/targets", requestBody)
	if err != nil {
		return false, err
	}

	switch response.StatusCode {
	case 201:
		return true, nil
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("The request was malformed and could no be acted upon")
		return false, errors.New("bad request made")
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error("The same target already exists in the upstream")
		return false, errors.New("target already exists in the upstream")
	default:
		logger.WithField("httpCode", response.StatusCode).WithField("httpStatus",
			response.Status).Error("Unexpected http status received in response")
		return false, errors.New("unexpected http status")
	}
}
