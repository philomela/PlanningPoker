package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"./ServerPlanningPoker"
	"./session"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	cookieName = "sessionId"
)

func main() {
	var currentServersSettings ServerPlanningPoker.ServersSettings

	//var err error
	type connection struct {
		conn    *websocket.Conn
		rooomId string
		UUID    uuid.UUID
	}
	var conns []connection
	type connections []*websocket.Conn
	//var rooms []Room //Коллекция комнат, вынести в отдельный пакет
	//collectionRoom := make(map[string]connections)
	/* Открываем, читаем и парсим json */
	jsonFile, err := os.Open("serversSettings.json")

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	} else {
		defer jsonFile.Close()

		byteArrayJsonSettings, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteArrayJsonSettings, &currentServersSettings)
	}

	currentSqlServer, err := ServerPlanningPoker.ServerSql{DSN: currentServersSettings.SQLServer.DSN, TypeSql: currentServersSettings.SQLServer.TypeSql}.OpenConnection()

	fmt.Println("Server succsesful configured. ©Roman Solovyev")

	http.HandleFunc("/echo/", func(w http.ResponseWriter, r *http.Request) {
		var URL = r.URL.Query()["roomId"][0]
		var Conn connection
		Conn.conn, _ = upgrader.Upgrade(w, r, nil)
		Conn.rooomId = URL
		Conn.UUID, err = uuid.NewUUID()
		if err != nil {
			fmt.Println(err)
			return //Добавить реализацию логирования ошибок в базу
		}
		currentSqlServer.Query(`EXEC [CreateConnection] @UUID=?, @RoomId=?`, Conn.UUID.String(), Conn.rooomId)
		conns = append(conns, Conn)

		//collectionRoom[URL] = conns
		//conn.WriteMessage(1, []byte(URL))
		//conn.WriteMessage(1, []byte("Hello, I'm TCP server"))
		println(conns)
		println(URL)
		for {
			msgType, msg, err := Conn.conn.ReadMessage()
			//Добавить логику работы когда мы закрываем соединение, чтобы удалять ненужные
			if err != nil {
				return
			}

			currentSqlServer.Exec("INSERT INTO ServerPlanningPoker.Connections(Id, [GUID], RoomId) VALUES(1, NEWID(), 1)")
			if err != nil {
				fmt.Println(err)
				return //Добавить реализацию логирования ошибок в базу
			}
			currentSqlServer.Exec("EXECUTE SaveError 'Error test'")
			if err != nil {
				fmt.Println(err)
				return //Добавить реализацию логирования ошибок в базу
			}

			for _, value := range conns {
				if value.rooomId == URL {
					msg = []byte(URL)
					value.conn.WriteMessage(msgType, msg)
				}

			}

			//return
		}
	})

	http.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query()["roomId"][0]
		if len(roomId) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			//http.Error(w, "Bad Request", http.StatusBadRequest)
			//http.ServeFile(w, r, "templates/error_bad_request.html") //реализовать вызов методов страниц ошибок
			return //Реализовать обертку для ошибок
		}
		success := "false"

		resultSQL, err := currentSqlServer.Query("SELECT 'true' WHERE EXISTS (SELECT [GUID] FROM ServerPlanningPoker.Rooms WHERE [GUID] = '" + roomId + "')")
		defer resultSQL.Close()
		if err != nil {
			log.Println(err)
		}
		for resultSQL.Next() {
			err := resultSQL.Scan(&success)
			if err != nil {
				fmt.Println(err)
				continue //Добавить реализацию логирования ошибок в базу
			}
		}
		type ViewData struct {
			Success     string
			RoomId      string
			CurrentHost string
		}
		data := ViewData{
			Success:     success,
			RoomId:      roomId,
			CurrentHost: currentServersSettings.ServerHost.Host,
		}
		fmt.Println(roomId)
		fmt.Println(success)
		if success == "true" {
			//http.Redirect(w, r, "/room/?roomId="+roomId, 301)
			tmpl, _ := template.ParseFiles("templates/rooms.html")
			tmpl.Execute(w, data)
		}

	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login := r.FormValue("login")
		password := r.FormValue("password")
		if len(login) > 0 && len(password) > 0 {
			resultSP, err := currentSqlServer.Query(`EXEC [CheckUser] ?, ?`, login, password)
			if err != nil {
				log.Println(err)
			}
			var resultCkeck bool
			for resultSP.Next() {
				err := resultSP.Scan(&resultCkeck)
				if err != nil {
					fmt.Println(err)
					continue //Добавить реализацию логирования ошибок в базу
				}
			}
			if resultCkeck {
				session := session.NewSession()
				sessionId := session.InitSession(login)
				cookie := &http.Cookie{
					Name:    cookieName,
					Value:   sessionId,
					Expires: time.Now().Add(5 * time.Minute),
				}
				http.SetCookie(w, cookie)
				w.Write([]byte("Sucsess"))
				return
			}
		}
		w.Write([]byte("Wrong"))

	})

	http.HandleFunc("/loginform", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	})

	http.HandleFunc("/room/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/NewRoom", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/newPlanningPokerRoom.html")
	})

	http.HandleFunc("/create-room", func(w http.ResponseWriter, r *http.Request) {
		nameRoom := r.FormValue("nameRoom")
		xmlTasks := r.FormValue("xmlTasks")
		if len(nameRoom) == 0 && len(xmlTasks) == 0 {
			http.ServeFile(w, r, "templates/error_bad_request.html") //реализовать вызов методов страниц ошибок
			return
		}
		var roomId string
		fmt.Println(nameRoom)
		fmt.Println(xmlTasks)
		//ctx := context.Background()
		resultSP, err := currentSqlServer.Query(`EXEC [NewPlanningPokerRoom] @NameRoom=?, @Tasks=?;`, nameRoom, xmlTasks)

		if err != nil {
			log.Println(err)
		}
		var hello string
		for resultSP.Next() {
			err := resultSP.Scan(&hello)
			if err != nil {
				fmt.Println(err)
				continue //Добавить реализацию логирования ошибок в базу
			}
		}
		fmt.Println(currentServersSettings.ServerHost.Host)
		w.Write([]byte(currentServersSettings.ServerHost.Host + currentServersSettings.ServerHost.Room + strings.ToLower(hello)))
		fmt.Println(resultSP)
		fmt.Println(hello)
		fmt.Println(roomId)

	})

	http.ListenAndServe(":80", nil)

}
