package ServerPlanningPoker

type ServersSettings struct {
	SQLServer  ServerSql
	ServerHost ServerHost
}

type ServerHost struct {
	Host string
	Room string
}
