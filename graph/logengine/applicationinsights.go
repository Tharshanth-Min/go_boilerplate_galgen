package logengine

import (
	"os"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var telemetryClient appinsights.TelemetryClient

func CreateTelemetryClient() {
	instrumentationKey := os.Getenv("APPLICATION_INSIGHTS_KEY")
	telemetryConfig := appinsights.NewTelemetryConfiguration(instrumentationKey)

	/*// Configure how many items can be sent in one call to the data collector:
	telemetryConfig.MaxBatchSize = 8192

	// Configure the maximum delay before sending queued telemetry:
	telemetryConfig.MaxBatchInterval = 2 * time.Second

	// diagnostic command
	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		log.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})*/

	client := appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	telemetryClient = client
	telemetryClient.Context().Tags.Cloud().SetRole("Download")
}

func GetTelemetryClient() appinsights.TelemetryClient {
	if telemetryClient == nil {
		CreateTelemetryClient()
	}
	return telemetryClient
}
