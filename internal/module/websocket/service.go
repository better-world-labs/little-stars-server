package websocket

import (
	"aed-api-server/internal/pkg/network"

	"gitlab.openviewtech.com/openview-pub/gopkg/log"
)

var server *network.WSServer

type WsAgent struct {
	conn *network.WSConn
}

func (a WsAgent) Run() {
	for {
		log.DefaultLogger().Debug("NewAgent:run")
		str, err := a.conn.ReadMsg()
		if err != nil {
			log.DefaultLogger().Errorf("NewAgent:error %v", err)
			return
		}

		log.DefaultLogger().Debugf("******* NewAgent:run str:%v", string(str))
	}
}

func (a WsAgent) OnClose() {
	a.conn.Close()
}

func NewAgent(s *network.WSConn) network.Agent {
	log.DefaultLogger().Debugf("NewAgent: %s", s.RemoteAddr().String())
	return WsAgent{conn: s}
}

func StartListen() {
	server = &network.WSServer{}
	server.NewAgent = NewAgent
	server.MaxConnNum = 1000
	server.PendingWriteNum = 100
	server.MaxMsgLen = 1024 * 10
	server.Addr = "127.0.0.1:3653"
	server.Start()
}
