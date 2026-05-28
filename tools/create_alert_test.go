package tools

import (
	"encoding/json"
	"testing"

	"github.com/metoro-io/metoro-mcp-server/model"
)

func TestCreateAlertHandlerArgsUnmarshalsInvestigateOnFire(t *testing.T) {
	var args CreateAlertHandlerArgs
	if err := json.Unmarshal([]byte(`{"investigate_on_fire":true}`), &args); err != nil {
		t.Fatalf("failed to unmarshal args: %v", err)
	}

	if !args.InvestigateOnFire {
		t.Fatal("expected investigate_on_fire to unmarshal to true")
	}
}

func TestCreateAlertPayloadUsesTopLevelInvestigateOnFire(t *testing.T) {
	alert := model.Alert{
		Metadata: model.MetadataObject{
			Name: "High 5XX Rate",
			Id:   "high-5xx-rate",
		},
		Timeseries: model.TimeseriesConfig{},
	}
	alert.SetInvestigateOnFire(true)

	request := model.CreateUpdateAlertRequest{Alert: alert}
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal request JSON: %v", err)
	}
	alertPayload, ok := payload["alert"].(map[string]any)
	if !ok {
		t.Fatalf("expected alert payload object, got %T", payload["alert"])
	}

	if got := alertPayload["investigateOnFire"]; got != true {
		t.Fatalf("expected top-level alert.investigateOnFire=true, got %#v", got)
	}
	if metadataPayload, ok := alertPayload["metadata"].(map[string]any); ok {
		if _, exists := metadataPayload["investigateOnFire"]; exists {
			t.Fatal("did not expect investigateOnFire under alert.metadata")
		}
	}
	if timeseriesPayload, ok := alertPayload["timeseries"].(map[string]any); ok {
		if _, exists := timeseriesPayload["investigateOnFire"]; exists {
			t.Fatal("did not expect investigateOnFire under alert.timeseries")
		}
	}
}

func TestNewAlertOmitsInvestigateOnFireByDefault(t *testing.T) {
	alert := model.NewAlert(model.MetadataObject{
		Name: "High 5XX Rate",
		Id:   "high-5xx-rate",
	}, model.TimeseriesConfig{})

	request := model.CreateUpdateAlertRequest{Alert: *alert}
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal request JSON: %v", err)
	}
	alertPayload, ok := payload["alert"].(map[string]any)
	if !ok {
		t.Fatalf("expected alert payload object, got %T", payload["alert"])
	}
	if _, exists := alertPayload["investigateOnFire"]; exists {
		t.Fatal("did not expect investigateOnFire to be marshaled by default")
	}
}
