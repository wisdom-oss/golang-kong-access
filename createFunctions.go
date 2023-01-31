package golangkongaccess

import (
	"fmt"
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
		return false, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(upstreamName) == "" {
		return false, ErrEmptyFunctionParameter
	}
	// Build the request body for the upstream creation request
	requestBody := url.Values{}
	requestBody.Set("name", upstreamName)
	requestBody.Set("healthchecks.active.http_path", "/ping")
	requestBody.Set("healthchecks.active.timeout", "2")
	requestBody.Add("healthchecks.active.http_statuses", "204")
	requestBody.Set("healthchecks.active.concurrency", "2")
	requestBody.Set("healthchecks.active.healthy.interval", "1")
	requestBody.Set("healthchecks.active.unhealthy.interval", "1")
	// Make the request to the gateway
	response, err := http.PostForm(gatewayAPIURL+"/upstreams", requestBody)
	// Close the response body since it is not being read by the function
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return false, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("the upstream could not be created")
		return false, wrapHttpClientError(err)
	}
	switch response.StatusCode {
	case 201:
		logger.WithField("httpCode", response.StatusCode).Debug("successfully created the upstream in the gateway")
		return true, nil
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error(
			"a upstream with this name already exists in the gateway",
		)
		return false, ErrResourceExists
	default:
		logger.WithFields(
			log.Fields{
				"upstream":   upstreamName,
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("The gateway responded with an unexpected status code")
		return false, ErrUnexpectedHttpCode
	}
}

/*
CreateTargetInUpstream creates a new target with the supplied targetAddress on the supplied upstream
*/
func CreateTargetInUpstream(targetAddress string, upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(targetAddress) == "" {
		return false, ErrEmptyFunctionParameter
	}
	if strings.TrimSpace(upstreamName) == "" {
		return false, ErrEmptyFunctionParameter
	}
	// Build the request body
	requestBody := url.Values{}
	requestBody.Set("target", targetAddress)

	// Send the request body to the gateway
	response, err := http.PostForm(gatewayAPIURL+"/upstreams/"+upstreamName+"/targets", requestBody)
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return false, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		return false, wrapHttpClientError(err)
	}

	switch response.StatusCode {
	case 201:
		return true, nil
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("The request was malformed and could no be acted upon")
		return false, ErrBadRequest
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error("The same target already exists in the upstream")
		return false, ErrResourceExists
	default:
		logger.WithField("httpCode", response.StatusCode).WithField(
			"httpStatus",
			response.Status,
		).Error("Unexpected http status received in response")
		return false, ErrUnexpectedHttpCode
	}
}

/*
CreateService creates a new service object which uses the supplied upstream as host to which the requests are routed
*/
func CreateService(serviceName string, upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, ErrEmptyFunctionParameter
	}

	// Build the request body
	requestBody := url.Values{}
	requestBody.Set("name", serviceName)
	requestBody.Set("host", upstreamName)

	// Make the request to the gateway
	response, err := http.PostForm(gatewayAPIURL+"/services", requestBody)
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the api gateway")
		return false, wrapHttpClientError(err)
	}
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return false, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	switch response.StatusCode {
	case 201:
		return true, nil
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("A bad request was made to the api gateway")
		return false, ErrBadRequest
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error(
			"A service with the same name already exists in the api gateway",
		)
		return false, ErrResourceExists
	default:
		logger.WithFields(
			log.Fields{
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("An unexpected response code was received from the api gateway")
		return false, ErrUnexpectedHttpCode
	}
}

// CreateNewRoute sets up a new route entry in the gateway allowing the service to be reached under the path
func CreateNewRoute(serviceName string, path string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(serviceName) == "" {
		return false, ErrEmptyFunctionParameter
	}
	if strings.TrimSpace(path) == "" {
		return false, ErrEmptyFunctionParameter
	}

	// Build the request body
	requestBody := url.Values{}
	requestBody.Set("paths", path)
	requestBody.Set("protocols", "http")
	requestBody.Set("request_buffering", "false")
	requestBody.Set("response_buffering", "false")

	// Send the request to the api gateway
	response, err := http.PostForm(gatewayAPIURL+"/services/"+serviceName+"/routes", requestBody)
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return false, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the api gateway")
		return false, fmt.Errorf("http client error: %w", err)
	}

	switch response.StatusCode {
	case 200, 201:
		routeCreated, err := ServiceHasRouteWithPathSetUp(serviceName, path)
		if err != nil {
			logger.WithError(err).Error("An error occurred while checking if the route has been created")
			return false, err
		}
		if !routeCreated {
			logger.Error("The route has not been created")
			return false, ErrResourceNotCreated
		}
		return true, nil
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("A bad request was made to the api gateway")
		return false, ErrBadRequest
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error(
			"A route with the same path already exists in the api gateway",
		)
		return false, ErrResourceExists
	default:
		logger.WithFields(
			log.Fields{
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("An unexpected response code was received from the api gateway")
		return false, ErrUnexpectedHttpCode
	}
}

/*
AddPluginToService adds a new plugin to the supplied service with the supplied configuration.
The configuration is required by default. You may disable this by passing true to configurationOptional
*/
func AddPluginToService(serviceName string, pluginName string, pluginConfiguration url.Values) (bool, error) {
	// This is required to us
	if gatewayAPIURL == "" {
		return false, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(serviceName) == "" || strings.TrimSpace(pluginName) == "" {
		return false, ErrEmptyFunctionParameter
	}

	// Build the request body
	requestBody := url.Values{}
	requestBody.Set("name", pluginName)

	// Iterate through the plugin configuration and prepend every configuration option with "config."
	for configurationKey, configurationValueArray := range pluginConfiguration {
		for _, configurationValue := range configurationValueArray {
			requestBody.Add("config."+configurationKey, configurationValue)
		}
	}

	// Send a request to the api gateway
	response, err := http.PostForm(gatewayAPIURL+"/services/"+serviceName+"/plugins", requestBody)
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return false, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the api gateway")
		return false, wrapHttpClientError(err)
	}
	switch response.StatusCode {
	case 201:
		pluginConfigured, err := ServiceHasPlugin(serviceName, pluginName)
		if err != nil {
			logger.WithError(err).Error("An error occurred while checking the plugin creation")
			return false, fmt.Errorf("unable to check plugin creation: %w", err)
		}
		return pluginConfigured, nil
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("A bad request was made to the api gateway")
		return false, ErrBadRequest
	case 409:
		logger.WithField("httpCode", response.StatusCode).Error(
			"The same plugin exists for this service",
		)
		return false, ErrResourceExists
	default:
		logger.WithFields(
			log.Fields{
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("An unexpected response code was received from the api gateway")
		return false, ErrUnexpectedHttpCode
	}
}
