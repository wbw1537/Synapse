package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
)

const DiscoveryVersion = "v1"

var protocolID = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// DecodeDiscovery rejects unsupported shapes before they can silently lose data.
func DecodeDiscovery(payload []byte) (*ServicePayload, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil || fields == nil {
		return nil, fmt.Errorf("discovery must be a JSON object")
	}
	for _, field := range []string{"widgets", "actions"} {
		if _, ok := fields[field]; ok {
			return nil, fmt.Errorf("legacy %s is unsupported; use layout/components with api_version v1", field)
		}
	}
	var p ServicePayload
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&p); err != nil {
		return nil, fmt.Errorf("invalid discovery: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("discovery must contain exactly one JSON object")
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

// Validate checks the supported wire version and capability references.
func (p *ServicePayload) Validate() error {
	if p.APIVersion != DiscoveryVersion {
		return fmt.Errorf("api_version must be %q", DiscoveryVersion)
	}
	if !protocolID.MatchString(p.ID) {
		return fmt.Errorf("service id must contain only letters, digits, underscores or hyphens")
	}
	if p.TTL <= 0 {
		return fmt.Errorf("ttl must be positive")
	}
	switch p.Status {
	case "online", "warning", "error", "offline":
	default:
		return fmt.Errorf("status must be online, warning, error or offline")
	}
	if p.Layout.Type != "sections" || p.Layout.Root == nil {
		return fmt.Errorf("layout requires type sections and a root array")
	}
	if p.Components == nil {
		return fmt.Errorf("components must be an object (empty is allowed)")
	}
	actions := make(map[string]bool)
	for id, comp := range p.Components {
		if !protocolID.MatchString(id) || comp.ID != id {
			return fmt.Errorf("component %q requires an identical id using letters, digits, underscores or hyphens", id)
		}
		switch comp.Type {
		case "stat", "gauge", "status_indicator", "log_stream", "action_group", "link":
		default:
			return fmt.Errorf("component %q has unsupported type %q", id, comp.Type)
		}
		if comp.ActionID != "" || (comp.Type != "action_group" && comp.Items != nil) {
			return fmt.Errorf("component %q must declare actions in action_group.items", id)
		}
		if comp.Type == "action_group" {
			if comp.Items == nil {
				return fmt.Errorf("action_group %q requires an items array", id)
			}
			for _, item := range comp.Items {
				if !protocolID.MatchString(item.ActionID) || actions[item.ActionID] {
					return fmt.Errorf("action_id %q must be valid and unique within the service", item.ActionID)
				}
				actions[item.ActionID] = true
			}
		}
	}
	for i, section := range p.Layout.Root {
		if section.Type != "section" || section.Children == nil {
			return fmt.Errorf("layout.root[%d] requires type section and a children array", i)
		}
		for _, id := range section.Children {
			if _, ok := p.Components[id]; !ok {
				return fmt.Errorf("layout.root[%d] references unknown component %q", i, id)
			}
		}
	}
	return nil
}
