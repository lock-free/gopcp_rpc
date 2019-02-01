package gopcp_rpc

import (
	"github.com/idata-shopee/goaio"
	"github.com/idata-shopee/gopcp"
	"net"
	"sync"
)

type OnCloseHandler = func(error)

func GetPCPServer(port int, sandbox *gopcp.Sandbox) (goaio.TcpServer, error) {
	return goaio.GetTcpServer(port, func(conn net.Conn) goaio.ConnectionHandler {
		pcpClient := gopcp.PcpClient{}
		pcpServer := gopcp.NewPcpServer(sandbox)
		var remoteCallMap sync.Map
		pcpConnectionHandler := PCPConnectionHandler{GetPackageProtocol(), pcpClient, pcpServer, nil, remoteCallMap}
		connHandler := goaio.ConnectionHandler{conn, pcpConnectionHandler.OnData, func(err error) {}}

		pcpConnectionHandler.connHandler = &connHandler
		return connHandler
	})
}

func GetPCPClient(host string, port int, onClose OnCloseHandler) (*PCPConnectionHandler, error) {
	pcpClient := gopcp.PcpClient{}
	pcpServer := gopcp.NewPcpServer(gopcp.GetSandbox(map[string]*gopcp.BoxFunc{}))
	var remoteCallMap sync.Map
	pcpConnectionHandler := PCPConnectionHandler{GetPackageProtocol(), pcpClient, pcpServer, nil, remoteCallMap}
	tcpClient, err := goaio.GetTcpClient(host, port, pcpConnectionHandler.OnData, onClose)

	pcpConnectionHandler.connHandler = &tcpClient

	if err != nil {
		return nil, err
	}

	return &pcpConnectionHandler, nil
}
