package golangkongaccess

type TargetListResponse struct {
	Targets    []UpstreamTargetInformation `json:"data"`
	NextTarget string                      `json:"next"`
}

type RouteConfigurationList struct {
	RouteConfigurations []RouteConfiguration `json:"data"`
	Next                string               `json:"next"`
}

type PluginList struct {
	Plugins []PluginInformation `json:"data"`
	Next    string              `json:"next"`
}
