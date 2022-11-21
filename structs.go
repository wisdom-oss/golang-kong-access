package golangkongaccess

// UpstreamConfiguration specifies how the upstream is configured in the gateway.
// This model is only used for reading such information from the gateway
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

// UpstreamTargetInformation contains all information the gateway has stored about a upstream target.
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

// ServiceConfiguration contains all data stored in the gateway for the service configuration
type ServiceConfiguration struct {
	Id                string   `json:"id"`
	CreationTimestamp float64  `json:"created_at"`
	UpdateTimestamp   float64  `json:"updated_at"`
	Name              string   `json:"name"`
	Retries           int      `json:"retries"`
	Protocol          string   `json:"protocol"`
	Host              string   `json:"host"`
	Port              int      `json:"port"`
	Path              string   `json:"path"`
	ConnectTimeout    int      `json:"connect_timeout"`
	WriteTimeout      int      `json:"write_timeout"`
	ReadTimeout       int      `json:"read_timeout"`
	Tags              []string `json:"tags"`
	ClientCertificate struct {
		Id string `json:"id"`
	} `json:"client_certificate"`
	TlsVerify      bool        `json:"tls_verify"`
	TlsVerifyDepth interface{} `json:"tls_verify_depth"`
	CaCertificates []string    `json:"ca_certificates"`
	Enabled        bool        `json:"enabled"`
}

// RouteConfiguration contains all information about a route set up in the gateways
type RouteConfiguration struct {
	Id                      string            `json:"id"`
	CreationTimestamp       float64           `json:"created_at"`
	UpdateTimestamp         float64           `json:"updated_at"`
	Name                    string            `json:"name"`
	Protocols               []string          `json:"protocols"`
	Methods                 []string          `json:"methods"`
	Hosts                   []string          `json:"hosts"`
	Paths                   []string          `json:"paths"`
	Headers                 map[string]string `json:"headers"`
	HttpsRedirectStatusCode int               `json:"https_redirect_status_code"`
	RegexPriority           int               `json:"regex_priority"`
	StripPath               bool              `json:"strip_path"`
	PathHandling            string            `json:"path_handling"`
	PreserveHost            bool              `json:"preserve_host"`
	RequestBuffering        bool              `json:"request_buffering"`
	ResponseBuffering       bool              `json:"response_buffering"`
	Tags                    []string          `json:"tags"`
	Service                 struct {
		Id string `json:"id"`
	} `json:"service"`
}

// PluginInformation contains the information about a plugin installed on either a service, path or globally
type PluginInformation struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	CreationTimestamp int                    `json:"created_at"`
	Route             interface{}            `json:"route"`
	Service           interface{}            `json:"service"`
	Consumer          interface{}            `json:"consumer"`
	Configuration     map[string]interface{} `json:"config"`
	Protocols         []string               `json:"protocols"`
	Enabled           bool                   `json:"enabled"`
	Tags              []string               `json:"tags"`
}
