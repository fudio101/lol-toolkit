package lcu

import (
	lcuclient "github.com/its-haze/lcu-gopher"
)

// WebSocketClient wraps lcu-gopher's WebSocket functionality.
// Use the Client's GetLCUClient() method to access lcu-gopher's Subscribe methods directly.
type WebSocketClient struct {
	client *lcuclient.Client
}

// NewWebSocketClient creates a new WebSocket client using lcu-gopher.
// Note: The underlying lcu-gopher client already has WebSocket support built-in.
// Use Client.GetLCUClient() to access Subscribe methods directly.
func (c *Client) NewWebSocketClient() (*WebSocketClient, error) {
	return &WebSocketClient{
		client: c.client,
	}, nil
}

// Subscribe subscribes to an event using lcu-gopher's event system.
func (ws *WebSocketClient) Subscribe(endpoint string, handler func(*lcuclient.Event), eventTypes ...lcuclient.EventType) error {
	return ws.client.Subscribe(endpoint, handler, eventTypes...)
}

// SubscribeToAll subscribes to all events.
func (ws *WebSocketClient) SubscribeToAll(handler func(*lcuclient.Event)) error {
	return ws.client.SubscribeToAll(handler)
}

// Start is a no-op since lcu-gopher handles WebSocket automatically after Connect().
func (ws *WebSocketClient) Start() error {
	// WebSocket is already connected when Client.Connect() is called
	return nil
}

// Stop disconnects the WebSocket (disconnects the entire client).
func (ws *WebSocketClient) Stop() {
	ws.client.Disconnect()
}
