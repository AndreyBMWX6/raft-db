package main

import (
	"fmt"
	//"net"
)

type Node struct {  // Ноды
	//ip net.IPAddr
	//port uint
	status string
	termNumber uint
	votedAmount uint
	nodesPtrs [] *Node // общее кол-во нод = размер nodesPtrs + 1
	nodeVote bool
	entiresLog []Entire
}

type Cluster []Node // Кластер, кол-во нод = размер слайса

type ClusterInterface interface {
	Cluster(nodesCount_ uint) *Cluster
}

func (c *Cluster) Cluster(nodesCount_ uint) *Cluster {
	*c = make([]Node, nodesCount_)
	var nodesPtrs [] *Node
	nodesPtrs = make([]*Node, len(*c))
	for i := 0; i < len(*c); i++ {
		nodesPtrs[i] = &(*c)[i]
	}
	for _, node := range nodesPtrs {
		if node == nodesPtrs[0] {
			node.status = "l"
		} else {
			node.status = "f"
		}
		node.Node(0, &nodesPtrs)
	}
	return c
}

type Entire struct {
	query []byte
	termNumber uint
	logIndex uint
}

type Response struct {
	termNumber uint
	response bool
}

type NodeInterface interface {
	Node(termNumber_ uint, nodesPtr_ *[]Node) *Node
	SendRequest(termNumber uint) []Response
	RequestVote(node Node, termNumber uint) Response
	ReplyVote(termNumber uint) bool
	CheckTermNumber(termNumber uint) bool
	LeaderCheck() bool
	//Election() // для корректной работы необходимо реализовать heartbeat и инициацию выборов
	// 				в случае "молчания" лидера
	AppendEntires()
	ReplyAppend(termNumber uint, entires [] Entire) Response
	ProcessQuery(query []byte)
}

func (n *Node) Node(termNumber_ uint, nodesPtrs_ *[] *Node) *Node {
	n.termNumber = termNumber_
	//n.nodesPtrs = make([]*Node, len(*nodesPtrs_) - 1)
	for _, nodePtr_ := range *nodesPtrs_ {
		if nodePtr_ == n {
			continue
		}
		n.nodesPtrs = append(n.nodesPtrs, nodePtr_)
	}
	n.nodeVote = true
	return n
}

func (n*Node) CheckTermNumber(termNumber uint) bool {
	if termNumber > n.termNumber {
		if n.status == "l" || n.status == "c" {
			n.status = "f"
		}
		n.termNumber = termNumber
	}
	return termNumber >= n.termNumber
}

/*
func (n *Node) Election() {
	n.SendRequest(n.termNumber)
	if n.LeaderCheck() {
		n.Election()
	}
}
 */

func (n* Node) AppendEntires() {
	for _, follower := range n.nodesPtrs {
		if len(n.entiresLog) < 2 {
			follower.ReplyAppend(n.termNumber, n.entiresLog[:])
		} else {
			follower.ReplyAppend(n.termNumber, n.entiresLog[(len(n.entiresLog)-2):]) // сделать бескончный ретрай при возвращении false для конкретной ноды
			 }
	}
}


/* не нужны, т.к можно достать из entires
почему-то передачу лога по указателю сделать не получилось
метаданные: ,lastTermNumber uint, logLength uint*/
func (n* Node) ReplyAppend(termNumber uint, entires [] Entire) Response {
	if n.entiresLog == nil {
		n.entiresLog = append(n.entiresLog, entires[len(entires) - 1])
		return Response{n.entiresLog[len(n.entiresLog) - 1].termNumber, true}
	}
	if n.entiresLog[len(n.entiresLog) - 1].termNumber == entires[len(entires) - 2].termNumber { // непонятно сравниваются записи или их сроки
		n.entiresLog = append(n.entiresLog, entires[len(entires) - 1])
		return Response{n.entiresLog[len(n.entiresLog) - 1].termNumber, true}
	}
	return Response{n.entiresLog[len(n.entiresLog) - 1].termNumber, false}
}

func (n *Node) SendRequest(termNumber uint) {
	n.status = "c"
	n.votedAmount++ // node votes for himself
	n.nodeVote = false
	responses := make([]Response, 0)
	i := 0
	for _, node := range n.nodesPtrs {
		if  node.CheckTermNumber(termNumber) {
			responses = append(responses, n.RequestVote(node, termNumber))
			if responses[i].response {
				n.votedAmount++
			}
			i++
		}
	}
	if int(n.votedAmount) > len(n.nodesPtrs) + 1 {
		n.status = "l"
	}
	return
}

func (n *Node) LeaderCheck() bool {
	if n.status == "l" {
		return true
	}
	for _, node := range n.nodesPtrs {
		if node.status == "l" {
			return true
		}
	}
	return false
}

func (n *Node) RequestVote(node *Node, termNumber uint) Response {
	return Response{node.termNumber, node.ReplyVote()}
}

func (n *Node) ReplyVote(/*lastRecordTermNumber uint, recordsLogLength uint*/) bool {

	if n.nodeVote {
		reply := n.nodeVote
		n.nodeVote = false
		return reply
	}
	return n.nodeVote // заменить на сравнение метаданных(срок посл. эл-та, длина) лога
}

func (n *Node) ProcessQuery(query []byte) {
	if n.status != "l" {
		// нужно в каждую ноду добавить указатель на лидера так будет удобнее
		// (return) n.leaderPtr.ProcessQuery(query)
	}
	n.entiresLog = append(n.entiresLog, Entire{query, n.termNumber, uint(len(n.entiresLog) + 1)})
}

func main() {
	var c Cluster
	c.Cluster(5)
	c[1].termNumber = 1
	c[2].termNumber = 2
	c[3].termNumber = 2
	c[4].termNumber = 3

	c[1].SendRequest(c[1].termNumber)
	c[4].SendRequest(c[4].termNumber)
	c[0].ProcessQuery([]byte("x->3"))
	c[0].AppendEntires()
	fmt.Println("Voted for Node 1: ", c[1].votedAmount)
	fmt.Println("Voted for Node 4: ", c[4].votedAmount)
}
