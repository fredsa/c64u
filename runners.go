package main

import (
	"fmt"
	"net/url"
)

// SIDPlay plays a SID file from the Ultimate's file system.
// The optional songNr specifies which song to play (0 for default).
func (c *Client) SIDPlay(file string, songNr int) (*ErrorResponse, error) {
	params := url.Values{"file": {file}}
	if songNr > 0 {
		params.Set("songnr", fmt.Sprintf("%d", songNr))
	}
	var resp ErrorResponse
	if err := c.put("/v1/runners:sidplay", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SIDPlayData plays a SID file attached to the request.
// The optional songNr specifies which song to play (0 for default).
// An optional songLengths attachment can also be provided (pass nil to omit).
func (c *Client) SIDPlayData(sidData []byte, songNr int, songLengths []byte) (*ErrorResponse, error) {
	params := url.Values{}
	if songNr > 0 {
		params.Set("songnr", fmt.Sprintf("%d", songNr))
	}
	// For SID play with possible song lengths, use multipart
	var resp ErrorResponse
	if err := c.postFile("/v1/runners:sidplay", params, "file.sid", sidData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MODPlay plays an Amiga MOD file from the Ultimate's file system.
func (c *Client) MODPlay(file string) (*ErrorResponse, error) {
	params := url.Values{"file": {file}}
	var resp ErrorResponse
	if err := c.put("/v1/runners:modplay", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MODPlayData plays an Amiga MOD file attached to the request.
func (c *Client) MODPlayData(modData []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postFile("/v1/runners:modplay", nil, "file.mod", modData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoadPRG loads a program into C64 memory from the Ultimate's file system.
// The machine resets and loads via DMA but does not automatically run.
func (c *Client) LoadPRG(file string) (*ErrorResponse, error) {
	params := url.Values{"file": {file}}
	var resp ErrorResponse
	if err := c.put("/v1/runners:load_prg", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LoadPRGData loads a program into C64 memory from attached data.
// The machine resets and loads via DMA but does not automatically run.
func (c *Client) LoadPRGData(prgData []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postFile("/v1/runners:load_prg", nil, "file.prg", prgData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RunPRG loads and runs a program from the Ultimate's file system.
// The machine resets, loads via DMA, and automatically runs the program.
func (c *Client) RunPRG(file string) (*ErrorResponse, error) {
	params := url.Values{"file": {file}}
	var resp ErrorResponse
	if err := c.put("/v1/runners:run_prg", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RunPRGData loads and runs a program from attached data.
// The machine resets, loads via DMA, and automatically runs the program.
func (c *Client) RunPRGData(prgData []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postFile("/v1/runners:run_prg", nil, "file.prg", prgData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RunCRT starts a cartridge file from the Ultimate's file system.
// The machine resets with the specified cartridge active.
func (c *Client) RunCRT(file string) (*ErrorResponse, error) {
	params := url.Values{"file": {file}}
	var resp ErrorResponse
	if err := c.put("/v1/runners:run_crt", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RunCRTData starts a cartridge file from attached data.
// The machine resets with the attached cartridge active.
func (c *Client) RunCRTData(crtData []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.postFile("/v1/runners:run_crt", nil, "file.crt", crtData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
