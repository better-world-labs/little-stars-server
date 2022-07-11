package game

import (
	"aed-api-server/internal/interfaces/domains"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/utils"
	"math"
)

func playersMapUserIds(players []*entities.GameProcess) []int64 {
	set := utils.NewInt64Set()

	for _, p := range players {
		set.Add(p.UserId)
	}

	return set.ToSlice()
}

func GameProcessAddUser(process *domains.GameProcess, user *entities.SimpleUser) {
	process.User = user
}

func GameProcessesAddUsers(processes []*domains.GameProcess, users map[int64]*entities.SimpleUser) {
	for _, p := range processes {
		GameProcessAddUser(p, users[p.UserId])
	}
}

func caculateGamePoints(completedCount, totalPoints int) int {
	if completedCount == 0 {
		return totalPoints
	}

	return int(math.Round(float64(totalPoints) / float64(completedCount)))
}

func caculateUserRankPercents(userId int64, rankedUsers []int64) int {
	rank := len(rankedUsers)

	for i, u := range rankedUsers {
		if u == userId {
			rank = i + 1
			break
		}
	}

	percent := int(math.Round(float64(rank) / float64(len(rankedUsers)) * 100))
	if percent > 99 {
		percent = 99
	}

	if percent < 1 {
		percent = 1
	}

	return percent
}
