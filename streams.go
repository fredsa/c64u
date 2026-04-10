package main

import (
	"net/url"
)

// StreamName represents a valid data stream name.
type StreamName string

const (
	StreamVideo StreamName = "video"
	StreamAudio StreamName = "audio"
	StreamDebug StreamName = "debug"
)

// StartStream starts a data stream to the given IP address. U64 only.
// The ip parameter can include a port number (e.g. "192.168.1.100:6789").
// Default ports: video=11000, audio=11001, debug=11002.
// Note: turning on the video stream automatically turns off the debug stream.
func (c *Client) StartStream(stream StreamName, ip string) (*ErrorResponse, error) {
	params := url.Values{"ip": {ip}}
	var resp ErrorResponse
	if err := c.put("/v1/streams/"+url.PathEscape(string(stream))+":start", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StopStream stops a data stream. U64 only.
func (c *Client) StopStream(stream StreamName) (*ErrorResponse, error) {
	var resp ErrorResponse
	if err := c.put("/v1/streams/"+url.PathEscape(string(stream))+":stop", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
