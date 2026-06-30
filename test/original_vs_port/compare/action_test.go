package main

import "testing"

func TestNormalizeRawAction_IgnoresLog(t *testing.T) {
	got, comparable, err := normalizeRawAction([]byte(`{"type":"dahai","actor":1,"pai":"5m","tsumogiri":false,"log":"discard"}`))
	if err != nil {
		t.Fatalf("normalizeRawAction() failed: %v", err)
	}
	if !comparable {
		t.Fatal("normalizeRawAction() comparable = false, want true")
	}
	want := normalizedAction{Type: "dahai", Actor: new(1), Pai: "5m", Tsumogiri: new(false)}
	if !actionsEqual(got, want) {
		t.Errorf("action = %+v, want %+v", got, want)
	}
}

func TestNormalizeRawAction_RepresentativeActions(t *testing.T) {
	tests := []string{
		`{"type":"reach","actor":1}`,
		`{"type":"chi","actor":1,"target":0,"pai":"3m","consumed":["1m","2m"]}`,
		`{"type":"pon","actor":1,"target":0,"pai":"3m","consumed":["3m","3m"]}`,
		`{"type":"hora","actor":1,"target":0,"pai":"3m"}`,
		`{"type":"none","actor":1}`,
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			action, comparable, err := normalizeRawAction([]byte(input))
			if err != nil {
				t.Fatalf("normalizeRawAction() failed: %v", err)
			}
			if !comparable {
				t.Fatal("normalizeRawAction() comparable = false, want true")
			}
			if action.Type == "" {
				t.Error("action type is empty")
			}
		})
	}
}
