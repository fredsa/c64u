package c64u

import (
	"encoding/json"
	"net/url"
)

// CategoriesResponse is returned by ListCategories.
type CategoriesResponse struct {
	Categories []string `json:"categories"`
	Errors     []string `json:"errors"`
}

// ListCategories returns all configuration categories.
func (c *Client) ListCategories() (*CategoriesResponse, error) {
	var resp CategoriesResponse
	if err := c.get("/v1/configs", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetConfig returns configuration items for the given category. Wildcards are allowed.
// The result is a map of category -> item -> value.
func (c *Client) GetConfig(category string) (map[string]any, error) {
	var resp map[string]any
	if err := c.get("/v1/configs/"+url.PathEscape(category), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetConfigItem returns detailed information about a specific config item.
// Both category and item support wildcards.
func (c *Client) GetConfigItem(category, item string) (map[string]any, error) {
	var resp map[string]any
	path := "/v1/configs/" + url.PathEscape(category) + "/" + url.PathEscape(item)
	if err := c.get(path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SetConfigItem sets a specific configuration item to the given value.
// Both category and item support wildcards.
func (c *Client) SetConfigItem(category, item, value string) (*ErrorResponse, error) {
	path := "/v1/configs/" + url.PathEscape(category) + "/" + url.PathEscape(item)
	params := url.Values{"value": {value}}
	var resp ErrorResponse
	if err := c.put(path, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetConfigs sets multiple configuration items at once.
// The data should be a map of category -> item -> value, e.g.:
//
//	map[string]any{
//	    "Drive A Settings": map[string]any{
//	        "Drive": "Enabled",
//	        "Drive Bus ID": 8,
//	    },
//	}
func (c *Client) SetConfigs(data map[string]any) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postJSON("/v1/configs", nil, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetConfigsJSON sets multiple configuration items from raw JSON.
func (c *Client) SetConfigsJSON(data json.RawMessage) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postJSON("/v1/configs", nil, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoadConfigFromFlash restores the configuration from non-volatile memory.
func (c *Client) LoadConfigFromFlash() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/configs:load_from_flash", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SaveConfigToFlash saves the current configuration to non-volatile memory.
func (c *Client) SaveConfigToFlash() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/configs:save_to_flash", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetConfigToDefault resets the current settings to factory defaults.
// This does not affect values stored in non-volatile memory.
func (c *Client) ResetConfigToDefault() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/configs:reset_to_default", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
