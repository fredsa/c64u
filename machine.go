package c64u

import (
	"fmt"
	"net/url"
)

// Reset sends a reset to the machine without changing configuration.
func (c *Client) Reset() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/machine:reset", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Reboot restarts the machine, re-initializing the cartridge configuration and sending a reset.
func (c *Client) Reboot() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/machine:reboot", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Pause pauses the machine by pulling the DMA line low, stopping the CPU.
// Note: timers are not stopped.
func (c *Client) Pause() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/machine:pause", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resume resumes the machine from a paused state.
func (c *Client) Resume() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/machine:resume", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PowerOff powers off the machine. U64 only.
// Note: a valid response may not be received.
func (c *Client) PowerOff() (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/machine:poweroff", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WriteMem writes data to C64 memory via DMA using URL parameters.
// The address is in hexadecimal (e.g. "D020"). The data is a hex string (e.g. "0504").
// Maximum 128 bytes can be written with this method.
func (c *Client) WriteMem(address string, data string) (*ErrorResponse, error) {
	params := url.Values{
		"address": {address},
		"data":    {data},
	}
	var resp ErrorResponse
	if err := c.put("/v1/machine:writemem", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WriteMemData writes binary data to C64 memory via DMA.
// The address is in hexadecimal (e.g. "0400"). The data must not wrap around $FFFF.
func (c *Client) WriteMemData(address string, data []byte) (*ErrorResponse, error) {
	params := url.Values{"address": {address}}
	var resp ErrorResponse
	if err := c.postBinary("/v1/machine:writemem", params, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReadMem reads data from C64 memory via DMA.
// The address is in hexadecimal (e.g. "0400"). Length specifies the number
// of bytes to read (default 256 if 0). Returns raw binary data.
func (c *Client) ReadMem(address string, length int) ([]byte, error) {
	params := url.Values{"address": {address}}
	if length > 0 {
		params.Set("length", fmt.Sprintf("%d", length))
	}
	return c.getRaw("/v1/machine:readmem", params)
}

// DebugRegResponse is returned by ReadDebugReg and WriteDebugReg.
type DebugRegResponse struct {
	Value  string   `json:"value"`
	Errors []string `json:"errors"`
}

// ReadDebugReg reads the debug register ($D7FF). U64 only.
func (c *Client) ReadDebugReg() (*DebugRegResponse, error) {
	var resp DebugRegResponse
	if err := c.get("/v1/machine:debugreg", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WriteDebugReg writes a value (hex) to the debug register ($D7FF) and returns
// the read-back value. U64 only.
func (c *Client) WriteDebugReg(value string) (*DebugRegResponse, error) {
	params := url.Values{"value": {value}}
	var resp DebugRegResponse
	if err := c.put("/v1/machine:debugreg", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
