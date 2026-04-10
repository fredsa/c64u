package main

// VersionResponse is returned by the Version endpoint.
type VersionResponse struct {
	Version string   `json:"version"`
	Errors  []string `json:"errors"`
}

// Version returns the current version of the REST API.
func (c *Client) Version() (*VersionResponse, error) {
	var resp VersionResponse
	if err := c.get("/v1/version", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
