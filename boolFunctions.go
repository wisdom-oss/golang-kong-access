package golangkongaccess

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

/*
IsUpstreamSetUp

Check if an upstream with the name is configured in the Kong API Gateway
*/
func IsUpstreamSetUp(upstreamName string) (bool, error) {
	// check if the gateway connection was set up
	if strings.TrimSpace(gatewayAPIURL) == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(upstreamName) == "" {
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
			"The gateway responded with 200 OK -> the upstream is configured",
		)
		return true, nil
	case 404:
		logger.WithField("upstream", upstreamName).Warning("The supplied upstream is not configured on the gateway")
		return false, nil
	default:
		logger.WithField("upstream", upstreamName).WithField(
			"httpCode",
			response.StatusCode,
		).Error("The gateway responded with an unexpected status code")
		return false, errors.New("unexpected status code")
	}
}

// IsAddressInUpstreamTargetList checks if the supplied target address is in the target list of the supplied upstream
func IsAddressInUpstreamTargetList(targetAddress string, upstreamName string) (bool, error) {
	if strings.TrimSpace(gatewayAPIURL) == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(targetAddress) == "" {
		return false, errors.New("empty target address")
	}
	if strings.TrimSpace(upstreamName) == "" {
		return false, errors.New("empty upstreamName supplied")
	}
	// Request the targets of the upstream
	response, err := http.Get(gatewayAPIURL + "/upstreams/" + upstreamName + "/targets")
	if err != nil {
		logger.WithError(err).Error("Unable to check if the ip address is listed in upstream targets")
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
			if target.Address == targetAddress {
				return true, nil
			}
		}
		return false, nil
	case 404:
		logger.WithField("upstream", upstreamName).Error("The supplied upstream is not configured on the gateway")
		return false, errors.New("upstream not found")
	default:
		logger.WithField("upstream", upstreamName).WithField(
			"httpCode",
			response.StatusCode,
		).Error("The gateway responded with an unexpected status code")
		return false, errors.New("unexpected status code")
	}
}

/*
IsServiceSetUp checks if a service with the supplied service name exists in the api gateway. The function does not check
if the service is correctly configured. For checking the configuration please use the ServiceHasUpstream function
*/
func IsServiceSetUp(serviceName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, errors.New("empty service name supplied")
	}
	response, err := http.Get(gatewayAPIURL + "/services/" + serviceName)
	if err != nil {
		logger.WithError(err).Error("Unable to check if the service is configured")
		return false, err
	}
	// Check the status code of the response
	switch response.StatusCode {
	case 200:
		return true, nil
	case 404:
		logger.WithField("serviceName", serviceName).Error("The supplied service is not configured on the gateway")
		return false, nil
	default:
		logger.WithField("serviceName", serviceName).WithField(
			"httpCode",
			response.StatusCode,
		).Error("The gateway responded with an unexpected status code")
		return false, errors.New("unexpected status code")
	}
}

/*
ServiceHasUpstream tests if the supplied service name's configuration has the upstream in its host field
*/
func ServiceHasUpstream(serviceName string, upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, errors.New("empty service name supplied")
	}
	if strings.TrimSpace(upstreamName) == "" {
		return false, errors.New("empty upstream name supplied")
	}

	// Use the read function to access the service configuration
	serviceConfiguration, err := ReadServiceConfiguration(serviceName)
	if err != nil {
		logger.WithError(err).Error("An error occurred while reading the service configuration from the api gateway")
		return false, err
	}

	// Now access the service configuration and check the host entry
	return serviceConfiguration.Host == upstreamName, nil
}

/*
ServiceHasRouteSetUp checks if the service has any routes set up.
To check for a route with a path use ServiceHasRouteWithPathSetUp
*/
func ServiceHasRouteSetUp(serviceName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, errors.New("empty service name supplied")
	}
	routeConfigurationList, err := ReadRouteConfigurationList(serviceName)
	if err != nil {
		logger.WithError(err).Error("An error occurred while reading the route configuration list")
		return false, err
	}
	return len(routeConfigurationList.RouteConfigurations) > 0, nil
}

/*
ServiceHasRouteWithPathSetUp checks if the service has a route set up matching the supplied path
HINT: You need to include a leading slash in the path
*/
func ServiceHasRouteWithPathSetUp(serviceName string, path string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, errors.New("empty service name supplied")
	}
	if strings.TrimSpace(path) == "" {
		return false, errors.New("empty path supplied")
	}

	routeConfigurationList, err := ReadRouteConfigurationList(serviceName)
	if err != nil {
		logger.WithError(err).Error("An error occurred while reading the route configuration list")
		return false, err
	}
	for _, routeConfiguration := range routeConfigurationList.RouteConfigurations {
		if stringArrayContains(routeConfiguration.Paths, path) {
			return true, nil
		}
	}
	return false, nil
}

/*
ServiceHasPlugin checks if a service has a plugin set up with the supplied name
*/
func ServiceHasPlugin(serviceName string, pluginName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, errors.New("the connection to the api gateway was not set up")
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, errors.New("empty service name supplied")
	}
	if strings.TrimSpace(pluginName) == "" {
		return false, errors.New("empty plugin name supplied")
	}
	pluginList, err := ReadServicePlugins(serviceName)
	if err != nil {
		logger.WithError(err).Error("An error occurred while reading the plugin list")
		return false, err
	}

	for _, plugin := range pluginList.Plugins {
		if plugin.Name == pluginName {
			return true, nil
		}
	}
	return false, nil
}
