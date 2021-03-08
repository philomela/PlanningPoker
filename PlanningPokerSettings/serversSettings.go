package PlanningPokerSettings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ServerSettings struct {
	SQLServer  ServerSql
	ServerHost ServerHost
	SmtpServer ServerSmtp
}

type ServerHost struct {
	Host                     string
	Port                     string
	InternalHostName         string
	ExternalHostName         string
	Room                     string
	WebSocketProtocol        string
	WebSocketExternalAddress string
}

type ServerSmtp struct {
	Host      string
	LoginHost string
	PassHost  string
	PortHost  string
}

func InitServerSetting() *ServerSettings {
	currentServerSettings := &ServerSettings{}
	currentServerSettings.InitSettingFromLocalFile()
	currentServerSettings.InitSettingFromEnvVariables(os.Getenv("APP_ENV_NAME_PLANNING_POKER"))
	return currentServerSettings
}

func (s *ServerSettings) InitSettingFromLocalFile() {
	jsonFile, err := os.Open("./serversSettings.json")

	defer jsonFile.Close()

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	} else {
		byteArrayJsonSettings, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteArrayJsonSettings, s)
	}
}

func (s *ServerSettings) InitSettingFromEnvVariables(envParam string) {
	switch envParam {
	case "docker":
		s.ServerHost.Host = os.Getenv("PLANNING_POKER_PROTOCOL")
		s.ServerHost.Host = os.Getenv("PLANNING_POKER_HOST")
		s.ServerHost.Port = os.Getenv("PLANNING_POKER_PORT")
		s.ServerHost.InternalHostName = os.Getenv("PLANNING_POKER_INTERNALHOSTNAME")
		s.ServerHost.ExternalHostName = os.Getenv("PLANNING_POKER_EXTERNALHOSTNAME")
		s.ServerHost.WebSocketProtocol = os.Getenv("PLANNING_POKER_WEBSOCKET_PROTOCOL")
		s.ServerHost.WebSocketExternalAddress = os.Getenv("PLANNING_POKER_WEBSOCKET_EXTERNALADDRESS")
	}
}
