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
	meta    *NodeMeta
	buckets *BucketList
	ctx     internal.ContextWithCancel
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
	out := make(chan Message, 100)

	go func() {
		
	}

	var cancel context.CancelFunc
	if rootCtx == nil {
		rootCtx, cancel = context.WithCancel(context.Background())
	} else {
		rootCtx, cancel = context.WithCancel(rootCtx)
	}
	return &Node{
		meta:    newMeta,
		buckets: NewBucketList(newMeta),
		ctx:     internal.ContextWithCancel{Ctx: rootCtx, Cancel: cancel},
	}
}

/*
NODE INTERFACE
*/

func (n *Node) Start() {
	nodeLogger.Log("Starting node...")
}

func (n *Node) Stop() {
	nodeLogger.Log("Stopping node...")
}

func (n *Node) JoinNetwork() {
	nodeLogger.Log("Joining network...")
}

func (n *Node) Broadcast() {
	nodeLogger.Log("Broadcasting...")
}

func (n *Node) SendToNode() {
	nodeLogger.Log("Sending to node...")
}

/*
PRIVATE FUNCTIONS
*/

func (n *Node) receivedFromNode(node NodeMeta) {

}

func (n *Node) run() {
	wg.Add(1)
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

func (n *Node) ping() {

}
