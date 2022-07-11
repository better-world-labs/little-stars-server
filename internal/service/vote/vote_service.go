package vote

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"github.com/go-xorm/xorm"
)

type voteService struct{}

//go:inject-component
func NewVoteService() service.VoteService {
	return &voteService{}
}

func (s *voteService) GetVoteById(id int64) (*entities.Vote, bool, error) {
	var v entities.Vote
	exists, err := db.GetById("vote", id, &v)
	return &v, exists, err
}

func (s *voteService) GetVoteOptionById(id int64) (*entities.VoteOptionDetail, bool, error) {
	var v entities.VoteOption
	exists, err := db.GetById("vote_option", id, &v)
	if err != nil {
		return nil, exists, err
	}

	if !exists {
		return nil, exists, nil
	}

	detail := &entities.VoteOptionDetail{VoteOption: v}
	rank, err := s.ListVoteOptionsRank(v.VoteId)
	if err != nil {
		return nil, exists, err
	}

	for i, o := range rank {
		if o.Id == v.Id {
			detail.Rank = i + 1
		}
	}

	return detail, exists, err
}

func (s *voteService) GetVoteOptionByIdForUpdate(session *xorm.Session, id int64) (*entities.VoteOption, bool, error) {
	var v entities.VoteOption
	exists, err := session.Table("vote_option").Where("id = ?", id).ForUpdate().Get(&v)
	return &v, exists, err
}

func (s *voteService) ListVoteOptions(voteId int64) ([]*entities.VoteOption, error) {
	var opts []*entities.VoteOption
	err := db.Table("vote_option").Where("vote_id = ?", voteId).Find(&opts)
	return opts, err
}

func (s *voteService) ListVoteOptionsRank(voteId int64) ([]*entities.VoteOption, error) {
	var opts []*entities.VoteOption
	err := db.Table("vote_option").Where("vote_id = ?", voteId).Desc("vote").Find(&opts)
	return opts, err
}

func (s *voteService) CountUserRecords(session *xorm.Session, voteId, userId int64) (int64, error) {
	return session.Table("vote_user_record").Where("vote_id = ? and user_id = ?", voteId, userId).Count()
}

func (s *voteService) DoVote(voteId, userId int64, options []int64) error {
	set := utils.NewInt64Set()
	set.AddAll(options)
	options = set.ToSlice()

	// TODO 排序，防止死锁

	if len(options) == 0 {
		return errors.New("no options set")
	}

	vote, exists, err := s.GetVoteById(voteId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no vote found")
	}

	if vote.Status() != entities.VoteProjectStatusStarted {
		return response.ErrorVoteCompleted
	}

	if vote.OptionType == entities.VoteOptionTypeSingle {
		options = options[:1]
	}

	return db.Transaction(func(session *xorm.Session) error {
		remainTimes, err := s.GetUserRemainTimes(session, vote, userId)
		if err != nil {
			return err
		}

		if remainTimes < 1 {
			return response.ErrorNoVoteChance
		}

		for _, o := range options {
			err = s.vote(session, voteId, o)
			if err != nil {
				return err
			}
		}

		return s.SaveUserRecord(session, &entities.VoteRecord{
			UserId:    userId,
			VoteId:    voteId,
			OptionIds: options,
		})
	})
}

func (s *voteService) vote(session *xorm.Session, voteId, optionId int64) error {
	option, exists, err := s.GetVoteOptionByIdForUpdate(session, optionId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no option found")
	}

	if voteId != option.VoteId {
		return errors.New("invalid optionId for vote")
	}

	option.Vote++

	return s.UpdateOption(session, option)
}

func (s *voteService) UpdateOption(session *xorm.Session, option *entities.VoteOption) error {
	_, err := session.Table("vote_option").ID(option.Id).Update(option)
	return err
}

func (s *voteService) SaveUserRecord(session *xorm.Session, record *entities.VoteRecord) error {
	_, err := session.Table("vote_user_record").Insert(record)
	return err
}

func (s *voteService) GetUserRemainTimes(session *xorm.Session, vote *entities.Vote, userId int64) (int, error) {
	timesVoted, err := s.CountUserRecords(session, vote.Id, userId)
	if err != nil {
		return 0, err
	}

	return vote.MaxTimes - int(timesVoted), nil
}
