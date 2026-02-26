package tools

import (
	"encoding/json"
	"fmt"
)

const profileDurationPruneRatio = 0.01

type flameGraphNode struct {
	Name     string            `json:"name"`
	Start    int64             `json:"start"`
	Duration int64             `json:"duration"`
	Children []*flameGraphNode `json:"children"`
	Color    *string           `json:"color,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

func trimProfilesResponse(response []byte) ([]byte, error) {
	var root flameGraphNode
	if err := json.Unmarshal(response, &root); err != nil {
		return nil, fmt.Errorf("error unmarshaling profiles response: %v", err)
	}

	changed := pruneFlameGraphNodesBelowDurationThreshold(&root)
	if !changed {
		return response, nil
	}

	trimmed, err := json.Marshal(root)
	if err != nil {
		return nil, fmt.Errorf("error marshaling trimmed profiles response: %v", err)
	}

	return trimmed, nil
}

func pruneFlameGraphNodesBelowDurationThreshold(root *flameGraphNode) bool {
	if root == nil || root.Duration <= 0 {
		return false
	}

	minDuration := float64(root.Duration) * profileDurationPruneRatio
	return pruneFlameGraphChildren(root, minDuration)
}

func pruneFlameGraphChildren(node *flameGraphNode, minDuration float64) bool {
	if node == nil || len(node.Children) == 0 {
		return false
	}

	changed := false
	filteredChildren := node.Children[:0]
	for _, child := range node.Children {
		if child == nil || float64(child.Duration) < minDuration {
			changed = true
			continue
		}

		if pruneFlameGraphChildren(child, minDuration) {
			changed = true
		}
		filteredChildren = append(filteredChildren, child)
	}

	node.Children = filteredChildren
	return changed
}
