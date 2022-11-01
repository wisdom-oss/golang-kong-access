package golangkongaccess

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return nil, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(upstreamName) == "" {
		return nil, ErrEmptyFunctionParameter
	}
	// Create a new upstream information object
	upstreamConfiguration := &UpstreamConfiguration{}
	// Make a http request for information about the upstream
	response, err := http.Get(gatewayAPIURL + "/upstreams/" + upstreamName)
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return nil, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the gateway")
		return nil, wrapHttpClientError(err)
	}
	switch response.StatusCode {
	case 200:
		logger.Debug("Got response from the api gateway.")
		jsonParseError := json.NewDecoder(response.Body).Decode(upstreamConfiguration)
		if jsonParseError != nil {
			logger.WithError(jsonParseError).Error("An error occurred while parsing the gateways response")
			return nil, fmt.Errorf("unable to parse json: %w", jsonParseError)
		}
		return upstreamConfiguration, nil
	case 404:
		logger.WithField("upstream", upstreamName).Warning("The supplied upstream is not configured on the gateway")
		return nil, ErrResourceNotFound
	default:
		logger.WithFields(
			log.Fields{
				"upstream":   upstreamName,
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("The gateway responded with an unexpected status code")
		return nil, ErrUnexpectedHttpCode
	}
}

func ReadServiceConfiguration(serviceName string) (*ServiceConfiguration, error) {
	if gatewayAPIURL == "" {
		return nil, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(serviceName) == "" {
		return nil, errors.New("empty serviceName supplied")
	}
	// Create a new instance of the service configuration
	serviceConfiguration := &ServiceConfiguration{}

	// Make a request to the api gateway
	response, err := http.Get(gatewayAPIURL + "/services/" + serviceName)
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return nil, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while sending the request to the api gateway")
		return nil, wrapHttpClientError(err)
	}
	switch response.StatusCode {
	case 200:
		jsonDecodeError := json.NewDecoder(response.Body).Decode(serviceConfiguration)
		if jsonDecodeError != nil {
			logger.WithError(jsonDecodeError).Error("Unable to parse the response sent by the api gateway")
			return nil, fmt.Errorf("unable to parse json: %w", jsonDecodeError)
		}
		return serviceConfiguration, nil
	case 400:
		logger.WithField("httpCode", response.StatusCode).Error("A bad request was made to the api gateway")
		return nil, ErrBadRequest
	case 404:
		logger.WithField("httpCode", response.StatusCode).Error(
			"The supplied service name is not present in the api gateway",
		)
		return nil, ErrResourceNotFound
	default:
		logger.WithFields(
			log.Fields{
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("An unexpected response code was received from the api gateway")
		return nil, ErrUnexpectedHttpCode
	}
}

func ReadRouteConfigurationList(serviceName string) (*RouteConfigurationList, error) {
	if gatewayAPIURL == "" {
		return nil, errors.New("the connection to the api gateway was not set up")
	}
	if serviceName == "" || strings.TrimSpace(serviceName) == "" {
		return nil, errors.New("empty serviceName supplied")
	}
	// Make the request to the api gateway
	response, err := http.Get(gatewayAPIURL + "/services/" + serviceName + "/routes")
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return nil, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while reading the routes from the service")
		return nil, wrapHttpClientError(err)
	}

	switch response.StatusCode {
	case 200:
		routeConfigurationList := &RouteConfigurationList{}
		jsonDecodeError := json.NewDecoder(response.Body).Decode(routeConfigurationList)
		if jsonDecodeError != nil {
			logger.WithError(jsonDecodeError).Error("Unable to parse the response sent by the api gateway")
			return nil, fmt.Errorf("unable to parse json: %w", jsonDecodeError)
		}
		return routeConfigurationList, nil
	case 404:
		logger.WithField("httpCode", response.StatusCode).Error(
			"The supplied service name is not present in the api gateway",
		)
		return nil, ErrResourceNotFound
	default:
		logger.WithFields(
			log.Fields{
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("An unexpected response code was received from the api gateway")
		return nil, ErrUnexpectedHttpCode
	}
}

/*
ReadServicePlugins returns a list of all configured plugins for a service
*/
func ReadServicePlugins(serviceName string) (*PluginList, error) {
	if gatewayAPIURL == "" {
		return nil, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(serviceName) == "" {
		return nil, ErrEmptyFunctionParameter
	}
	// Make the request to the api gateway
	response, err := http.Get(gatewayAPIURL + "/services/" + serviceName + "/plugins")
	bodyCloseErr := response.Body.Close()
	if bodyCloseErr != nil {
		return nil, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while reading the routes from the service")
		return nil, wrapHttpClientError(err)
	}

	switch response.StatusCode {
	case 200:
		pluginList := &PluginList{}
		jsonDecodeError := json.NewDecoder(response.Body).Decode(pluginList)
		if jsonDecodeError != nil {
			logger.WithError(jsonDecodeError).Error("Unable to parse the response sent by the api gateway")
			return nil, jsonDecodeError
		}
		return pluginList, nil
	case 404:
		logger.WithField("httpCode", response.StatusCode).Error(
			"The supplied service name is not present in the api gateway",
		)
		return nil, ErrResourceNotFound
	default:
		logger.WithFields(
			log.Fields{
				"httpCode":   response.StatusCode,
				"httpStatus": response.Status,
			},
		).Error("An unexpected response code was received from the api gateway")
		return nil, ErrUnexpectedHttpCode
	}
}
