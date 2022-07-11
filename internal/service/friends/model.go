package friends

import "time"

type Relationship struct {
	ID       int64 `xorm:"id pk autoincr"`
	UserId   int64
	ParentId int64
	CreateAt time.Time
}

type RelationshipTreeNode struct {
	Relationship `xorm:"extends"`

	children []*RelationshipTreeNode `xorm:"-"`
}

func ParseTrees(slice []*RelationshipTreeNode) ([]*RelationshipTreeNode, error) {
	nodeMap := relationshipNodesSlice2Map(slice)

	//TODO Loop Check
	var roots []*RelationshipTreeNode
	for _, node := range slice {
		if node.ParentId == 0 {
			roots = append(roots, node)
		}

		if parent, exists := nodeMap[node.ParentId]; exists {
			parent.AddChild(node)
		}
	}

	return roots, nil
}

func (n *RelationshipTreeNode) AddChild(child *RelationshipTreeNode) {
	n.children = append(n.children, child)
}

func (n *RelationshipTreeNode) GetChildren() []*RelationshipTreeNode {
	return n.children
}

type Prospective struct {
	ID       int64     `xorm:"id pk autoincr"`
	OpenId   string    `xorm:"open_id"`
	ParentId int64     `xorm:"parent_id"`
	CreateAt time.Time `xorm:"create_at"`
}
