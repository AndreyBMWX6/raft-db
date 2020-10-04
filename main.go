package main

import (
	"fmt"
	"container/list"
	//"net"
)

type Cluster struct {  // Кластер
	nodeCount uint
	nodes []Node
}

type ClusterInterface interface {
	Cluster(nodesCount_ uint) *Cluster
}

func (c *Cluster) Cluster(nodesCount_ uint) *Cluster {
	c.nodeCount = nodesCount_
	c.nodes = make([]Node, c.nodeCount)
	var nodesPtrs [] *Node
	nodesPtrs = make([]*Node, c.nodeCount)
	for i := 0; uint(i) < c.nodeCount; i++ {
		nodesPtrs[i] = &c.nodes[i]
	}
	for _, node := range nodesPtrs {
		node.Node(0, &nodesPtrs)
	}
	return c
}

type Record struct {
	query string
	termNumber uint
	logIndex uint
}

type Node struct {  // Ноды
	//ip net.IPAddr
	//port uint
	termNumber uint
	votedAmount uint
	nodesPtrs [] *Node
	nodesCount uint
	nodeVote bool
	recordsLog list.List
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
		n.termNumber = termNumber
	}
	return termNumber >= n.termNumber
}

func (n *Node) SendRequest(termNumber uint) {
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
	return
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
	return n.nodeVote // заменить на сравнение метаданных(срок посл. эл-таб длина) лога
}

func main() {
	var c Cluster
	c.Cluster(5)
	c.nodes[1].termNumber = 1
	c.nodes[2].termNumber = 2
	c.nodes[3].termNumber = 2
	c.nodes[4].termNumber = 3

	c.nodes[1].SendRequest(c.nodes[1].termNumber)
	c.nodes[4].SendRequest(c.nodes[4].termNumber)
	fmt.Println("Voted for Node 1: ", c.nodes[1].votedAmount)
	fmt.Println("Voted for Node 4: ", c.nodes[4].votedAmount)
}
