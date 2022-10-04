package golang_kong_access

const (
	ConsistentHashing LoadBalancingAlgorithm = "consistent-hashing"
	LeastConnections                         = "least-connections"
	RoundRobin                               = "round-robin"
)

const (
	Consumer      HashingInput = "consumer"
	IpAddress                  = "ip"
	Header                     = "header"
	Cookie                     = "cookie"
	Path                       = "path"
	QueryArgument              = "query_arg"
	UriCapture                 = "uri_capture"
	None                       = "none"
)
