package entities

import (
	"aed-api-server/internal/pkg/response"
	"time"
)

const (
	StatusNotStarted = 0
	StatusIng        = 1
	StatusCompleted  = 2
	StatusExpired    = 3
)

type (
	DonationRecord struct {
		Id         int64
		DonationId int64 `json:"donationId"`
		Points     int   `json:"points"`
		UserId     int64 `json:"userId"`
	}

	DonationStat struct {
		//积分捐出的总积分
		DonationTotalPoints int

		//积分捐献项目总数
		DonationProjectCount int

		//积分捐献次数
		DonationCount int
	}

	Donation struct {
		Id             int64      `xorm:"id pk autoincr" json:"id"`
		Title          string     `xorm:"title" json:"title"`
		Images         []string   `xorm:"images" json:"images"`
		Description    string     `xorm:"description" json:"description"`
		TargetPoints   int        `xorm:"target_points" json:"targetPoints"`
		ActualPoints   int        `xorm:"actual_points" json:"actualPoints"`
		StartAt        time.Time  `xorm:"start_at" json:"-"`
		CompleteAt     *time.Time `xorm:"complete_at" json:"-"`
		ExpiredAt      time.Time  `xorm:"expired_at" json:"-"`
		Status         int        `xorm:"-" json:"status"`
		ArticleId      int64      `xorm:"article_id" json:"articleId"`
		Executor       string     `xorm:"executor" json:"executor"`
		ExecutorNumber string     `xorm:"executor_number" json:"executorNumber"`
		Feedback       string     `xorm:"feedback" json:"feedback"`
		Plan           string     `xorm:"plan" json:"plan"`
		PlanImage      string     `xorm:"plan_image" json:"planImage"`
		Budget         string     `xorm:"budget" json:"budget"`
		CreatedAt      time.Time  `xorm:"created_at" json:"createdAt"`
		RecordsCount   *int       `xorm:"-" json:"recordsCount,omitempty"`
	}

	DonationWithUserDonated struct {
		Donation

		DonatedPoints int `json:"donatedPoints"`
	}

	DonationEvidence struct {
		ViewLink         string `json:"viewLink,omitempty"`
		EvidenceImageUrl string `json:"evidenceImageUrl,omitempty"`
	}

	DonationApply struct {
		Name   string `json:"name" xorm:"name"`
		Sex    string `json:"sex" xorm:"sex"`
		Region string `json:"region" xorm:"region"`
		Mobile string `json:"mobile" xorm:"mobile"`
		Job    string `json:"job" xorm:"job"`
	}
)

func NewDonation(
	id int64,
	title string,
	images []string,
	description string,
	targetPoints int,
	actualPoints int,
	startAt time.Time,
	completeAt *time.Time,
	expiredAt time.Time,
	articleId int64,
	executor string,
	executorNumber string,
	feedback string,
	status int,
	plan string,
	planImage string,
	budget string,
) *Donation {
	return &Donation{
		Id:             id,
		Title:          title,
		Images:         images,
		Description:    description,
		TargetPoints:   targetPoints,
		ActualPoints:   actualPoints,
		StartAt:        startAt,
		CompleteAt:     completeAt,
		ExpiredAt:      expiredAt,
		Executor:       executor,
		ExecutorNumber: executorNumber,
		Feedback:       feedback,
		Status:         status,
		ArticleId:      articleId,
		Plan:           plan,
		PlanImage:      planImage,
		Budget:         budget,
	}
}

func (d *Donation) Donate(points int) (int, error) {
	if d.Status == StatusNotStarted {
		return 0, response.ErrorDonationNotStartYet
	}

	if d.Status == StatusCompleted {
		return 0, response.ErrorDonationCompleted
	}

	if d.Status == StatusExpired {
		return 0, response.ErrorDonationExpired
	}

	left := d.TargetPoints - d.ActualPoints
	need := min(left, points)
	d.ActualPoints += need
	if d.CheckCompleted() {
		d.Complete()
	}

	return need, nil
}

func (d *Donation) CheckCompleted() bool {
	return d.ActualPoints == d.TargetPoints
}

func (d *Donation) Complete() {
	d.Status = StatusCompleted
	//now := global.FormattedTime(time.Now())
	n := time.Now()
	d.CompleteAt = &n
}

func (d *Donation) GetProcessPercents() float32 {
	return float32(d.ActualPoints) / float32(d.TargetPoints) * 100
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
