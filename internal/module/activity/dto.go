package activity

type BorrowDevice struct {
	AidId int64 `json:"aidId,omitempty,string" binding:"required"`
}

type GoingToDevice struct {
	AidId int64 `json:"aidId,omitempty,string" binding:"required"`
}
