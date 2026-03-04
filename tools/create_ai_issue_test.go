package tools

import (
	"context"
	"strings"
	"testing"
)

func TestCreateAIIssueHandlerRejectsInvalidPriority(t *testing.T) {
	priority := "P0"
	_, err := CreateAIIssueHandler(context.Background(), CreateAIIssueHandlerArgs{
		Title:       "title",
		Description: "description",
		Summary:     "summary",
		Priority:    &priority,
	})
	if err == nil {
		t.Fatalf("expected error for invalid priority")
	}
	if !strings.Contains(err.Error(), "invalid priority") {
		t.Fatalf("expected invalid priority error, got: %v", err)
	}
}

func TestCreateAIIssueHandlerRejectsInvalidCategory(t *testing.T) {
	category := "network"
	_, err := CreateAIIssueHandler(context.Background(), CreateAIIssueHandlerArgs{
		Title:       "title",
		Description: "description",
		Summary:     "summary",
		Category:    &category,
	})
	if err == nil {
		t.Fatalf("expected error for invalid category")
	}
	if !strings.Contains(err.Error(), "invalid category") {
		t.Fatalf("expected invalid category error, got: %v", err)
	}
}
