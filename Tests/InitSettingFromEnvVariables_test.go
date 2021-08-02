package PlanningPoker_tests

import (
	"PlanningPoker/PlanningPokerSettings"
	"testing"
)

func TestInitSessionFromEnvVariables(t *testing.T) {
	var sessionsTool = PlanningPokerSettings.ServerSettings{}
	sessionsTool.InitSettingFromEnvVariables("docker")
	if sessionsTool.ServerHost.InternalHostName == "" {
		t.Error("Env variable didn't fill")
	}
}
