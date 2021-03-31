package RoomInteraction

import (
	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomPatternHtml struct {
	CreatorTools             template.HTML
	CreatorScripts           template.HTML
	CreatorStyles            template.HTML
	WebSocketExternalAddress template.HTML
}

type Connection struct {
	Conn      *websocket.Conn
	RoomGUID  string
	UUID      uuid.UUID
	UserEmail string
}

type Change interface {
	GetChange(nameEvent string) (commandOnChange string)
}

type ChangesViewModel struct {
	Changes map[string]string
}

func NewChangesViewModel() *ChangesViewModel {
	ChangesViewModelOut := ChangesViewModel{
		Changes: make(map[string]string),
	}
	ChangesViewModelOut.Changes["ChangeVote"] = `EXEC ServerPlanningPoker.[Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`
	ChangesViewModelOut.Changes["ChangeGetVM"] = `EXEC ServerPlanningPoker.[Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`
	ChangesViewModelOut.Changes["StartVoting"] = `EXEC ServerPlanningPoker.[Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`
	ChangesViewModelOut.Changes["StopVoting"] = `EXEC ServerPlanningPoker.[Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`
	ChangesViewModelOut.Changes["FinishPlanning"] = `EXEC ServerPlanningPoker.[Push_And_Get_Changes] @xmlChanges=?, @nameChanges=?, @roomGUID=?, @email=?`

	return &ChangesViewModelOut
}

func (c *ChangesViewModel) GetChange(nameEvent string) (commandOnChange string) {
	commandOnChange = c.Changes[nameEvent]
	return
}
