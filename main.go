package main

import (
	"PlanningPoker/PlanningPokerSettings"
	"PlanningPoker/RoomInteraction"
	"PlanningPoker/Sessions"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/go-passwd/validator"
	"github.com/google/uuid"

	"crypto/md5"
	"encoding/hex"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var currentServerSettings *PlanningPokerSettings.ServerSettings

func init() {
	currentServerSettings = PlanningPokerSettings.InitServerSetting()
}

var conns []RoomInteraction.Connection

type connections []*websocket.Conn

var sessionsTool = Sessions.InitSessionsTool()

var currentSqlServer *sql.DB
var err error

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	cookieName         = "sessionId"
	lengthGUID         = 36
	regexpPatternEmail = `^[-\w.]+@([A-z0-9][-A-z0-9]+\.)+[A-z]{2,4}$`
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/unknownroom", unknownroomHandler).Methods("GET")
	router.HandleFunc("/bad-request", badRequestHandler).Methods("GET")
	router.HandleFunc("/restore-password", validateDataMiddleware(createRestoreAccountLink)).Methods("POST")
	router.HandleFunc("/restore-account", restoreAccountHandler).Methods("GET")
	router.HandleFunc("/change-password-form", changePasswordFormHandler).Methods("GET")
	router.HandleFunc("/update-password", validateDataMiddleware(restoreAccountUpdateHandler)).Methods("POST")
	router.HandleFunc("/rooms", checkAuthMiddleware(roomsHandler)).Queries("roomId", "")
	router.HandleFunc("/loginform", checkAuthMiddleware(loginFormHandler)).Methods("GET")
	router.HandleFunc("/create-room", checkAuthMiddleware(createRoomHandler)).Methods("POST")
	router.HandleFunc("/newroom", checkAuthMiddleware(newRoomHandler)).Methods("GET")
	router.HandleFunc("/login", validateDataMiddleware(loginHandler)).Methods("POST")
	router.HandleFunc("/echo", echoSocket).Methods("GET")
	router.HandleFunc("/room", checkAuthMiddleware(roomHandler)).Methods("GET")
	router.HandleFunc("/registrationform", registrationFormHandler).Methods("GET")
	router.HandleFunc("/registration", validateDataMiddleware(registrationHandler)).Methods("POST")
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	currentSqlServer, err = PlanningPokerSettings.ServerSql{DSN: currentServerSettings.SQLServer.DSN, TypeSql: currentServerSettings.SQLServer.TypeSql}.OpenConnection()

	fmt.Println("Server succsesful configured. ©Roman Solovyev")

	fmt.Println("Server started on:" + currentServerSettings.ServerHost.Host + currentServerSettings.ServerHost.Port)
	server := &http.Server{Addr: currentServerSettings.ServerHost.Host + currentServerSettings.ServerHost.Port, Handler: router}

	go func() {
		server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println("To stop the server enter the command: stop")

	var key string
	for {
		_, err = fmt.Scan(&key)
		if key == "stop" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			return
		}
	}
}

/*End point to check exist room for id*/
func roomsHandler(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query()["roomId"][0]
	if len(roomId) != lengthGUID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resultSQL, err := currentSqlServer.Query("SELECT 'true' WHERE EXISTS (SELECT [GUID] FROM ServerPlanningPoker.Rooms WHERE [GUID] = '" + roomId + "')")
	defer resultSQL.Close()
	if err != nil {
		log.Println(err)
	}

	var successFindRoom string
	for resultSQL.Next() {
		err := resultSQL.Scan(&successFindRoom)
		if err != nil {
			fmt.Println(err)
		}
	}

	type ViewData struct {
		RoomId      string
		CurrentHost string
	}
	data := ViewData{
		RoomId:      roomId,
		CurrentHost: currentServerSettings.ServerHost.ExternalHostName,
	}

	if successFindRoom == "true" {
		tmpl, _ := template.ParseFiles("templates/rooms.html")
		tmpl.Execute(w, data)
	}
	return
}

