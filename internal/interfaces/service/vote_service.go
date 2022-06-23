package service

import (
	"aed-api-server/internal/interfaces/entities"
	"github.com/go-xorm/xorm"
)

type VoteService interface {
	GetVoteById(id int64) (*entities.Vote, bool, error)

	GetVoteOptionById(id int64) (*entities.VoteOptionDetail, bool, error)

	ListVoteOptions(voteId int64) ([]*entities.VoteOption, error)

	ListVoteOptionsRank(voteId int64) ([]*entities.VoteOption, error)

	DoVote(voteId, userId int64, options []int64) error

	GetUserRemainTimes(session *xorm.Session, vote *entities.Vote, userId int64) (int, error)
}
