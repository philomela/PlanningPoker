package main

import (
	RoomInteraction "ServerPlanningPoker/RoomInteraction"
	ServerPlanningPoker "ServerPlanningPoker/ServerPlanningPokerSettings"
	ServerPlanningPokerSettings "ServerPlanningPoker/ServerPlanningPokerSettings"
	Sessions "ServerPlanningPoker/Sessions"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var conns []RoomInteraction.Connection

type connections []*websocket.Conn

//var rooms []Room //Коллекция комнат, вынести в отдельный пакет
//collectionRoom := make(map[string]connections)
/* Открываем, читаем и парсим json */

var client *redis.Client
var sessionsTool = Sessions.InitSessionsTool()

var currentServersSettings ServerPlanningPokerSettings.ServersSettings
var currentSqlServer *sql.DB
var err error

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	cookieName = "sessionId"
)

func main() {
	//eventSthutdown := make(chan string)
	router := mux.NewRouter()
	router.HandleFunc("/rooms", roomsHandler).Methods("GET")
	router.HandleFunc("/loginform", loginFormHandler).Methods("GET")
	router.HandleFunc("/create-room", createRoomHandler).Methods("POST")
	router.HandleFunc("/newroom", newRoomHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/echo", echoSocket).Methods("GET")
	router.HandleFunc("/room", roomHandler).Methods("GET")
	router.HandleFunc("/registrationform", registrationFormHandler).Methods("GET")
	router.HandleFunc("/registration", registrationHandler).Methods("POST")
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	jsonFile, err := os.Open("./serversSettings.json")

	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	} else {
		defer jsonFile.Close()

		byteArrayJsonSettings, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteArrayJsonSettings, &currentServersSettings)
	}

	currentSqlServer, err = ServerPlanningPoker.ServerSql{DSN: currentServersSettings.SQLServer.DSN, TypeSql: currentServersSettings.SQLServer.TypeSql}.OpenConnection()

	fmt.Println("Server succsesful configured. ©Roman Solovyev")

	//APP_IP := os.Getenv("APP_IP")
	//APP_PORT := os.Getenv("APP_PORT")

	APP_IP := "localhost"
	APP_PORT := "8080"
	fmt.Println("Server started on:" + APP_IP + ":" + APP_PORT)
	server := &http.Server{Addr: APP_IP + ":" + APP_PORT, Handler: router}

	go func() {
		server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	var key string
	for {
		fmt.Println("To stop the server enter the command: stop")
		_, err = fmt.Scan(&key)
		if key == "stop" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}
	}

}

/*Метод проверки существования комнаты по id*/
func roomsHandler(w http.ResponseWriter, r *http.Request) {
	resultCheckCookie := sessionsTool.CheckAndUpdateSession(r, &w)
	println(resultCheckCookie)
	if !resultCheckCookie {
		return //Проверка сессии в cookie
	}
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

		tmpl, _ := template.ParseFiles("templates/rooms.html")
		tmpl.Execute(w, data)
	}
}

/*Метод создания новой комнаты с задачами*/
func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	resultCheckCookie := sessionsTool.CheckAndUpdateSession(r, &w)
	if !resultCheckCookie {
		return //Добавить сообщение не авторизованного пользователя
	}
	userLogin := sessionsTool.GetUserLoginSession(r)
	nameRoom := r.FormValue("nameRoom")
	xmlTasks := r.FormValue("xmlTasks")
	if len(nameRoom) == 0 && len(xmlTasks) == 0 {
		http.ServeFile(w, r, "templates/error_bad_request.html") //реализовать вызов методов страниц ошибок
		return
	}
	var roomId string
	fmt.Println(nameRoom)
	fmt.Println(xmlTasks)

	resultSP, err := currentSqlServer.Query(`EXEC [NewPlanningPokerRoom] @NameRoom=?, @Tasks=?, @Creator=?;`, nameRoom, xmlTasks, userLogin)

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
}

/*Метод входа пользователей*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("loginUser")
	password := r.FormValue("password")
	//Проверки логин/пароль
	if len(login) > 0 && len(password) > 0 {
		resultSP, err := currentSqlServer.Query(`EXEC [CheckUser] ?, ?`, login, password)
		if err != nil {
			log.Println(err)
		}
		var resultCkeck bool
		for resultSP.Next() {
			err := resultSP.Scan(&resultCkeck)
			fmt.Println(resultCkeck)
			if err != nil {
				fmt.Println(err)
				continue //Добавить реализацию логирования ошибок в базу
			}
		}
		if resultCkeck {

			sessionsTool.CreateNewSession(login, r, &w)
			requestURL := r.Header.Get("Referer")
			println(requestURL)
			if requestURL == currentServersSettings.ServerHost.Host+"newroom" {

				w.Write([]byte(requestURL))
			}
			return
		}
	}

	w.Write([]byte("Wrong"))
}

/*Метод отображения формы входа*/
func loginFormHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "templates/loginForm.html")
}

/*Метод отображения страницы создания комнаты*/
func newRoomHandler(w http.ResponseWriter, r *http.Request) {
	resultCheckCookie := sessionsTool.CheckAndUpdateSession(r, &w)

	fmt.Println(resultCheckCookie)

	if resultCheckCookie {
		tmpl, _ := template.ParseFiles("templates/newPlanningPokerRoom.html")
		tmpl.Execute(w, nil)
		return

	} else {
		tmpl, _ := template.ParseFiles("templates/loginForm.html")
		tmpl.Execute(w, nil)
		return
	}

}

