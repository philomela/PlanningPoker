package ServerPlanningPoker

import "net"

type Server struct {
	ProtocolServer string
	Port           string
}

func (s Server) GetServer() (net.Listener, error) {
	return net.Listen(s.ProtocolServer, s.Port)
}
