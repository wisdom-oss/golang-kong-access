package golangkongaccess

import (
	"fmt"
	"net/http"
	"strings"
)

func DeleteUpstreamTarget(targetAddress string, upstreamName string) (bool, error) {
	if gatewayAPIURL == "" {
		return false, ErrConnectionNotSetUp
	}
	if strings.TrimSpace(targetAddress) == "" {
		return false, ErrEmptyFunctionParameter
	}
	if strings.TrimSpace(upstreamName) == "" {
		return false, ErrEmptyFunctionParameter
	}
	request, err := http.NewRequest(
		"DELETE", gatewayAPIURL+"/upstreams/"+upstreamName+"/targets/"+targetAddress,
		nil,
	)
	if err != nil {
		logger.WithError(err).Error("An error occurred while building the request to delete the target")
		return false, wrapHttpClientError(err)
	}

	response, err := httpClient.Do(request)
	if bodyCloseErr := response.Body.Close(); bodyCloseErr != nil {
		return false, fmt.Errorf("error while closing response body: %w", bodyCloseErr)
	}
	if err != nil {
		logger.WithError(err).Error("An error occurred while deleting the target from the upstream")
		return false, wrapHttpClientError(err)
	}

	if response.StatusCode != 204 {
		return false, ErrUnexpectedHttpCode
	}
	return true, nil
}
