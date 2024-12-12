package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type GetK8sEventAttributeValueHandlerArgs struct {
	TimeConfig     utils.TimeConfig    `json:"time_config" jsonschema:"required,description=The time period to get the possible values of K8 event attributes values. e.g. if you want to see the possible values for the attributes in the last 5 minutes you would set time_period=5 and time_window=Minutes. You can also set an absoulute time range by setting start_time and end_time"`
	Attribute      string              `json:"attribute" jsonschema:"required,description=The attribute key to get the possible values for"`
	Filters        map[string][]string `json:"filters" jsonschema:"description=The filters to apply before getting the possible values. For example if you want to get the possible values for attribute key service.name where the environment is X you would set the Filters as {environment: [X]}"`
	ExcludeFilters map[string][]string `json:"excludeFilters" jsonschema:"description=The exclude filters to exclude/eliminate possible values an attribute can take. Events matching the exclude filters will not be returned. For example if you want the possible values for attribute key service.name where the attribute environment is not X then you would set the ExcludeFilters as {environment: [X]}"`
	Regexes        []string            `json:"regexes" jsonschema:"description=The regexes to apply to the event messages. Only the attribute values (for a given attribute key) of events messages that match these regexes will be returned. For example if you want the possible values for attribute key service.name where the event message contains the word 'error' you would set the regexes as ['error']"`
	ExcludeRegexes []string            `json:"excludeRegexes" jsonschema:"description=The exclude regexes to apply to the event messages. The attribute values (for a given attribute key) of events messages that match these regexes will not be returned. For example if you want the possible values for attribute key service.name where the event message does not contain the word 'error' you would set the exclude regexes as ['error']"`
	Environments   []string            `json:"environments" jsonschema:"description=The environments to get events from. If empty, events from all environments will be returned"`
	Ascending      bool                `json:"ascending" jsonschema:"description=If true, events will be returned in ascending order, otherwise in descending order"`
}

func GetK8sEventAttributeValuesForIndividualAttributeHandler(arguments GetK8sEventAttributeValueHandlerArgs) (*mcpgolang.ToolResponse, error) {
	startTime, endTime, err := utils.CalculateTimeRange(arguments.TimeConfig)
	if err != nil {
		return nil, fmt.Errorf("error calculating time range: %v", err)
	}
	request := model.GetSingleK8sEventSummaryRequest{
		GetK8sEventsRequest: model.GetK8sEventsRequest{
			StartTime:      startTime,
			EndTime:        endTime,
			Filters:        arguments.Filters,
			ExcludeFilters: arguments.ExcludeFilters,
			Regexes:        arguments.Regexes,
			ExcludeRegexes: arguments.ExcludeRegexes,
			Environments:   arguments.Environments,
			Ascending:      arguments.Ascending,
		},
		Attribute: arguments.Attribute,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := utils.MakeMetoroAPIRequest("POST", "k8s/events/summaryIndividualAttribute", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error making Metoro call: %v", err)
	}

	return mcpgolang.NewToolReponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}
