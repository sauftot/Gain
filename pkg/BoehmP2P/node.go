package BoehmP2P

import (
	"Gain/internal"
	log "Gain/pkg/Logger"
	"context"
	"fmt"
	"net"
	"sync"
)

/*
PACKAGE-GLOBAL CONFIGURATION VARIABLES
*/

var (
	idBits = 160
	k      = 20
	// alpha  = 3
	port = 40001
)

// package-global logger
var nodeLogger *log.Logger
var logLevel = log.DEBUG

// package-global waitGroup
var wg sync.WaitGroup

/*
PUBLIC STRUCTS AND FUNCTIONS
*/

type Node struct {
	meta     *NodeMeta
	buckets  *BucketList
	ctx      internal.ContextWithCancel
	register chan internal.ContextWithCancel
}

func NewNode(rootCtx context.Context) *Node {
	var err error
	newMeta := GenerateOwnMeta()
	nodeLogger, err = log.NewLogger("Node_" + newMeta.ID.String())
	nodeLogger.SetLogLevel(logLevel)
	if err != nil {
		fmt.Println("Error creating logger: ", err)
		panic(err)
	}

	var cancel context.CancelFunc
	if rootCtx == nil {
		rootCtx, cancel = context.WithCancel(context.Background())
	} else {
		rootCtx, cancel = context.WithCancel(rootCtx)
	}

	register := make(chan internal.ContextWithCancel, 10)

	return &Node{
		meta:    newMeta,
		buckets: NewBucketList(newMeta, register),
		ctx:     internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel},
	}
}

/*
NODE INTERFACE
*/

func (n *Node) Start(boot string) {
	nodeLogger.Log("Starting node...")
	// resolve boot
	var ip net.IP
	ip = net.ParseIP(boot)
	if ip == nil {
		i, err := net.ResolveIPAddr("ip4", boot)
		if err != nil {
			nodeLogger.Error("Error resolving boot node: ", err)
			n.ctx.Cancel()
			return
		}
		ip = i.IP
	}

	// dial boot
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		nodeLogger.Error("Error dialing boot node: ", err)
	}

	// start the network handler TODO: move to separate function and implement callbacks
	wg.Add(1)
	go n.run()

	// send initial message to boot node TODO: make this into context and pass to register for network handler
	initialMsgID := GenerateID()
	initialMsg := NewMessage(&initialMsgID, n.meta, FIND_NODE, n.meta.ID.Bytes())
	_, err = conn.Write(initialMsg.ToBytes())
	if err != nil {
		nodeLogger.Error("Error sending initial message to boot node: ", err)
		n.ctx.Cancel()
		return
	}
}

func (n *Node) Stop() {
	nodeLogger.Log("Stopping node...")
}

func (n *Node) Store() {

}

func (n *Node) Find() {

}

/*
PRIVATE FUNCTIONS
*/

func (n *Node) receivedFromNode(node NodeMeta) {

}

func (n *Node) run() {
	defer wg.Done()
	nodeLogger.Log("run started")

	l, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		nodeLogger.Error("Error listening on UDP: ", err)
		panic(err)
	}
	defer func(l *net.UDPConn) {
		err := l.Close()
		if err != nil {
			nodeLogger.Error("Error closing UDP listener: ", err)
		}
	}(l)
}

func (n *Node) broadcast() {

}

func (n *Node) findNode() {

}

func (n *Node) store() {

}