/*End point to create new room with tasks*/
func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	userLogin := sessionsTool.GetUserLoginSession(r)
	nameRoom := r.FormValue("nameRoom")
	xmlTasks := r.FormValue("xmlTasks")
	if len(nameRoom) == 0 && len(xmlTasks) == 0 {
		http.Redirect(w, r, "/bad-request", 301)
		return
	}

	resultSP, err := currentSqlServer.Query(`EXEC ServerPlanningPoker.[NewPlanningPokerRoom] @NameRoom=?, @Tasks=?, @Creator=?;`, nameRoom, xmlTasks, userLogin)
	if err != nil {
		log.Println(err)
	}

	var newGUIDRoom string
	for resultSP.Next() {
		err := resultSP.Scan(&newGUIDRoom)
		if err != nil {
			fmt.Println(err)
		}
	}
	if newGUIDRoom == "error" {
		w.Write([]byte("error"))
		return
	}
	w.Write([]byte(currentServerSettings.ServerHost.ExternalHostName + currentServerSettings.ServerHost.Room + strings.ToLower(newGUIDRoom)))
	return
}

/*End point for to enter users*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	passwordMd5 := md5.Sum([]byte(r.FormValue("password")))
	passHash := hex.EncodeToString(passwordMd5[:])
	resultSP, err := currentSqlServer.Query(`EXEC ServerPlanningPoker.[CheckUser] ?, ?`, email, passHash)
	if err != nil {
		log.Println(err)
	}

	var resultCkeckUser bool

	for resultSP.Next() {
		err := resultSP.Scan(&resultCkeckUser)
		if err != nil {
			fmt.Println(err)
		}
	}

	if resultCkeckUser {
		sessionsTool.CreateNewSession(email, r, &w)
		if r.Header.Get("Referer") == currentServerSettings.ServerHost.ExternalPathToLoginForm {
			w.Write([]byte(currentServerSettings.ServerHost.ExternalHostName))
		}
		return
	} else {
		w.Write([]byte("Unsuccsess"))
		return
	}
}

/*End point to display invitation to enter the room*/
func roomHandler(w http.ResponseWriter, r *http.Request) {
	var creatorOrUser string
	userName := sessionsTool.GetUserLoginSession(r)
	roomUID := r.URL.Query()["roomId"][0]

	resultSP, err := currentSqlServer.Query(`EXEC ServerPlanningPoker.[CheckCreator] @email=?, @roomUID=?`, userName, roomUID)
	if err != nil {
		fmt.Println(err)
	}
	for resultSP.Next() {
		err := resultSP.Scan(&creatorOrUser)
		if err != nil {
			fmt.Println(err)
		}
	}

	roomPatterns := RoomInteraction.RoomPatternHtml{
		WebSocketExternalAddress: template.HTML(currentServerSettings.ServerHost.WebSocketExternalAddress),
	}

	if creatorOrUser == "Creator" {
		creatorToolsDataFile, err := ioutil.ReadFile("templates/creatorPatterns/creatorTools.html")
		if err != nil {
			fmt.Println(err)
		}
		creatorScriptsDataFile, err := ioutil.ReadFile("templates/creatorPatterns/creatorRoomScripts.html")
		if err != nil {
			fmt.Println(err)
		}
		creatorStylesDataFile, err := ioutil.ReadFile("templates/creatorPatterns/creatorStyles.html")
		if err != nil {
			fmt.Println(err)
		}
		roomPatterns.CreatorTools = template.HTML(creatorToolsDataFile)
		roomPatterns.CreatorScripts = template.HTML(creatorScriptsDataFile)
		roomPatterns.CreatorStyles = template.HTML(creatorStylesDataFile)

		tmpl, _ := template.ParseFiles("templates/room.html")
		tmpl.Execute(w, roomPatterns)
		return
	} else {
		tmpl, _ := template.ParseFiles("templates/room.html")
		tmpl.Execute(w, roomPatterns)
		return
	}
}

