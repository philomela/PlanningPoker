package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net"
	"os"

	"./ServerPlanningPoker"
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	var currentServersSettings ServerPlanningPoker.ServersSettings
	var currentServer net.Listener
	var err error

	/* Открываем, читаем и парсим json */
	jsonFile, err := os.Open("serversSettings.json")

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	} else {
		defer jsonFile.Close()

		byteArrayJsonSettings, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteArrayJsonSettings, &currentServersSettings)
		fmt.Println("Server succsesful configured. ©Roman Solovyev/Andrew Zabolotniy")
	}

	/* Получаем экземпляры серверов */
	currentServer, err = ServerPlanningPoker.Server{ProtocolServer: currentServersSettings.ServerPlanningPoker.ProtocolServer, Port: currentServersSettings.ServerPlanningPoker.Port}.GetServer()

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	} else {
		fmt.Println("Server started and listening port: " + currentServersSettings.ServerPlanningPoker.Port)
	}

	currenSqlServer, err := ServerPlanningPoker.ServerSql{DSN: currentServersSettings.SQLServer.DSN, TypeSql: currentServersSettings.SQLServer.TypeSql}.OpenConnection()

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}

	/* Общаемся с базой */
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

	currenSqlServer.Exec("EXECUTE SaveError 'Error test'")
	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}

	currentServer.Accept()
	//currentServer.
}
