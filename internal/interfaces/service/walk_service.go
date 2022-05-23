package service

import "aed-api-server/internal/interfaces/entities"

type WalkConvertInfoReq struct {
	Code          string `json:"code" form:"code"`
	Iv            string `json:"iv" form:"iv"`
	EncryptedData string `json:"encryptedData" form:"encryptedData"`
}

type WalkConvertInfo struct {
	TodayWalk       int `json:"todayWalk"`
	UnConvertWalk   int `json:"unConvertWalk"`
	ConvertedPoints int `json:"convertedPoints"`
	ConvertRatio    int `json:"convertRatio"`
}

type ConvertWalkToPointsRst struct {
	UnConvertWalk        int `json:"unConvertWalk"`
	ConvertedPoints      int `json:"convertedPoints"`
	CurrentConvertPoints int `json:"currentConvertPoints"`

	DealPointsRst *entities.DealPointsEventRst `json:"dealPointsRst"`
}

type WalkService interface {
	//GetWalkConvertInfo 获取积分兑换信息
	GetWalkConvertInfo(userId int64, req *WalkConvertInfoReq) (*WalkConvertInfo, error)

	//ConvertWalkToPoints 兑换积分
	ConvertWalkToPoints(userId int64, todayWalk int) (*ConvertWalkToPointsRst, error)
}