/*End point to upgare request for protocol WebSocket*/
func echoSocket(w http.ResponseWriter, r *http.Request) {
	var (
		URL    = r.URL.Query()["roomId"][0]
		Conn   RoomInteraction.Connection
		Change RoomInteraction.Change = RoomInteraction.NewChangesViewModel()
	)
	Conn.Conn, _ = upgrader.Upgrade(w, r, nil)
	Conn.RoomGUID = URL
	Conn.UUID, _ = uuid.NewUUID()

	Conn.UserEmail = sessionsTool.GetUserLoginSession(r)
	if len(Conn.UserEmail) < 0 {
		http.Redirect(w, r, "/login", 301)
		return
	}

	currentSqlServer.Query(`EXEC ServerPlanningPoker.[CreateConnection] @UUID=?, @RoomGUID=?, @Email=?`, Conn.UUID.String(), Conn.RoomGUID, Conn.UserEmail)
	if err != nil {
		fmt.Println(err)
	}

	conns = append(conns, Conn)

	for {
		msgType, msg, err := Conn.Conn.ReadMessage()
		//Добавить логику работы когда мы закрываем соединение, чтобы удалять ненужные
		if err != nil {
			return
		}
		msgConnection := string(msg)
		msgConnArray := strings.Split(msgConnection, "==")

		var (
			msgConnKey, msgConnValue string
		)
		if len(msgConnArray) == 2 {
			msgConnKey = msgConnArray[0]
			msgConnValue = msgConnArray[1]
		} else if len(msgConnArray) == 1 {
			msgConnKey = msgConnArray[0]
		} else {
			return
		}

		commandSql := Change.GetChange(msgConnKey)

		resultSP, err := currentSqlServer.Query(commandSql, msgConnValue, msgConnKey, Conn.RoomGUID, Conn.UserEmail)
		if err != nil {
			log.Println(err)
		}

		var resultSQL string

		for resultSP.Next() {
			err := resultSP.Scan(&resultSQL)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		for _, value := range conns {
			if value.RoomGUID == URL {
				msg = []byte(resultSQL)
				value.Conn.WriteMessage(msgType, msg)
			}

		}
	}
}

/*End point for registration new users across post request*/
func registrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	userName := r.FormValue("userName")
	passwordMd5 := md5.Sum([]byte(r.FormValue("password")))
	passHash := hex.EncodeToString(passwordMd5[:])

	structCh := make(chan struct{})

	resultSP, err := currentSqlServer.Query(`EXEC ServerPlanningPoker.[Add_User] @LoginName=?, @Email=?, @Password=?`, userName, email, passHash)
	if err != nil {
		log.Println(err)
	}

	var resultCheck string

	for resultSP.Next() {
		err := resultSP.Scan(&resultCheck)
		if err != nil {
			fmt.Println(err)
		}
	}
	if resultCheck == "Succsess" {
		w.Write([]byte(resultCheck))

		go func() {
			defer close(structCh)
			from := currentServerSettings.SmtpServer.LoginHost
			pass := currentServerSettings.SmtpServer.PassHost
			to := email

			content := fmt.Sprintf("Dear User!, %s. You have successfully registered with Planning-poker with login: %s", email, userName)
			msg := "From: " + from + "\n" +
				"To: " + to + "\n" +
				"Subject: Successful registration\n\n" +
				content

			if err := smtp.SendMail(currentServerSettings.SmtpServer.Host+currentServerSettings.SmtpServer.PortHost,
				smtp.PlainAuth("", from, pass, currentServerSettings.SmtpServer.Host),
				from, []string{to}, []byte(msg)); err != nil {
				fmt.Println("Error SendMail: ", err)
			}
		}()

		<-structCh

		return
	} else {
		w.Write([]byte(resultCheck))
		return
	}
}

/*Middleware for auth*/
func checkAuthMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resultCheckCookie := sessionsTool.CheckAndUpdateSession(r, &w)
		if resultCheckCookie {
			nextHandler(w, r)
		} else {
			tmpl, _ := template.ParseFiles("templates/loginForm.html")
			tmpl.Execute(w, nil)
			return
		}
	}
}

