package tools

import (
	"encoding/json"
	"testing"
)

func TestTrimProfilesResponsePrunesNodesBelowOnePercent(t *testing.T) {
	root := flameGraphNode{
		Name:     "root",
		Duration: 1000,
		Children: []*flameGraphNode{
			{
				Name:     "drop-nine",
				Duration: 9,
				Children: []*flameGraphNode{
					{
						Name:     "drop-nine-child",
						Duration: 100,
					},
				},
			},
			{
				Name:     "keep-ten",
				Duration: 10,
			},
			{
				Name:     "keep-large",
				Duration: 300,
				Children: []*flameGraphNode{
					{
						Name:     "drop-five",
						Duration: 5,
					},
					{
						Name:     "keep-twenty",
						Duration: 20,
						Children: []*flameGraphNode{
							{
								Name:     "drop-one",
								Duration: 1,
							},
							{
								Name:     "keep-ten-nested",
								Duration: 10,
							},
						},
					},
				},
			},
		},
	}

	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("failed to marshal test payload: %v", err)
	}

	trimmed, err := trimProfilesResponse(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var parsed flameGraphNode
	if err := json.Unmarshal(trimmed, &parsed); err != nil {
		t.Fatalf("failed to unmarshal trimmed response: %v", err)
	}

	if findFlameGraphChild(&parsed, "drop-nine") != nil {
		t.Fatalf("expected node below 1%% to be pruned")
	}
	if findFlameGraphChild(&parsed, "keep-ten") == nil {
		t.Fatalf("expected node at exactly 1%% to be retained")
	}

	keepLarge := findFlameGraphChild(&parsed, "keep-large")
	if keepLarge == nil {
		t.Fatalf("expected keep-large to be retained")
	}
	if findFlameGraphChild(keepLarge, "drop-five") != nil {
		t.Fatalf("expected nested node below threshold to be pruned")
	}

	keepTwenty := findFlameGraphChild(keepLarge, "keep-twenty")
	if keepTwenty == nil {
		t.Fatalf("expected keep-twenty to be retained")
	}
	if findFlameGraphChild(keepTwenty, "drop-one") != nil {
		t.Fatalf("expected deep nested node below threshold to be pruned")
	}
	if findFlameGraphChild(keepTwenty, "keep-ten-nested") == nil {
		t.Fatalf("expected deep nested node at threshold to be retained")
	}
}

func TestTrimProfilesResponseLeavesZeroDurationRootUnchanged(t *testing.T) {
	root := flameGraphNode{
		Name:     "root",
		Duration: 0,
		Children: []*flameGraphNode{
			{
				Name:     "child",
				Duration: 1,
			},
		},
	}

	raw, err := json.Marshal(root)
	if err != nil {
		t.Fatalf("failed to marshal test payload: %v", err)
	}

	trimmed, err := trimProfilesResponse(raw)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(trimmed) != string(raw) {
		t.Fatalf("expected zero-duration root payload to remain unchanged")
	}
}

func TestTrimProfilesResponseReturnsErrorForInvalidJSON(t *testing.T) {
	_, err := trimProfilesResponse([]byte(`{"name":`))
	if err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

func findFlameGraphChild(node *flameGraphNode, name string) *flameGraphNode {
	if node == nil {
		return nil
	}

	for _, child := range node.Children {
		if child != nil && child.Name == name {
			return child
		}
	}

	return nil
}
