package broker

import (
	"crypto/subtle"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"regexp"
)

var clientID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type accessHook struct {
	mqtt.HookBase
	axonToken string
	coreToken string
}

func (h *accessHook) ID() string { return "synapse-access" }
func (h *accessHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnACLCheck
}
func equalSecret(a, b string) bool {
	return a != "" && b != "" && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
func (h *accessHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	switch string(pk.Connect.Username) {
	case "synapse-core":
		return cl.ID == "synapse_core" && equalSecret(string(pk.Connect.Password), h.coreToken)
	case "axon":
		return clientID.MatchString(cl.ID) && cl.ID != "synapse_core" && equalSecret(string(pk.Connect.Password), h.axonToken)
	default:
		return false
	}
}
func (h *accessHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if string(cl.Properties.Username) == "synapse-core" && cl.ID == "synapse_core" {
		if write {
			return validTopic(topic, "synapse/v1/command/")
		}
		return topic == "synapse/v1/discovery/#" || validTopic(topic, "synapse/v1/discovery/")
	}
	if string(cl.Properties.Username) != "axon" {
		return false
	}
	if write {
		return topic == "synapse/v1/discovery/"+cl.ID
	}
	return topic == "synapse/v1/command/"+cl.ID
}
func validTopic(topic, prefix string) bool {
	if len(topic) <= len(prefix) || topic[:len(prefix)] != prefix {
		return false
	}
	return clientID.MatchString(topic[len(prefix):])
}
