package types

type HealthResp struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ReadyResp struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks,omitempty"`
}
