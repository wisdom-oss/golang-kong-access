package golang_kong_access

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

func UpdateServiceHost(serviceName string, newHost string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if serviceName == "" || strings.TrimSpace(serviceName) == "" {
		return false, errors.New("empty service name supplied")
	}
	if newHost == "" || strings.TrimSpace(newHost) == "" {
		return false, errors.New("empty new host supplied")
	}

	// Build the request body
	requestBody := url.Values{}
	requestBody.Set("host", newHost)

	// Build the PATCH request required by the api gateway
	request, err := http.NewRequest("PATCH", gatewayAPIURL+"/services/"+serviceName,
		strings.NewReader(requestBody.Encode()))
	if err != nil {
		logger.WithError(err).Error("An error occurred while building the request to update the service entry")
		return false, err
	}

	// Set the correct Content-Type header for the gateway
	request.Header.Set("Content-Type", "application/x-www-urlencoded")

	// Send the request to the gateway
	response, err := httpClient.Do(request)

	if err != nil {
		logger.WithError(err).Error("An error occurred while updating the host of the service")
		return false, err
	}

	switch response.StatusCode {
	case 200:
		// Check if the host is set correctly
		hostIsSet, err := ServiceHasUpstream(serviceName, newHost)
		if err != nil {
			logger.WithError(err).Error("An error occurred while validating the update of the host")
			return false, err
		}
		if !hostIsSet {
			logger.Error("The service host update was not successful")
			return false, errors.New("update unsuccessful")
		} else {
			return true, nil
		}
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("A bad request was made to the api gateway")
		return false, errors.New("bad request sent to the gateway")
	case 404:
		logger.WithField("httpCode", response.StatusCode).Error(
			"The supplied service name is not present in the api gateway")
		return false, errors.New("service not found")
	default:
		logger.WithFields(log.Fields{"httpCode": response.StatusCode,
			"httpStatus": response.Status}).Error("An unexpected response code was received from the api gateway")
		return false, errors.New("unexpected response code")
	}
}
