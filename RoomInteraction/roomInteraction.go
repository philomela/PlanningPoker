package RoomInteraction

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Conn      *websocket.Conn
	RoomGUID  string
	UUID      uuid.UUID
	UserEmail string
}

type Change interface {
	GetChange(nameEvent string, conn *Connection) (commandOnChange string)
}

type ChangesViewModel struct {
	Changes map[string]string
}

func NewChangesViewModel() *ChangesViewModel {
	ChangesViewModelOut := ChangesViewModel{
		Changes: make(map[string]string),
	}
	ChangesViewModelOut.Changes["ChangeVote"] = `EXEC [Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`
	ChangesViewModelOut.Changes["ChangeGetVM"] = `EXEC [Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`

	return &ChangesViewModelOut
}

func (c *ChangesViewModel) GetChange(nameEvent string, conn *Connection) (commandOnChange string) {
	commandOnChange = c.Changes[nameEvent]
	return
}
