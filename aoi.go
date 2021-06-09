package aoi

// Observer will get AOI event
type AOIEvent interface {
	// The node is entering
	OnEnterAOI(nodes []*NodeAOI)

	// The node is moving
	OnUpdateAOI(node *NodeAOI)

	// The node is leaving
	OnLeaveAOI(nodes []*NodeAOI)
}

// AOI actions called by entity of the scene
type AOIAction interface {
	// The node is entering aoi
	EnterAOI(node *NodeAOI)

	// The node is leaving aoi
	LeaveAOI(node *NodeAOI)

	// The node is moving to (destX,destY)
	MoveAOI(node *NodeAOI, destX, destY float32)

	// Get the specified _range neighbors
	FindNeighbors(node *NodeAOI, _range float32) map[*NodeAOI]struct{}

	// Debug only
	Print()
}