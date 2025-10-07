package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/metoro-mcp-server/model"
	"github.com/metoro-io/metoro-mcp-server/utils"
)

type CreateAlertHandlerArgs struct {
	AlertName         string                          `json:"alert_name" jsonschema:"required,description=The name of the alert to create"`
	AlertDescription  string                          `json:"alert_description" jsonschema:"required,description=The description of the alert to create"`
	Timeseries        []model.SingleTimeseriesRequest `json:"timeseries" jsonschema:"required,description=Array of timeseries data to get. Each item in this array corresponds to a single timeseries. You can then use the formulas to combine these timeseries. If you only want to see the combination of timeseries via defining formulas and if you dont want to see the individual timeseries data when setting formulas you can set shouldNotReturn to true"`
	Formula           model.Formula                   `json:"formula" jsonschema:"description=Optional formula to combine timeseries. Formula should only consist of formulaIdentifier of the timeseries in the timeseries array. e.g. a + b + c if a b c appears in the formulaIdentifier of the timeseries array. You can ONLY do the following operations: Arithmetic operations:+ (for add) - (for substract) * (for multiply) / (for division) % (for modulus) ^ or ** (for exponent). Comparison: == != < > <= >= . Logical:! (for not) && (for AND) || (for OR). Conditional operations: ?: (ternary) e.g. (a || b) ? 1 : 0. Do not guess the operations. Just use these available ones!"`
	Condition         string                          `json:"condition" jsonschema:"required,enum=GreaterThan,enum=LessThan,enum=GreaterThanOrEqual,enum=LessThanOrEqual,description=the arithmetic comparison to use to evaluate whether an alert is firing or not. This is used to determine whether the alert should be triggered based on the threshold value."`
	Threshold         float64                         `json:"threshold" jsonschema:"required,description=The threshold value for the alert. This is the value that will be used together with the the arithmetic condition to see whether the alert should be triggered or not. For example if you set the condition to GreaterThan and the threshold to 100 then the alert will fire if the value of the timeseries is greater than 100."`
	DatapointsToAlarm int64                           `json:"datapoints_to_alarm" jsonschema:"required,description=The number of datapoints that need to breach the threshold for the alert to be triggered"`
	EvaluationWindow  int64                           `json:"evaluation_window" jsonschema:"required,description=The evaluation window in number of datapoints. This is the number of datapoints that will be considered for evaluating the alert condition. For example if you set this to then the last 5 datapoints will be considered for evaluating the alert condition. This is useful for smoothing out spikes in the data and preventing false positives."`
}

func CreateAlertHandler(ctx context.Context, arguments CreateAlertHandlerArgs) (*mcpgolang.ToolResponse, error) {
	alert, err := createAlertFromTimeseries(ctx, arguments.AlertName, arguments.AlertDescription, arguments.Timeseries, arguments.Formula, arguments.Condition, arguments.Threshold, arguments.DatapointsToAlarm, arguments.EvaluationWindow)
	if err != nil {
		return nil, fmt.Errorf("error creating alert properties: %v", err)
	}

	newAlertRequest := model.SetAlertRequest{
		Alert: alert,
	}

	resp, err := setAlertMetoroCall(ctx, newAlertRequest)
	if err != nil {
		return nil, fmt.Errorf("error setting dashboard: %v", err)
	}
	return mcpgolang.NewToolResponse(mcpgolang.NewTextContent(fmt.Sprintf("%s", string(resp)))), nil
}

func createAlertFromTimeseries(ctx context.Context, alertName, alertDescription string, timeseries []model.SingleTimeseriesRequest, formula model.Formula, condition model.StaticAlarmCondition, threshold float64, datapointsToAlarm int64, evaluationWindow int64) (model.Alert, error) {
	alert := model.Alert{
		UUID:        uuid.NewString(),
		Name:        alertName,
		Description: alertDescription,
		MultiMetricAlert: &model.MultiMetricAlert{
			Filters: model.MultiMetricFilters{
				MetricSpecifiers: timeseries,
				Formula:          &formula,
			},
			MonitorEvaluation: model.MonitorEvaluation{
				MonitorEvaluationType: model.MetricMonitorEventStaticThreshold,
				MonitorEvalutionPayload: model.MonitorEvaluationPayload{
					StaticMonitorEvaluationPayload: model.StaticMonitorEvaluationPayload{
						DatapointsToAlarm:        datapointsToAlarm,
						EvaluationWindow:         evaluationWindow,
						MissingDatapointBehavior: model.MissingDatapointNotBreaching,
					},
				},
			},
			AlarmCondition: model.AlarmCondition{
				Condition: condition,
				Threshold: threshold,
			},
		},
	}
	return alert, nil
}

func setAlertMetoroCall(ctx context.Context, request model.SetAlertRequest) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling alert request: %v", err)
	}
	return utils.MakeMetoroAPIRequest("POST", "alert", bytes.NewBuffer(requestBody), utils.GetAPIRequirementsFromRequest(ctx))
}
