package aoi

import (
	"fmt"
	"math"
	"sync"
)

// AOI node
type NodeAOI struct {
	xPrev 		*NodeAOI
	xNext 		*NodeAOI
	yPrev 		*NodeAOI
	yNext 		*NodeAOI

	x 			float32
	y 			float32
	r 			float32		//range
	listener 	AOIEvent
	object 		interface{}
}

// Create a node
func NewNodeAOI(x,y float32, r float32) *NodeAOI {
	return &NodeAOI{
		x:      x,
		y:      y,
		r:      r,
	}
}

// Range can be changed anytime
func (n *NodeAOI) SetListener(listener AOIEvent) {
	n.listener = listener
}

func (n *NodeAOI) BindObject(obj interface{}) {
	n.object = obj
}

func (n *NodeAOI) GetObject() interface{} {
	return n.object
}

// Range can be changed anytime
func (n *NodeAOI) SetRange(r float32) {
	n.r = r
}

// AOI List
type ListAOI struct {
	xList 	*NodeAOI
	yList 	*NodeAOI
	lock 	sync.RWMutex
}

// Create list AOI
func NewAOI() AOIAction {
	return new(ListAOI)
}

// Add a node to list
func (l *ListAOI) add(node *NodeAOI) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.xList == nil || l.yList == nil{
		l.xList, l.yList = node, node
		return
	}
	var tail *NodeAOI
	found := false
	for cur := l.xList; cur != nil; cur = cur.xNext {
		if cur.x > node.x {
			node.xNext = cur
			if cur.xPrev != nil{
				node.xPrev = cur.xPrev
				cur.xPrev.xNext = node
			} else {
				l.xList = node
			}
			cur.xPrev = node
			found = true
			break
		}
		tail = cur
	}
	if tail != nil && !found{
		tail.xNext = node
		node.xPrev = tail
	}

	tail = nil
	found = false
	for cur := l.yList; cur != nil; cur = cur.yNext {
		if cur.y > node.y {
			node.yNext = cur
			if cur.yPrev != nil{
				node.yPrev = cur.yPrev
				cur.yPrev.yNext = node
			} else {
				l.yList = node
			}
			cur.yPrev = node
			found = true
			break
		}
		tail = cur
	}
	if tail != nil && !found{
		tail.yNext = node
		node.yPrev = tail
	}
}

// Remove the node from list
func (l *ListAOI) remove(node *NodeAOI) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if node == l.xList {
		if node.xNext != nil {
			l.xList = node.xNext
			if l.xList.xPrev != nil {
				l.xList.xPrev = nil
			}
		} else {
			l.xList = nil
		}
	} else if node.xPrev != nil && node.xNext != nil {
		node.xPrev.xNext = node.xNext
		node.xNext.xPrev = node.xPrev
	} else if node.xPrev != nil {
		node.xPrev.xNext = nil
	}

	if node == l.yList {
		if node.yNext != nil {
			l.yList = node.yNext
			if l.yList.yPrev != nil {
				l.yList.yPrev = nil
			}
		} else {
			l.yList = nil
		}
	} else if node.yPrev != nil && node.yNext != nil {
		node.yPrev.yNext = node.yNext
		node.yNext.yPrev = node.yPrev
	} else if node.yPrev != nil {
		node.yPrev.yNext = nil
	}

	node.xPrev = nil
	node.xNext = nil
	node.yPrev = nil
	node.yNext = nil
}

// Enter event
func (l *ListAOI) EnterAOI(node *NodeAOI) {
	l.add(node)
	neighbors := l.findNeighbors(node)
	arr := make([]*NodeAOI, 0)
	for neighbor,_ := range neighbors {
		neighbor.listener.OnEnterAOI([]*NodeAOI{node})
		arr = append(arr, neighbor)
	}
	node.listener.OnEnterAOI(arr)
}

// Leave event
func (l *ListAOI) LeaveAOI(node *NodeAOI) {
	neighbors := l.findNeighbors(node)
	arr := make([]*NodeAOI, 0)
	for neighbor,_ := range neighbors {
		neighbor.listener.OnLeaveAOI([]*NodeAOI{node})
		arr = append(arr, neighbor)
	}
	node.listener.OnLeaveAOI(arr)
	l.remove(node)
}

// Move event
func (l *ListAOI) MoveAOI(node *NodeAOI, destX, destY float32) {
	oldNeighbors := l.findNeighbors(node)
	l.remove(node)
	node.x = destX
	node.y = destY
	l.add(node)
	newNeighbors := l.findNeighbors(node)
	arr := make([]*NodeAOI, 0)
	for old,_ := range oldNeighbors {
		if _,ok := newNeighbors[old]; !ok {
			old.listener.OnLeaveAOI([]*NodeAOI{node})
			arr = append(arr, old)
		}
	}
	node.listener.OnLeaveAOI(arr)
	arr = make([]*NodeAOI, 0)
	for _new,_ := range newNeighbors {
		if _,ok := oldNeighbors[_new]; !ok {
			_new.listener.OnEnterAOI([]*NodeAOI{node})
			arr = append(arr, _new)
		} else {
			_new.listener.OnUpdateAOI(node)
		}
	}
	node.listener.OnEnterAOI(arr)
}

// Find the watchers(neighbors)
func (l *ListAOI) findNeighbors(node *NodeAOI) map[*NodeAOI]struct{} {
	neighbors := make(map[*NodeAOI]struct{}, 0)

	l.lock.RLock()
	defer l.lock.RUnlock()

	for cur := node.xNext; cur != nil; cur = cur.xNext{
		if cur.x - node.x > node.r {
			break
		} else {
			if math.Abs(float64(cur.y - node.y)) <= float64(node.r) {
				neighbors[cur] = struct{}{}
			}
		}
	}

	for cur := node.xPrev; cur != nil; cur = cur.xPrev{
		if node.x - cur.x > node.r {
			break
		} else {
			if math.Abs(float64(cur.y - node.y)) <= float64(node.r) {
				neighbors[cur] = struct{}{}
			}
		}
	}
	return neighbors
}

func (l *ListAOI) FindNeighbors(node *NodeAOI, _range float32) map[*NodeAOI]struct{} {
	neighbors := make(map[*NodeAOI]struct{}, 0)

	l.lock.RLock()
	defer l.lock.RUnlock()

	for cur := node.xNext; cur != nil; cur = cur.xNext{
		if cur.x - node.x > _range {
			break
		} else {
			if math.Abs(float64(cur.y - node.y)) <= float64(_range) {
				neighbors[cur] = struct{}{}
			}
		}
	}

	for cur := node.xPrev; cur != nil; cur = cur.xPrev{
		if node.x - cur.x > _range {
			break
		} else {
			if math.Abs(float64(cur.y - node.y)) <= float64(_range) {
				neighbors[cur] = struct{}{}
			}
		}
	}
	return neighbors
}

// Debug only
func (l *ListAOI) Print() {
	for cur := l.xList; cur != nil; cur = cur.xNext{
		fmt.Printf("(%f,%f)->",cur.x,cur.y)
	}
	fmt.Println("x list end.")
	for cur := l.yList; cur != nil; cur = cur.yNext{
		fmt.Printf("(%f,%f)->",cur.x,cur.y)
	}
	fmt.Println("y list end.")
}
