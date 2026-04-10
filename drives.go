package main

import (
	"net/url"
)

// DriveInfo represents information about a single drive.
type DriveInfo struct {
	Enabled   bool   `json:"enabled"`
	BusID     int    `json:"bus_id"`
	Type      string `json:"type"`
	ROM       string `json:"rom,omitempty"`
	ImageFile string `json:"image_file,omitempty"`
	ImagePath string `json:"image_path,omitempty"`
	LastError string `json:"last_error,omitempty"`

	Partitions []DrivePartition `json:"partitions,omitempty"`
}

// DrivePartition represents a partition entry for a soft IEC drive.
type DrivePartition struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

// DrivesResponse is returned by ListDrives.
type DrivesResponse struct {
	Drives []map[string]DriveInfo `json:"drives"`
	Errors []string              `json:"errors"`
}

// ListDrives returns information about all internal drives on the IEC bus.
func (c *Client) ListDrives() (*DrivesResponse, error) {
	var resp DrivesResponse
	if err := c.get("/v1/drives", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MountMode specifies the mount mode for a disk image.
type MountMode string

const (
	MountReadWrite MountMode = "readwrite"
	MountReadOnly  MountMode = "readonly"
	MountUnlinked  MountMode = "unlinked"
)

// MountImage mounts a disk image from the Ultimate's file system onto the specified drive.
// Valid drive names are "a", "b", etc.
// imageType is optional (e.g. "d64", "g64", "d71", "g71", "d81") — pass "" to auto-detect from extension.
// mode is optional — pass "" for the default.
func (c *Client) MountImage(drive, image string, imageType string, mode MountMode) (*ErrorResponse, error) {
	params := url.Values{"image": {image}}
	if imageType != "" {
		params.Set("type", imageType)
	}
	if mode != "" {
		params.Set("mode", string(mode))
	}
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":mount", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MountImageData mounts a disk image from attached data onto the specified drive.
// imageType and mode are optional — pass "" to omit.
func (c *Client) MountImageData(drive string, data []byte, filename string, imageType string, mode MountMode) (*ErrorResponse, error) {
	params := url.Values{}
	if imageType != "" {
		params.Set("type", imageType)
	}
	if mode != "" {
		params.Set("mode", string(mode))
	}
	var resp ErrorResponse
	if err := c.postFile("/v1/drives/"+url.PathEscape(drive)+":mount", params, filename, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetDrive resets the specified drive.
func (c *Client) ResetDrive(drive string) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":reset", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveDisk removes the mounted disk from the specified drive.
func (c *Client) RemoveDisk(drive string) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":remove", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnlinkDisk breaks the link between the drive and the mounted image file.
// Further writes will no longer be reflected in the image file.
func (c *Client) UnlinkDisk(drive string) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":unlink", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DriveOn turns on the specified drive. If already on, it is reset.
func (c *Client) DriveOn(drive string) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":on", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DriveOff turns off the specified drive. It will no longer be accessible on the serial bus.
func (c *Client) DriveOff(drive string) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":off", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoadDriveROM loads a ROM file from the Ultimate's file system into the specified drive.
// The ROM must be 16K or 32K depending on the drive type. This is a temporary action.
func (c *Client) LoadDriveROM(drive, file string) (*ErrorResponse, error) {
	params := url.Values{"file": {file}}
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":load_rom", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoadDriveROMData loads a ROM from attached data into the specified drive.
// The ROM must be 16K or 32K depending on the drive type. This is a temporary action.
func (c *Client) LoadDriveROMData(drive string, romData []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postFile("/v1/drives/"+url.PathEscape(drive)+":load_rom", nil, "drive.rom", romData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DriveMode represents a floppy drive mode.
type DriveMode string

const (
	DriveMode1541 DriveMode = "1541"
	DriveMode1571 DriveMode = "1571"
	DriveMode1581 DriveMode = "1581"
)

// SetDriveMode changes the drive mode. This also loads the corresponding ROM.
// A temporary ROM loaded with LoadDriveROM will be lost.
func (c *Client) SetDriveMode(drive string, mode DriveMode) (*ErrorResponse, error) {
	params := url.Values{"mode": {string(mode)}}
	var resp ErrorResponse
	if err := c.put("/v1/drives/"+url.PathEscape(drive)+":set_mode", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
