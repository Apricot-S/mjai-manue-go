package main

import (
	"encoding/json/v2"
	"fmt"
	"slices"
)

type normalizedAction struct {
	Type      string   `json:"type"`
	Actor     *int     `json:"actor,omitempty"`
	Target    *int     `json:"target,omitempty"`
	Pai       string   `json:"pai,omitempty"`
	Consumed  []string `json:"consumed,omitempty"`
	Tsumogiri *bool    `json:"tsumogiri,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

func normalizeRawAction(raw []byte) (normalizedAction, bool, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return normalizedAction{}, false, err
	}
	t, _ := m["type"].(string)
	action := normalizedAction{Type: t}

	switch t {
	case "dahai":
		action.Actor = intPtrFromMap(m, "actor")
		action.Pai, _ = m["pai"].(string)
		action.Tsumogiri = boolPtrFromMap(m, "tsumogiri")
	case "reach":
		action.Actor = intPtrFromMap(m, "actor")
	case "chi", "pon", "daiminkan":
		action.Actor = intPtrFromMap(m, "actor")
		action.Target = intPtrFromMap(m, "target")
		action.Pai, _ = m["pai"].(string)
		action.Consumed = stringSliceFromMap(m, "consumed")
	case "ankan":
		action.Actor = intPtrFromMap(m, "actor")
		action.Consumed = stringSliceFromMap(m, "consumed")
	case "kakan":
		action.Actor = intPtrFromMap(m, "actor")
		action.Pai, _ = m["pai"].(string)
		action.Consumed = stringSliceFromMap(m, "consumed")
	case "hora":
		action.Actor = intPtrFromMap(m, "actor")
		action.Target = intPtrFromMap(m, "target")
		action.Pai, _ = m["pai"].(string)
	case "ryukyoku":
		action.Actor = intPtrFromMap(m, "actor")
		action.Reason, _ = m["reason"].(string)
		if action.Reason != "kyushukyuhai" {
			return normalizedAction{}, false, nil
		}
	case "none":
		action.Actor = intPtrFromMap(m, "actor")
	default:
		return normalizedAction{}, false, nil
	}
	return action, true, nil
}

func intPtrFromMap(m map[string]any, key string) *int {
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch n := v.(type) {
	case int:
		return &n
	case int64:
		i := int(n)
		return &i
	case float64:
		i := int(n)
		return &i
	}
	return nil
}

func boolPtrFromMap(m map[string]any, key string) *bool {
	v, ok := m[key].(bool)
	if !ok {
		return nil
	}
	return &v
}

func stringSliceFromMap(m map[string]any, key string) []string {
	values, ok := m[key].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		s, ok := value.(string)
		if !ok {
			continue
		}
		out = append(out, s)
	}
	return out
}

func actionsEqual(a, b normalizedAction) bool {
	if a.Type != b.Type || !intPtrEqual(a.Actor, b.Actor) || !intPtrEqual(a.Target, b.Target) {
		return false
	}
	if a.Pai != b.Pai || !boolPtrEqual(a.Tsumogiri, b.Tsumogiri) || a.Reason != b.Reason {
		return false
	}
	return slices.Equal(a.Consumed, b.Consumed)
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func mustActionJSON(a normalizedAction) string {
	b, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%+v", a)
	}
	return string(b)
}

func isZeroAction(a normalizedAction) bool {
	return a.Type == "" && a.Actor == nil && a.Target == nil && a.Pai == "" && len(a.Consumed) == 0 && a.Tsumogiri == nil && a.Reason == ""
}
