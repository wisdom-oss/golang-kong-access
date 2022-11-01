package golangkongaccess

const (
	ConsistentHashing LoadBalancingAlgorithm = "consistent-hashing"
	LeastConnections  LoadBalancingAlgorithm = "least-connections"
	RoundRobin        LoadBalancingAlgorithm = "round-robin"
)

const (
	Consumer      HashingInput = "consumer"
	IpAddress     HashingInput = "ip"
	Header        HashingInput = "header"
	Cookie        HashingInput = "cookie"
	Path          HashingInput = "path"
	QueryArgument HashingInput = "query_arg"
	UriCapture    HashingInput = "uri_capture"
	None          HashingInput = "none"
)
