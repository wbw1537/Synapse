package notification

import (
	"fmt"
	"sync"
	"time"
)

type AlertState struct {
	LastStatus    string
	LastAlertTime time.Time
}

type AlertManager struct {
	sender Sender
	states map[string]*AlertState // Key: "serviceID:widgetIndex:monitorIndex"
	mu     sync.Mutex
}

func NewAlertManager(sender Sender) *AlertManager {
	return &AlertManager{
		sender: sender,
		states: make(map[string]*AlertState),
	}
}

// CheckAndAlert evaluates if a notification should be sent based on state change
func (am *AlertManager) CheckAndAlert(key string, isTriggered bool, severity, message, serviceName string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	state, exists := am.states[key]
	if !exists {
		state = &AlertState{LastStatus: "ok"}
		am.states[key] = state
	}

	currentStatus := "ok"
	if isTriggered {
		currentStatus = severity
	}

	// Logic: Alert only on state change
	if state.LastStatus != currentStatus {
		if isTriggered {
			// Resolved -> Error/Warning
			subject := fmt.Sprintf("%s: %s - %s", severity, serviceName, message)
			body := fmt.Sprintf("Service: %s\nAlert: %s\nSeverity: %s\nTime: %s", serviceName, message, severity, time.Now().Format(time.RFC1123))
			go am.sender.Send(subject, body)
		} else {
			// Error/Warning -> Resolved
			// Optional: Send "Resolved" email? For MVP, let's skip to reduce noise, or enable if requested.
			// Let's print log.
			fmt.Printf("Alert Resolved: %s - %s\n", serviceName, message)
		}
		state.LastStatus = currentStatus
		state.LastAlertTime = time.Now()
	}
}
