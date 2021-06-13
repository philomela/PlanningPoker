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
	Protocol                       string
	Host                           string
	Port                           string
	InternalHostName               string
	ExternalHostName               string
	ExternalPathToLoginForm        string
	ExternalPathToNewRoom          string
	ExternalPathToRestoreAcc       string
	ExternalPathToChangePass       string
	Room                           string
	WebSocketProtocol              string
	WebSocketExternalAddress       string
	RestoreAccount                 string
	ExternalPathToRegistrationForm string
}

type ServerSmtp struct {
	Host      string
	LoginHost string
	PassHost  string
	PortHost  string
	ApiKey    string
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
		return
	} else {
		byteArrayJsonSettings, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteArrayJsonSettings, s)
	}
}

func (s *ServerSettings) InitSettingFromEnvVariables(envParam string) {
	switch envParam {
	case "docker":
		s.ServerHost.Protocol = os.Getenv("PLANNING_POKER_PROTOCOL")
		s.ServerHost.Host = os.Getenv("PLANNING_POKER_HOST")
		s.ServerHost.Port = os.Getenv("PLANNING_POKER_PORT")
		s.ServerHost.InternalHostName = os.Getenv("PLANNING_POKER_INTERNALHOSTNAME")
		s.ServerHost.ExternalHostName = os.Getenv("PLANNING_POKER_EXTERNALHOSTNAME")
		s.ServerHost.ExternalPathToLoginForm = os.Getenv("PLANNING_POKER_EXT_PATH_LOGINFROM")
		s.ServerHost.ExternalPathToNewRoom = os.Getenv("PLANNING_POKER_EXT_PATH_NEWROOM")
		s.ServerHost.ExternalPathToRegistrationForm = os.Getenv("PLANNING_POKER_EXT_PATH_REGFORM")
		s.ServerHost.ExternalPathToRestoreAcc = os.Getenv("PLANNING_POKER_EXT_PATH_RESTOREACC")
		s.ServerHost.ExternalPathToChangePass = os.Getenv("PLANNING_POKER_EXT_PATH_CHANGEPASS")
		s.ServerHost.WebSocketProtocol = os.Getenv("PLANNING_POKER_WEBSOCKET_PROTOCOL")
		s.ServerHost.WebSocketExternalAddress = os.Getenv("PLANNING_POKER_WEBSOCKET_EXTERNALADDRESS")
		s.SmtpServer.ApiKey = os.Getenv("PLANNING_POKER_SMTP_APIKEY")
	}
}
