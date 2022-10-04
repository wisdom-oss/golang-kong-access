package golang_kong_access

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
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
	if upstreamName == "" || strings.TrimSpace(upstreamName) == "" {
		return false, errors.New("empty upstreamName supplied")
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

func IsIPv4AddressInUpstreamTargetList(ipAddress string, upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if ipAddress == "" || strings.TrimSpace(ipAddress) == "" {
		return false, errors.New("empty ip address")
	}
	// Check if the supplied ip address is an ipv4 address
	match, regexError := regexp.MatchString("(?:[0-9]{1,3}\\.){3}[0-9]{1,3}", ipAddress)
	if regexError != nil {
		return false, regexError
	}
	if !match {
		return false, errors.New("invalid ipv4 address")
	}
	if upstreamName == "" || strings.TrimSpace(upstreamName) == "" {
		return false, errors.New("empty upstreamName supplied")
	}
	//Request the targets of the upstream
	response, err := http.Get(gatewayAPIURL + "/upstreams/" + upstreamName + "/targets")
	if err != nil {
		return false, err
	}
	// Check the status code of the response
	switch response.StatusCode {
	case 200:
		targetListResponse := &TargetListResponse{}
		jsonParseError := json.NewDecoder(response.Body).Decode(targetListResponse)
		if jsonParseError != nil {
			return false, jsonParseError
		}
		for _, target := range targetListResponse.Targets {
			if target.Address == ipAddress {
				return true, nil
			}
		}
		return false, nil
	case 404:
		logger.WithField("upstream", upstreamName).Error("The supplied upstream is not configured on the gateway")
		return false, errors.New("upstream not found")
	default:
		logger.WithField("upstream", upstreamName).WithField("httpCode",
			response.StatusCode).Error("The gateway responded with an unexpected status code")
		return false, errors.New("unexpected status code")
	}
}
