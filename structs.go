package golang_kong_access

type UpstreamConfiguration struct {
	UpstreamID                string                 `json:"upstreamID"`
	CreationTimestamp         float64                `json:"created_at"`
	Name                      string                 `json:"name"`
	LoadBalancingAlgorithm    LoadBalancingAlgorithm `json:"algorithm"`
	HashInput                 HashingInput           `json:"hash_on"`
	HashInputFallback         HashingInput           `json:"hash_fallback"`
	HashHeader                string                 `json:"hash_on_header"`
	HashHeaderFallback        string                 `json:"hash_fallback_header"`
	HashCookie                string                 `json:"hash_on_cookie"`
	HashCookiePath            string                 `json:"hash_on_cookie_path"`
	HashQueryArgument         string                 `json:"hash_on_query_arg"`
	HashQueryArgumentFallback string                 `json:"hash_fallback_query_arg"`
	HashUriCapture            string                 `json:"hash_on_uri_capture"`
	HashUriCaptureFallback    string                 `json:"hash_fallback_uri_capture"`
	Slots                     int                    `json:"slots"`
	HealthChecks              struct {
		Active struct {
			Timeout   int `json:"timeout"`
			Unhealthy struct {
				Interval     int   `json:"interval"`
				TcpFailures  int   `json:"tcp_failures"`
				Timeouts     int   `json:"timeouts"`
				HttpFailures int   `json:"http_failures"`
				HttpStatuses []int `json:"http_statuses"`
			} `json:"unhealthy"`
			Type        string `json:"type"`
			Concurrency int    `json:"concurrency"`
			Headers     []struct {
				XAnotherHeader []string `json:"x-another-header"`
				XMyHeader      []string `json:"x-my-header"`
			} `json:"headers"`
			Healthy struct {
				Interval     int   `json:"interval"`
				Successes    int   `json:"successes"`
				HttpStatuses []int `json:"http_statuses"`
			} `json:"healthy"`
			HttpPath               string `json:"http_path"`
			HttpsSni               string `json:"https_sni"`
			HttpsVerifyCertificate bool   `json:"https_verify_certificate"`
		} `json:"active"`
		Passive struct {
			Type      string `json:"type"`
			Unhealthy struct {
				HttpStatuses []int `json:"http_statuses"`
				HttpFailures int   `json:"http_failures"`
				Timeouts     int   `json:"timeouts"`
				TcpFailures  int   `json:"tcp_failures"`
			} `json:"unhealthy"`
			Healthy struct {
				HttpStatuses []int `json:"http_statuses"`
				Successes    int   `json:"successes"`
			} `json:"healthy"`
		} `json:"passive"`
		Threshold int `json:"threshold"`
	} `json:"healthchecks"`
	Tags              []string `json:"tags"`
	HostHeader        string   `json:"host_header"`
	ClientCertificate struct {
		Id string `json:"id"`
	} `json:"client_certificate"`
}

type UpstreamTargetInformation struct {
	Id        string  `json:"id"`
	CreatedAt float64 `json:"created_at"`
	Upstream  struct {
		Id string `json:"id"`
	} `json:"upstream"`
	Address string   `json:"target"`
	Weight  int      `json:"weight"`
	Tags    []string `json:"tags"`
}
