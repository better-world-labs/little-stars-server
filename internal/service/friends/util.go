package friends

func relationshipNodesSlice2Map(slice []*RelationshipTreeNode) map[int64]*RelationshipTreeNode {
	m := make(map[int64]*RelationshipTreeNode)
	for _, s := range slice {
		m[s.ID] = s
	}

	return m
}

func relationship2TreeNode(slice []*Relationship) []*RelationshipTreeNode {
	nodes := make([]*RelationshipTreeNode, len(slice))
	for i, s := range slice {
		nodes[i] = &RelationshipTreeNode{Relationship: *s}
	}
	return nodes
}

func groupingByParentId(slice []*Relationship) map[int64][]*Relationship {
	group := make(map[int64][]*Relationship)

	for _, s := range slice {
		group[s.ParentId] = append(group[s.ParentId], s)
	}

	return group
}

func mapUserIDs(slice []*Relationship) []int64 {
	ids := make([]int64, len(slice))

	for i := range slice {
		ids[i] = slice[i].UserId
	}

	return ids
}
