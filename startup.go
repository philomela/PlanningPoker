package main

import (
	"fmt"
	"net"

	"./ServerPlanningPoker"
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	var currentServer net.Listener
	var err error
	currentServer, err = ServerPlanningPoker.Server{ProtocolServer: "tcp", Port: ":4545"}.GetServer()

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}

	//currenSqlServer, err := ServerPlanningPoker.ServerSql{DSN: "sqlserver://u0932131_admin:RomanAndrey46@31.31.196.202/instance?database=u0932131_planningPoker", TypeSql: "mssql"}.OpenConnection()
	currenSqlServer, err := ServerPlanningPoker.ServerSql{DSN: "server=31.31.196.202;user id=u0932131_admin;password=RomanAndrey46;database=u0932131_planningPoker", TypeSql: "mssql"}.OpenConnection()

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}

	result, err := currenSqlServer.Query("SELECT * FROM test") //Делаем запрос к тестовой базе
	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}
	collection, err := result.Columns()
	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}

	fmt.Printf(collection[0]) //Выводим в сообщении, первое название поля из таблицы Test
	currentServer.Accept()
}
