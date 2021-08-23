package PlanningPoker_tests

import (
	"PlanningPoker/PlanningPokerSettings"
	"testing"
)

func TestInitSettingFromLocalFile(t *testing.T) {
	var settings = PlanningPokerSettings.ServerSettings{}
	settings.InitSettingFromLocalFile("../serversSettings.json")
	if settings.ServerHost.InternalHostName != "http://localhost:8080/" {
		t.Error("Env variable didn't fill")
	}
}
