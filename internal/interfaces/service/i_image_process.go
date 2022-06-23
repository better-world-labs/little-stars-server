package service

type IImageProcess interface {
	DrawDonationShareImage(recordId, userId int64) (string, error)
}
