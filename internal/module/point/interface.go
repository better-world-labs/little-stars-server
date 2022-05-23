package point

type Service interface {
	// Detail 积分明细
	// @Param accountID 用户id
	Detail(accountID int64) ([]*Point, error)

	// Total 总积分
	// @Param accountID 用户id
	TotalPoints(accountID int64) (float64, error)
}
