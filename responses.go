package golang_kong_access

type TargetListResponse struct {
	Targets    []UpstreamTargetInformation `json:"data"`
	NextTarget string                      `json:"next"`
}
