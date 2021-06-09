package aoi

import (
	"fmt"
	"testing"
)

type Stub struct {
	aoiNode		*NodeAOI
}

func NewNode(x,y,r float32) *NodeAOI {
	s := new(Stub)
	s.aoiNode = NewNodeAOI(x,y,r)
	s.aoiNode.SetListener(s)
	return s.aoiNode
}

func (s *Stub) OnEnterAOI(nodes []*NodeAOI)  {
	for _,node := range nodes {
		fmt.Printf("(%f,%f) OnEnterAOI: (%f,%f) \n", s.aoiNode.x, s.aoiNode.y, node.x, node.y)
	}
}

func (s *Stub) OnLeaveAOI(nodes []*NodeAOI)  {
	for _,node := range nodes {
		fmt.Printf("(%f,%f) OnLeaveAOI: (%f,%f) \n", s.aoiNode.x, s.aoiNode.y, node.x, node.y)
	}
}

func (s *Stub) OnUpdateAOI(node *NodeAOI) {

}

func TestAOI(t *testing.T)  {
	a := NewAOI()
	n1 := NewNode(1,5,2)
	n2 := NewNode(6,6,2)
	n3 := NewNode(3,1,2)
	n4 := NewNode(2,2,2)
	n5 := NewNode(5,3,2)
	a.EnterAOI(n1)
	a.EnterAOI(n2)
	a.EnterAOI(n3)
	a.EnterAOI(n4)
	a.EnterAOI(n5)

	//a.LeaveAOI(n5)

	n6 := NewNode(3,3,2)
	a.EnterAOI(n6)

	a.MoveAOI(n6,4,4)
	a.Print()
}