/*End point to display login form*/
func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/loginForm.html")
	tmpl.Execute(w, nil)
	return
}

/*End point to display the page to create room*/
func newRoomHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/newPlanningPokerRoom.html")
	tmpl.Execute(w, nil)
	return
}

/*End point to display the page unknown room*/
func unknownroomHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/unknownRoom.html")
	tmpl.Execute(w, nil)
	return
}

/*End point to display the page bad request*/
func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/error_bad_request.html")
	tmpl.Execute(w, nil)
	return
}

/*End point to display the page registration*/
func registrationFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/registrationForm.html")
	tmpl.Execute(w, nil)
	return
}

/*End point to display the index page*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
	return
}

/*End point to display the restore account  page*/
func restoreAccountHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/restoreAccountForm.html")
	tmpl.Execute(w, nil)
	return
}

/*End point for post request to restore password*/
func createRestoreAccountLink(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	type ViewData struct {
		Message string
	}

	tmpl, _ := template.ParseFiles("templates/restore_account.html")

	emailCh := make(chan string)
	tmplMessageCh := make(chan ViewData)

	go func() {
		resultSP, err := currentSqlServer.Query("EXEC ServerPlanningPoker.[Check_User_Email] @Email=?", email)
		if err != nil {
			fmt.Println(err)
		}

		var resultEmail string

		for resultSP.Next() {
			err := resultSP.Scan(&resultEmail)
			if err != nil {
				fmt.Println(err)
			}
		}

		emailCh <- resultEmail
	}()

	go func() {
		resultSP, err := currentSqlServer.Query("EXEC ServerPlanningPoker.[CreateAccountRecoveryLink] @Email=?", email)
		if err != nil {
			fmt.Println(err)
		}

		var recoveryLnk string

		for resultSP.Next() {
			err := resultSP.Scan(&recoveryLnk)
			if err != nil {
				fmt.Println(err)
			}
		}

		emailCurrCtx := <-emailCh

		if emailCurrCtx == email {

			from := currentServerSettings.SmtpServer.LoginHost
			pass := currentServerSettings.SmtpServer.PassHost
			to := emailCurrCtx

			content := fmt.Sprintf("Dear User!, %s. An attempt was made to restore your account. To change your password, follow the link: %s", email,
				currentServerSettings.ServerHost.ExternalHostName+currentServerSettings.ServerHost.RestoreAccount+strings.ToLower(recoveryLnk))
			msg := "From: " + from + "\n" +
				"To: " + to + "\n" +
				"Subject: An attempt was made to restore your account\n\n" +
				content

			if err := smtp.SendMail(currentServerSettings.SmtpServer.Host+currentServerSettings.SmtpServer.PortHost,
				smtp.PlainAuth("", from, pass, currentServerSettings.SmtpServer.Host),
				from, []string{to}, []byte(msg)); err != nil {
				fmt.Println("Error SendMail: ", err)
			}

			data := ViewData{Message: "An account recovery request has been generated, pls check in your email."}

			tmplMessageCh <- data

		} else {
			data := ViewData{Message: "Such account was not found."}
			tmplMessageCh <- data
		}

	}()

	tmpl.Execute(w, <-tmplMessageCh)
}

/*End point to display the restore account  page*/
func changePasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	linkRestore := r.URL.Query().Get("linkRestore")

	type ViewData struct {
		RestoreLink string
	}

	data := ViewData{RestoreLink: linkRestore}
	tmpl, _ := template.ParseFiles("templates/changePasswordForm.html")
	tmpl.Execute(w, data)
	return
}

