package essay

import (
	"aed-api-server/internal/interfaces/entities"
)

func userSliceToMap(users []*entities.SimpleUser) map[int64]*entities.SimpleUser {
	m := make(map[int64]*entities.SimpleUser)

	for _, u := range users {
		m[u.ID] = u
	}

	return m
}
