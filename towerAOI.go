package aoi

import "fmt"

// A grid of the map
type tower struct {
	objs map[*NodeAOI]struct{}
}

func (t *tower) init() {
	t.objs = map[*NodeAOI]struct{}{}
}

func (t *tower) addObject(node *NodeAOI) {
	t.objs[node] = struct{}{}
}

func (t *tower) removeObject(node *NodeAOI) {
	delete(t.objs, node)
}

type TowerAOI struct {
	minX, maxX, minY, maxY float32
	towerRange             float32
	towers                 [][]tower
	xTowerNum, yTowerNum   int
}

// Create Tower AOI
func NewTowerAOI(minX, maxX, minY, maxY float32, towerRange float32) AOIAction {
	aoi := &TowerAOI{
		minX: minX,
		maxX: maxX,
		minY: minY,
		maxY: maxY,
		towerRange: towerRange,
	}
	aoi.xTowerNum = int((maxX - minX)/towerRange) + 1
	aoi.yTowerNum = int((maxY - minY)/towerRange) + 1
	for i := 0; i < aoi.xTowerNum; i++ {
		aoi.towers[i] = make([]tower, aoi.yTowerNum)
		for j := 0; j < aoi.yTowerNum; j++ {
			aoi.towers[i][j].init()
		}
	}
	return aoi
}

// The node is entering aoi
func (aoi *TowerAOI) EnterAOI(node *NodeAOI) {
	arr := make([]*NodeAOI, 0)
	neighbors := aoi.findNeighbors(node)
	for obj := range neighbors {
		obj.listener.OnEnterAOI([]*NodeAOI{node})
		arr = append(arr, obj)
	}

	t := aoi.getTowerByPos(node.x, node.y)
	t.addObject(node)
	node.listener.OnEnterAOI(arr)
}

// The node is leaving aoi
func (aoi *TowerAOI) LeaveAOI(node *NodeAOI) {
	t := aoi.getTowerByPos(node.x, node.y)
	t.removeObject(node)

	arr := make([]*NodeAOI, 0)
	neighbors := aoi.findNeighbors(node)
	for obj := range neighbors {
		obj.listener.OnLeaveAOI([]*NodeAOI{node})
		arr = append(arr, obj)
	}
	node.listener.OnLeaveAOI(arr)
}

// The node is moving to (destX,destY)
func (aoi *TowerAOI) MoveAOI(node *NodeAOI, destX, destY float32) {
	oldNeighbors := aoi.findNeighbors(node)
	t := aoi.getTowerByPos(node.x, node.y)
	t.removeObject(node)
	node.x = destX
	node.y = destY
	t = aoi.getTowerByPos(node.x, node.y)
	t.addObject(node)
	newNeighbors := aoi.findNeighbors(node)
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

// Get the specified _range neighbors
func (aoi *TowerAOI) FindNeighbors(node *NodeAOI, _range float32) map[*NodeAOI]struct{} {
	neighbors := make(map[*NodeAOI]struct{}, 0)

	xiMin, yiMin := aoi.transPosToTowerCoord(node.x - _range, node.y - _range)
	xiMax, yiMax := aoi.transPosToTowerCoord(node.x + _range, node.y + _range)
	for xi := xiMin; xi <= xiMax; xi++ {
		for yi := yiMin; yi <= yiMax; yi++ {
			t := &aoi.towers[xi][yi]
			for obj := range t.objs {
				neighbors[obj] = struct{}{}
			}
		}
	}

	return neighbors
}

// Debug only
func (aoi *TowerAOI) Print() {
	for x := 0; x < aoi.xTowerNum; x++ {
		for y := 0; y < aoi.yTowerNum; y++ {
			t := aoi.towers[x][y]
			for cur := range t.objs {
				fmt.Printf("(%f,%f)->",cur.x,cur.y)
			}
		}
	}
}

// Transfer position of world coordinate to tower coordinate.
func (aoi *TowerAOI) transPosToTowerCoord(x,y float32) (int, int) {
	xi := int((x - aoi.minX) / aoi.towerRange)
	yi := int((y - aoi.minY) / aoi.towerRange)

	if xi < 0 {
		xi = 0
	} else if xi >= aoi.xTowerNum {
		xi = aoi.xTowerNum - 1
	}

	if yi < 0 {
		yi = 0
	} else if yi >= aoi.yTowerNum {
		yi = aoi.yTowerNum - 1
	}
	return xi, yi
}

func (aoi *TowerAOI) getTowerByPos(x,y float32) *tower {
	xi, yi := aoi.transPosToTowerCoord(x, y)
	return &aoi.towers[xi][yi]
}

func (aoi *TowerAOI) findNeighbors(node *NodeAOI) map[*NodeAOI]struct{} {
	neighbors := make(map[*NodeAOI]struct{}, 0)

	xiMin, yiMin := aoi.transPosToTowerCoord(node.x - node.r, node.y - node.r)
	xiMax, yiMax := aoi.transPosToTowerCoord(node.x + node.r, node.y + node.r)
	for xi := xiMin; xi <= xiMax; xi++ {
		for yi := yiMin; yi <= yiMax; yi++ {
			t := &aoi.towers[xi][yi]
			for obj := range t.objs {
				neighbors[obj] = struct{}{}
			}
		}
	}

	return neighbors
}