/*End point for change password*/
func restoreAccountUpdateHandler(w http.ResponseWriter, r *http.Request) {
	linkRestore := r.FormValue("linkRestore")
	passwordMd5 := md5.Sum([]byte(r.FormValue("password")))
	passHash := hex.EncodeToString(passwordMd5[:])
	resultSP, err := currentSqlServer.Query("EXEC ServerPlanningPoker.[RestoreAccount] @Link=?, @Password=?", linkRestore, passHash)
	if err != nil {
		fmt.Println(err)
	}

	var resultUpdPass bool

	for resultSP.Next() {
		err := resultSP.Scan(&resultUpdPass)
		if err != nil {
			fmt.Println(err)
		}
	}

	tmpl, _ := template.ParseFiles("templates/restore_account_success.html") //дописать проставку неактивной ссылки и проверку не протухшей ссылки для обновления пароля

	type ViewData struct {
		Message string
	}

	if resultUpdPass == true {
		data := ViewData{Message: "restored and change password"}
		tmpl.Execute(w, data)
	} else {
		data := ViewData{Message: "not restore because your link was use"}
		tmpl.Execute(w, data)
	}
	return
}

/*Middleware for check validation fields*/
func validateDataMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkEmailData = func(w http.ResponseWriter, r *http.Request) {
			email := r.FormValue("email")

			resultMatchEmail, err := regexp.MatchString(regexpPatternEmail, email)
			if err != nil {
				fmt.Println(err)
			}

			if resultMatchEmail == false {
				w.WriteHeader(400) //Проверить проверку пароля (не работает)
				return
			} else {
				nextHandler(w, r)
			}
		}
		var checkPasswordData = func(w http.ResponseWriter, r *http.Request) {
			password := r.FormValue("password")

			validator := validator.New(validator.MinLength(6, nil), validator.MaxLength(24, nil),
				validator.ContainsOnly(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!#$%&'()*+,-./:;<=>?@[\]^_{|}~`, nil),
				validator.ContainsAtLeast(`1234567890!#$%&'()*+,-./:;<=>?@[\]^_{|}~`, 1, nil))
			err := validator.Validate(password)
			if err != nil {
				w.WriteHeader(400) //Проверить проверку пароля (не работает)
				return
			} else {
				nextHandler(w, r)
			}
		}
		var checkLoginFormData = func(w http.ResponseWriter, r *http.Request) {
			email := r.FormValue("email")
			password := r.FormValue("password")

			validator := validator.New(validator.MinLength(6, nil), validator.MaxLength(24, nil),
				validator.ContainsOnly(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!#$%&'()*+,-./:;<=>?@[\]^_{|}~`, nil),
				validator.ContainsAtLeast(`1234567890!#$%&'()*+,-./:;<=>?@[\]^_{|}~`, 1, nil))
			err := validator.Validate(password)
			if err != nil {
				w.Write([]byte("Unsuccsess")) //Проверить проверку пароля (не работает)
				return
			}
			resultMatchEmail, err := regexp.MatchString(regexpPatternEmail, email)
			if err != nil {
				fmt.Println(err)
			}

			if resultMatchEmail == false || err != nil {
				w.Write([]byte("Unsuccsess")) //Проверить проверку пароля (не работает)
				return
			} else {
				nextHandler(w, r)
			}
		}
		var checkRegistrationFormData = func(w http.ResponseWriter, r *http.Request) {
			username := r.FormValue("userName")

			validator := validator.New(validator.MinLength(3, nil), validator.MaxLength(24, nil),
				validator.ContainsOnly(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890`, nil))
			err := validator.Validate(username)

			if err != nil {
				w.Write([]byte("pass or email or username no valid"))
				return
			} else {
				checkLoginFormData(w, r)
			}
		}
		rqstURL := r.Header.Get("Referer")

		switch {
		case currentServerSettings.ServerHost.ExternalPathToNewRoom == rqstURL:
			checkLoginFormData(w, r)
		case currentServerSettings.ServerHost.ExternalPathToLoginForm == rqstURL:
			checkLoginFormData(w, r)
		case strings.Contains(rqstURL, "roomId="):
			checkLoginFormData(w, r)
		case currentServerSettings.ServerHost.ExternalPathToRegistrationForm == rqstURL:
			checkRegistrationFormData(w, r)
		case currentServerSettings.ServerHost.ExternalPathToRestoreAcc == rqstURL:
			checkEmailData(w, r)
		case strings.Contains(rqstURL, "linkRestore="):
			checkPasswordData(w, r)
		default:
			return
		}
	}
}