/*Метод отображения приглашения входа в комнату*/
func roomHandler(w http.ResponseWriter, r *http.Request) {
	var creatorOrUser string
	resultCheckCookie := sessionsTool.CheckAndUpdateSession(r, &w)
	userName := sessionsTool.GetUserLoginSession(r)
	roomUID := r.URL.Query()["roomId"][0]

	resultSP, err := currentSqlServer.Query(`EXEC [CheckCreator] @email=?, @roomUID=?`, userName, roomUID)
	if err != nil {
		fmt.Println(err)
	}
	for resultSP.Next() {
		err := resultSP.Scan(&creatorOrUser)
		fmt.Println(creatorOrUser)
		fmt.Println(roomUID)
		if err != nil {
			fmt.Println(err)
			continue //Добавить реализацию логирования ошибок в базу
		}
	}

	if resultCheckCookie {
		if creatorOrUser == "Creator" {
			dataFile, err := ioutil.ReadFile("templates/creatorTools.html")
			if err != nil {
				fmt.Println(err)
			}
			dataHtml := template.HTML(dataFile)

			tmpl, _ := template.ParseFiles("templates/room.html")
			tmpl.Execute(w, dataHtml)
			return
		} else {
			tmpl, _ := template.ParseFiles("templates/room.html")
			tmpl.Execute(w, nil)
			return
		}

	} else {
		tmpl, _ := template.ParseFiles("templates/loginForm.html")
		tmpl.Execute(w, nil)
		return
	}
	return
}

type PreHandlerWebSocket struct {
	Map []string
}

func checkAuth(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		nextHandler(w, r)
	}
}

/*Метод upgare запроса до протокола WebSocket*/
func echoSocket(w http.ResponseWriter, r *http.Request) {
	var (
		URL    = r.URL.Query()["roomId"][0]
		Conn   RoomInteraction.Connection
		Change RoomInteraction.Change = RoomInteraction.NewChangesViewModel()
	)
	Conn.Conn, _ = upgrader.Upgrade(w, r, nil)
	Conn.RoomGUID = URL
	Conn.UUID, err = uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
		return //Добавить реализацию логирования ошибок в базу
	}
	Conn.UserEmail = sessionsTool.GetUserLoginSession(r)
	if len(Conn.UserEmail) < 0 {
		http.Redirect(w, r, "/login", 301) //пересмотреть, на сокете редирект может не работать
		return
	}

	currentSqlServer.Query(`EXEC [CreateConnection] @UUID=?, @RoomGUID=?, @Email=?`, Conn.UUID.String(), Conn.RoomGUID, Conn.UserEmail)
	conns = append(conns, Conn)
	println(conns)
	println(URL)
	for {
		msgType, msg, err := Conn.Conn.ReadMessage()
		//Добавить логику работы когда мы закрываем соединение, чтобы удалять ненужные
		if err != nil {
			return
		}
		msgConnection := string(msg)
		msgConnArray := strings.Split(msgConnection, "==")
		//msgConnValueArray := strings.Split(msgConnection, "==")
		fmt.Println(msgConnArray)
		var (
			msgConnKey, msgConnValue string
		)
		if len(msgConnArray) == 2 {
			msgConnKey = msgConnArray[0]
			msgConnValue = msgConnArray[1]
		} else if len(msgConnArray) == 1 {
			msgConnKey = msgConnArray[0]
		} else {
			return //Рассмотреть детальнее обработку и филтрацию команд
		}

		fmt.Println(string(msgConnKey))
		fmt.Println(string(msgConnValue))

		commandSql := Change.GetChange(msgConnKey, &Conn)
		fmt.Println(commandSql)
		resultSP, err := currentSqlServer.Query(commandSql, msgConnValue, msgConnKey, Conn.RoomGUID, Conn.UserEmail)
		if err != nil {
			log.Println(err)
		}
		var resultCkeck string
		for resultSP.Next() {
			err := resultSP.Scan(&resultCkeck)
			fmt.Println(resultCkeck)
			if err != nil {
				fmt.Println(err)
				continue //Добавить реализацию логирования ошибок в базу
			}
		}

		for _, value := range conns {
			if value.RoomGUID == URL {
				msg = []byte(resultCkeck)
				value.Conn.WriteMessage(msgType, msg)
			}

		}
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	userName := r.FormValue("userName")
	password := r.FormValue("password")

	if len(email) > 0 && len(password) > 0 && len(userName) > 0 {
		resultSP, err := currentSqlServer.Query(`EXECUTE Add_User @LoginName=?, @Email=?, @Password=?`, userName, email, password)
		if err != nil {
			log.Println(err)
		}
		var resultCheck string
		for resultSP.Next() {
			err := resultSP.Scan(&resultCheck)
			if err != nil {
				fmt.Println(err)
				continue //Добавить реализацию логирования ошибок в базу
			}
		}
		if resultCheck == "Succsess" {
			w.Write([]byte(resultCheck))
			go func(ss ServerPlanningPoker.ServersSettings) {
				from := ss.SmtpServer.LoginHost
				to := "romaphilomela@yandex.ru"
				host := ss.SmtpServer.Host
				auth := smtp.PlainAuth("", from, ss.SmtpServer.PassHost, host)
				fmt.Println(ss.SmtpServer.PassHost)
				message := "To: romaphilomela@yandex.ru\r\n" +
					"Subject: discount Gophers!\r\n" +
					"\r\n" +
					"This is the email body.\r\n"

				if err := smtp.SendMail(host+ss.SmtpServer.PortHost, auth, from, []string{to}, []byte(message)); err != nil {
					fmt.Println("Error SendMail: ", err)
				}
				fmt.Println("Email Sent!") //Перепроверить отправку сообщений
			}(currentServersSettings)
			return
		} else {
			w.Write([]byte(resultCheck))
			return
		}
	} else {
		w.Write([]byte("pass or email or login was empty"))
		return
	}
}

func registrationFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/registrationForm.html")
	tmpl.Execute(w, nil)
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
	return
}
