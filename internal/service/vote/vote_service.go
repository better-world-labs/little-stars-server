package vote

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"errors"
	"github.com/go-xorm/xorm"
)

type voteService struct {
	Points service.PointsService `inject:"-"`
}

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
	return session.Table("vote_user_record").Where("vote_id = ? and user_id = ? and mode = ?", voteId, userId, entities.VoteRecordModeNormal).Count()
}

func (s *voteService) VoteNormal(voteId, userId int64, options []int64) error {
	vote, exists, err := s.GetVoteById(voteId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("vote not found")
	}

	options, err = s.checkAndProcessOptions(vote, options)
	if err != nil {
		return err
	}

	return s.doVoteNormal(vote, userId, options)
}

func (s *voteService) VotePoints(voteId, userId int64, options []int64) error {
	vote, exists, err := s.GetVoteById(voteId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("vote not found")
	}

	options, err = s.checkAndProcessOptions(vote, options)
	if err != nil {
		return err
	}

	return s.doVotePoints(voteId, userId, options)
}

func (s *voteService) checkAndProcessOptions(vote *entities.Vote, options []int64) ([]int64, error) {
	set := utils.NewSet[int64]()
	set.AddAll(options)
	options = set.ToSlice()

	// TODO 排序，防止死锁

	if len(options) == 0 {
		return nil, errors.New("no options set")
	}

	if vote.Status() != entities.VoteProjectStatusStarted {
		return nil, response.ErrorVoteCompleted
	}

	if vote.OptionType == entities.VoteOptionTypeSingle {
		return options[:1], nil
	}

	return options, nil
}

func (s *voteService) doVoteNormal(vote *entities.Vote, userId int64, options []int64) error {
	return db.Transaction(func(session *xorm.Session) error {
		remainTimes, err := s.GetUserRemainTimes(session, vote, userId)
		if err != nil {
			return err
		}

		if remainTimes < 1 {
			return response.ErrorNoVoteChance
		}

		for _, o := range options {
			err = s.vote(session, vote.Id, o)
			if err != nil {
				return err
			}
		}

		return s.SaveUserRecord(session, &entities.VoteRecord{
			UserId:    userId,
			VoteId:    vote.Id,
			OptionIds: options,
			Mode:      entities.VoteRecordModeNormal,
		})
	})
}

func (s *voteService) doVotePoints(voteId, userId int64, options []int64) error {
	return db.Transaction(func(session *xorm.Session) error {
		err := s.Points.Pay(userId, pkg.VoteCostPoints, pkg.VoteCostPointsDescription)
		if err != nil {
			return err
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
			Mode:      entities.VoteRecordModePoints,
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
