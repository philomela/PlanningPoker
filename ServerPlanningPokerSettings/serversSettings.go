package ServerPlanningPoker

type ServersSettings struct {
	SQLServer  ServerSql
	ServerHost ServerHost
	SmtpServer ServerSmtp
}

type ServerHost struct {
	Host string
	Room string
}

type ServerSmtp struct {
	Host      string
	LoginHost string
	PassHost  string
	PortHost  string
}
