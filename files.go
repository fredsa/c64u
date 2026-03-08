package c64u

import (
	"fmt"
	"net/url"
)

// FileInfo returns basic information about a file on the Ultimate's file system.
// Supports wildcards. The path should start from the root of the file system.
func (c *Client) FileInfo(path string) (map[string]any, error) {
	var resp map[string]any
	if err := c.get("/v1/files/"+url.PathEscape(path)+":info", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateD64 creates a .d64 disk image file on the Ultimate's file system.
// The path should include the full path and filename.
// tracks: number of tracks (35 or 40, default 35 if 0).
// diskName: optional name for the disk header (pass "" to use the filename).
func (c *Client) CreateD64(path string, tracks int, diskName string) (*ErrorResponse, error) {
	params := url.Values{}
	if tracks > 0 {
		params.Set("tracks", fmt.Sprintf("%d", tracks))
	}
	if diskName != "" {
		params.Set("diskname", diskName)
	}
	var resp ErrorResponse
	if err := c.put("/v1/files/"+url.PathEscape(path)+":create_d64", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateD71 creates a .d71 disk image file (70 tracks) on the Ultimate's file system.
// diskName: optional name for the disk header (pass "" to use the filename).
func (c *Client) CreateD71(path string, diskName string) (*ErrorResponse, error) {
	params := url.Values{}
	if diskName != "" {
		params.Set("diskname", diskName)
	}
	var resp ErrorResponse
	if err := c.put("/v1/files/"+url.PathEscape(path)+":create_d71", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateD81 creates a .d81 disk image file (160 tracks) on the Ultimate's file system.
// diskName: optional name for the disk header (pass "" to use the filename).
func (c *Client) CreateD81(path string, diskName string) (*ErrorResponse, error) {
	params := url.Values{}
	if diskName != "" {
		params.Set("diskname", diskName)
	}
	var resp ErrorResponse
	if err := c.put("/v1/files/"+url.PathEscape(path)+":create_d81", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateDNP creates a .dnp disk image file on the Ultimate's file system.
// tracks is required (max 255, each track has 256 sectors, max ~16MB).
// diskName: optional name for the disk header (pass "" to use the filename).
func (c *Client) CreateDNP(path string, tracks int, diskName string) (*ErrorResponse, error) {
	params := url.Values{
		"tracks": {fmt.Sprintf("%d", tracks)},
	}
	if diskName != "" {
		params.Set("diskname", diskName)
	}
	var resp ErrorResponse
	if err := c.put("/v1/files/"+url.PathEscape(path)+":create_dnp", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
