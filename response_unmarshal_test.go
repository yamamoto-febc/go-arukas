package arukas

import (
	"encoding/json"
	"testing"
)

func assertPort(t *testing.T, port *Port, number int32, protocol string) {
	if port.Protocol != protocol || port.Number != number {
		t.Errorf("Expected %d/%s but got %d/%s", number, protocol, port.Number, port.Protocol)
	}
}

func TestUnmarshalOldPortFormat(t *testing.T) {
	responseBody := `
		{
			"data": {
				"type": "services",
				"id": "uuid-1",
				"attributes": {
					"image": "my-factorio",
					"ports": [
						{"protocol": "tcp", "number": 80},
						{"protocol": "tcp", "number": 443},
						{"protocol": "udp", "number": 34197}
					]
				}
			}
		}
	`

	service := new(ServiceData)
	if err := json.Unmarshal([]byte(responseBody), service); err != nil {
		t.Fatal("Failed to unmarshal old format port:", err)
	}

	assertPort(t, service.Ports()[0], 80, "tcp")
	assertPort(t, service.Ports()[1], 443, "tcp")
	assertPort(t, service.Ports()[2], 34197, "udp")
}

func TestUnmarshalNewPortFormat(t *testing.T) {
	responseBody := `
		{
			"data": {
				"type": "services",
				"id": "uuid-1",
				"attributes": {
					"image": "my-factorio",
					"ports": ["80", "443/tcp", "34197/udp"]
				}
			}
		}
	`

	service := new(ServiceData)
	if err := json.Unmarshal([]byte(responseBody), service); err != nil {
		t.Fatal("Failed to unmarshal old format port:", err)
	}

	assertPort(t, service.Ports()[0], 80, "tcp")
	assertPort(t, service.Ports()[1], 443, "tcp")
	assertPort(t, service.Ports()[2], 34197, "udp")
}
