package tools

import (
	"context"
	"fmt"
	"time"

	mcpgolang "github.com/metoro-io/mcp-golang"
)

type UnixToRFC3339HandlerArgs struct {
	UnixTimestamp int64 `json:"unix_timestamp" jsonschema:"required,description=Unix timestamp in seconds or milliseconds"`
}

func UnixToRFC3339Handler(ctx context.Context, arguments UnixToRFC3339HandlerArgs) (*mcpgolang.ToolResponse, error) {
	// Determine if the timestamp is in seconds or milliseconds
	// Unix timestamps in seconds are typically 10 digits (until year 2286)
	// Unix timestamps in milliseconds are typically 13 digits
	var t time.Time

	// Check if it's likely milliseconds (more than 10 digits or would result in a date far in the future)
	if arguments.UnixTimestamp > 9999999999 {
		// Treat as milliseconds
		t = time.Unix(0, arguments.UnixTimestamp*int64(time.Millisecond))
	} else {
		// Treat as seconds
		t = time.Unix(arguments.UnixTimestamp, 0)
	}

	// Convert to RFC3339 format
	rfc3339String := t.UTC().Format(time.RFC3339)

	// Create a response with both interpretations if the timestamp could be ambiguous
	var response string
	if arguments.UnixTimestamp <= 9999999999 && arguments.UnixTimestamp >= 1000000000 {
		// Could be either seconds or milliseconds, show both
		tAsMillis := time.Unix(0, arguments.UnixTimestamp*int64(time.Millisecond))
		response = fmt.Sprintf("Interpreted as seconds: %s\nInterpreted as milliseconds: %s",
			rfc3339String,
			tAsMillis.UTC().Format(time.RFC3339))
	} else {
		response = rfc3339String
	}

	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(response)), nil
}
