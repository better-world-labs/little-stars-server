package controller

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/asserts"
	"aed-api-server/internal/pkg/base"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/star"
	"aed-api-server/internal/pkg/utils"
	cert2 "aed-api-server/internal/service/cert"
	"aed-api-server/internal/service/imageprocessing"
	"aed-api-server/internal/service/img"
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ShareController TODO 后续整理代码
type ShareController struct {
	Process     service.IImageProcess  `inject:"-"`
	UserService service.UserServiceOld `inject:"-"`
	Skill       service.SkillService   `inject:"-"`

	ServerHost string `conf:"server.host"`
	UploadDir  string `conf:"alioss.upload-dir"`

	certCreator cert2.ImageCreator
}

//go:inject-component
func NewShareController() *ShareController {
	creator, err := cert2.NewImageCreatorDefaultAssert()
	if err != nil {
		panic(err)
	}
	return &ShareController{
		certCreator: creator,
	}
}

func (con *ShareController) MountNoAuthRouter(r *route.Router) {
	r.GET("image-processing/qrcode", con.RenderQrCode)
	r.GET("image-processing/share/essay", con.RenderShareEssay)
	r.GET("image-processing/share/donation", con.RenderShareDonation)
	r.GetOriginRouter().GET("image-processing/share/cert", con.RenderSharedCert)
	r.GetOriginRouter().GET("image-processing/share/medal", con.RenderSharedMedal)
	r.GetOriginRouter().GET("image-processing/resource/medal", con.RenderSharedMedal)
	r.GetOriginRouter().GET("image-processing/resource/cert", con.RenderResourceCert)
	r.GetOriginRouter().GET("image-processing/resource/evidence", con.RenderResourceEvidence)
}

func (con *ShareController) RenderSharedCert(c *gin.Context) {
	var certEntity *entities.UserCertEntity = nil
	var userId int64
	err := imageprocessing.LookUpAndGenPic(c, func() (string, string, *time.Time) {
		projectId, exists := c.GetQuery("projectId")
		utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))

		accountId, exists := c.GetQuery("accountId")
		utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))

		pId, err := strconv.ParseInt(projectId, 10, 64)
		userId, err = strconv.ParseInt(accountId, 10, 64)
		utils.MustNil(err, base.WrapError("imageprocessing", "invalid param", err))

		certEntity, exists, err = con.Skill.GetUserCertForProject(userId, pId)
		utils.MustTrue(exists, base.NewError("imageprocessing", "certEntity not exists"))
		return fmt.Sprintf("%s/redirect/cert/share_%s_%s.png", con.UploadDir, projectId, accountId), "", nil
	}, func(writer *io.PipeWriter) {
		err := con.DoRenderCert(certEntity.Img["origin"].(string), writer, userId)
		if err != nil {
			log.Errorf("gen pic err:%v", err)
		}
	})

	if err != nil {
		response.ReplyError(c, err)
	}
}

func helpTime() (today string, todayEnd time.Time) {
	now := time.Now()
	today = now.Format("2006-01-02")
	todayEnd, _ = time.ParseInLocation("2006-01-02 15:04:05", today+" 23:59:59", time.Local)
	return today, todayEnd
}

func (con *ShareController) RenderSharedMedal(c *gin.Context) {
	var mId int64
	var account *entities.User
	var err error
	medalId, exists := c.GetQuery("medalId")
	utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))

	accountId, exists := c.GetQuery("accountId")
	utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))
	err = imageprocessing.LookUpAndGenPic(c, func() (key string, url string, expired *time.Time) {

		mId, err = strconv.ParseInt(medalId, 10, 64)
		userId, err := strconv.ParseInt(accountId, 10, 64)
		utils.MustNil(err, base.WrapError("imageprocessing", "invalid param", err))

		account, err = con.UserService.GetUserByID(userId)
		utils.MustNil(err, err)

		today, todayEnd := helpTime()
		key = fmt.Sprintf("%s/redirect/medal/share_%s_%s", con.UploadDir, medalId, accountId)
		url = fmt.Sprintf("%s.png", key)
		key = fmt.Sprintf("%s_%s", key, today)
		return key, url, &todayEnd
	}, func(writer *io.PipeWriter) {
		err = imageprocessing.DrawMedalShare(mId, account, writer, con.ServerHost)
		if err != nil {
			log.Errorf("gen pic err:%v", err)
		}
	})

	if err != nil {
		response.ReplyError(c, err)
	}
}

func (con *ShareController) DoRenderCert(certUrl string, writer io.Writer, accountID int64) error {
	res, err := http.Get(certUrl)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("get avatar error for cert")
	}

	cert, _, err := image.Decode(res.Body)
	dst := image.NewRGBA(cert.Bounds())
	draw.Draw(dst, dst.Bounds(), cert, image.Point{}, draw.Over)

	err = img.DrawQrCodeRightBottom(dst, 128, star.GetPlaceCardSharedQrCodeContent(accountID), 10, 10)
	if err != nil {
		return err
	}

	return jpeg.Encode(writer, dst, &jpeg.Options{Quality: 90})
}

func (con *ShareController) RenderResourceCert(c *gin.Context) {
	var account *entities.User = nil
	var certEntity *entities.UserCertEntity = nil
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "public, max-age=86400")
	err := imageprocessing.LookUpAndGenPic(c, func() (string, string, *time.Time) {
		accountIdStr, exists := c.GetQuery("accountId")
		utils.MustTrue(exists, errors.New("renderResourceCert error: invalid param"))

		accountId, err := strconv.ParseInt(accountIdStr, 10, 64)
		utils.MustNil(err, err)

		account, err = con.UserService.GetUserByID(accountId)
		utils.MustNil(err, err)

		//TODO 这里暂时写死了projectId，后续改为直接传certId
		certEntity, exists, err = con.Skill.GetUserCertForProject(accountId, 1)
		if err != nil {
			response.ReplyError(c, err)
			return "", "", nil
		}
		if !exists {
			response.ReplyError(c, errors.New("证书不存在"))
			return "", "", nil
		}
		return fmt.Sprintf("%s/redirect/cert/resource_%s_%v.png", con.UploadDir, "account", accountId), "", nil
	}, func(writer *io.PipeWriter) {
		err := con.certCreator.Create(account.Avatar, account.Nickname, "\"茫茫人海之中，去挽救下一个倒地昏迷的人吧\"", time.Time(certEntity.Created), writer)
		if err != nil {
			log.Errorf("gen pic err:%v", err)
		}
	})

	if err != nil {
		response.ReplyError(c, err)
	}
}

func (con *ShareController) DoRenderResourceCert(account *entities.User, entity *entities.UserCertEntity, writer io.Writer) error {
	return con.certCreator.Create(account.Avatar, account.Nickname, "\"茫茫人海之中，去挽救下一个倒地昏迷的人吧\"", time.Time(entity.Created), writer)
}

func (con *ShareController) RenderResourceEvidence(c *gin.Context) {
	var account *entities.User = nil
	var evi *entities.Evidence = nil
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "public, max-age=86400")
	err := imageprocessing.LookUpAndGenPic(c, func() (string, string, *time.Time) {
		category, exists := c.GetQuery("category")
		utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))

		businessKey, exists := c.GetQuery("businessKey")
		utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))

		cate, err := strconv.ParseInt(category, 10, 64)
		utils.MustTrue(exists, base.NewError("imageprocessing", "invalid param"))

		evi, exists, err = interfaces.S.Evidence.GetEvidenceByBusinessKey(businessKey, entities.EvidenceCategory(cate))
		utils.MustNil(err, base.NewError("imageprocessing", "get evidence error"))
		utils.MustTrue(exists, base.NewError("imageprocessing", "not found"))

		account, err = con.UserService.GetUserByID(evi.AccountID)
		utils.MustNil(err, base.NewError("imageprocessing", "get account error"))
		return fmt.Sprintf("%s/redirect/evidence/%s_%v.png", con.UploadDir, category, businessKey), "", nil
	}, func(writer *io.PipeWriter) {
		err := con.DoRenderResourceEvidence(account, evi, writer)
		if err != nil {
			log.Errorf("gen pic err:%v", err)
		}
	})

	if err != nil {
		response.ReplyError(c, err)
	}
}

// DoRenderResourceEvidenceSimpleCategoryCertOrMedal  证书与勋章可以复用存证结构
func (con *ShareController) DoRenderResourceEvidenceSimpleCategoryCertOrMedal(account *entities.User, evi *entities.Evidence, writer io.Writer) error {
	bgBytes, _ := asserts.GetResource("evidence_background_simple.jpg")

	bgImg, _, err := image.Decode(bytes.NewReader(bgBytes))
	bg := image.NewRGBA(bgImg.Bounds())
	sealBytes, _ := asserts.GetResource("evidence_seal.png")

	seal, _, err := image.Decode(bytes.NewReader(sealBytes))
	if err != nil {
		return err
	}
	draw.Draw(bg, bg.Bounds(), bgImg, image.Point{}, draw.Over)

	err = img.DrawText(bg, evi.TransactionHash, 450, 922, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, evi.Time.Format("2006-01-02 15:04:05"), 390, 1101, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, account.Nickname, 390, 1175, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, "HASH 上链", 390, 1323, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, fmt.Sprintf("%d 字节", evi.FileBytes), 430, 1454, 22, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, evi.ContentHash, 390, 1624, 22, color.Black)
	if err != nil {
		return err
	}

	content := fmt.Sprintf("https://openscan.openviewtech.com/#/transaction/transactionDetail?pageSize=10&pageNumber=1&v_page=transaction&pkHash=%s", evi.TransactionHash)

	err = img.DrawQrCodeRightBottom(bg, 234, content, 172, 274)
	if err != nil {
		return err
	}

	resizedSeal := imaging.Resize(seal, 500, 500, imaging.Lanczos)
	draw.Draw(bg, bg.Bounds(), resizedSeal, image.Point{X: -1000, Y: -1040}, draw.Over)
	return jpeg.Encode(writer, bg, &jpeg.Options{Quality: 90})
}

func (con *ShareController) DoRenderResourceEvidence(account *entities.User, evi *entities.Evidence, writer io.Writer) error {
	switch evi.Category {
	case entities.EvidenceCategoryCert, entities.EvidenceCategoryMedal:
		return con.DoRenderResourceEvidenceSimpleCategoryCertOrMedal(account, evi, writer)

	case entities.EvidenceCategoryDonation:
		return con.DoRenderResourceEvidenceSimpleDonation(account, evi, writer)

	default:
		return errors.New("invalid evidence type")
	}
}

// DoRenderResourceEvidenceSimpleDonation  数据结构与其他两种不一样
func (con *ShareController) DoRenderResourceEvidenceSimpleDonation(account *entities.User, evi *entities.Evidence, writer io.Writer) error {
	bgBytes, _ := asserts.GetResource("evidence_background_first_donation.jpg")

	bgImg, _, err := image.Decode(bytes.NewReader(bgBytes))
	bg := image.NewRGBA(bgImg.Bounds())
	sealBytes, _ := asserts.GetResource("evidence_seal.png")

	seal, _, err := image.Decode(bytes.NewReader(sealBytes))
	if err != nil {
		return err
	}
	draw.Draw(bg, bg.Bounds(), bgImg, image.Point{}, draw.Over)

	recordId, err := strconv.ParseInt(evi.BusinessKey, 10, 64)
	if err != nil {
		return err
	}

	record, exists, err := interfaces.S.Donation.GetRecordById(recordId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("record not found")
	}

	donation, exists, err := interfaces.S.Donation.GetDonationDetail(record.DonationId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("donation not found")
	}

	err = img.DrawText(bg, fmt.Sprintf("【%s】", donation.Title), 155, 922, 28, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, "项目的所有捐积分记录都会上传至区块链，区块链由多方共同维护，永久保护，任何人", 155, 980, 28, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, "无法篡改，确保相关信息的公开，透明。", 155, 1038, 30, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, fmt.Sprintf("截止%s，已上传%s笔积分流水，共计%s积分。", evi.Time.Format("2006年01月02日"), utils.PointsString(*donation.RecordsCount), utils.PointsString(donation.ActualPoints)), 155, 1154, 28, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, "小星星区块链平台", 371, 1409, 28, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawText(bg, "0x0c", 371, 1530, 28, color.Black)
	if err != nil {
		return err
	}

	link := interfaces.S.Evidence.GetTransactionViewLink(evi.TransactionHash)

	err = img.DrawQrCodeRightBottom(bg, 234, link, 172, 274)
	if err != nil {
		return err
	}

	err = img.DrawTextAutoBreakASCII(bg, link, 350, 1655, 50, 30, 28, color.Black)
	if err != nil {
		return err
	}

	resizedSeal := imaging.Resize(seal, 500, 500, imaging.Lanczos)
	draw.Draw(bg, bg.Bounds(), resizedSeal, image.Point{X: -1000, Y: -1040}, draw.Over)
	return jpeg.Encode(writer, bg, &jpeg.Options{Quality: 90})
}

func (con *ShareController) RenderShareDonation(c *gin.Context) (interface{}, error) {
	var param struct {
		RecordId  int64 `form:"recordId"`
		AccountId int64 `form:"accountId"`
	}

	err := c.ShouldBindQuery(&param)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	url, err := con.Process.DrawDonationShareImage(param.RecordId, param.AccountId)
	if err != nil {
		return nil, err
	}

	c.Redirect(302, url)

	return nil, nil
}

func (con *ShareController) RenderShareEssay(ctx *gin.Context) (interface{}, error) {

	var err error
	var query = struct {
		Source  string `form:"source" binding:"required"`
		EssayId int64  `form:"essayId" binding:"required"`
		Url     string `form:"url" binding:"required"`
		Sharer  int64  `form:"sharer" binding:"required"`
	}{}

	err = ctx.ShouldBindQuery(&query)
	utils.MustNil(err, err)

	// TODO 缓存还是有问题，需要做调整，这里先这样
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Cache-Control", "public, max-age=86400")
	err = DrawEssayShare(query.EssayId, query.Source, query.Sharer, query.Url, ctx.Writer)
	//err = lookUpAndGenPic(ctx, func() (key string, url string, expired *time.Time) {
	//	today, todayEnd := helpTime()
	//	key = fmt.Sprintf("%s/redirect/medal/share_%s_%d", con.uploadDir, query.Source, query.Sharer)
	//	url = fmt.Sprintf("%s.png", key)
	//	key = fmt.Sprintf("%s_%s", key, today)
	//	return key, url, &todayEnd
	//}, func(writer *io.PipeWriter) {
	//	if err != nil {
	//		log.Errorf("gen pic err:%v", err)
	//	}
	//})

	return nil, err
}

func (con ShareController) RenderQrCode(c *gin.Context) (interface{}, error) {
	param := struct {
		Content string `form:"content" binding:"required,max=128"`
		Size    int    `form:"size" binding:"required"`
	}{}

	err := c.ShouldBindQuery(&param)
	if err != nil {
		return nil, err
	}

	code, err := qrcode.New(param.Content, qrcode.Medium)
	if err != nil {
		c.Status(400)
		_, _ = c.Writer.Write([]byte(err.Error()))
		return nil, err
	}
	code.DisableBorder = true

	c.Header("Cache-Control", "public, max-age=86400")
	err = code.Write(param.Size, c.Writer)
	if err != nil {
		log.Errorf("gen pic err:%v", err)
	}

	return nil, nil
}

func DrawEssayShare(essayId int64, source string, sharer int64, _url string, writer io.Writer) error {
	essay, err := interfaces.S.Essay.GetById(essayId)
	if err != nil {
		return err
	}

	if essay == nil {
		return errors.New("essay not found")
	}

	var bgBytes []byte
	if essay.Extra == "" {
		return drawEssayShareDefault(essay, source, sharer, _url, writer)
	}

	bgBytes, err = asserts.GetResourceFromUrl(essay.Extra)
	if err != nil {
		return err
	}

	bgImage, _, err := image.Decode(bytes.NewReader(bgBytes))
	if err != nil {
		return err
	}
	background := image.NewRGBA(bgImage.Bounds())
	draw.Draw(background, background.Bounds(), bgImage, image.Point{}, draw.Over)

	size := computeQrCodeSize(background.Bounds())
	margin := computeMargin(size)
	err = img.DrawQrCodeRightBottom(background, size, star.GetEssaySharedQrCodeContent(sharer, url.QueryEscape(_url), source), margin, margin)
	if err != nil {
		return err
	}

	return jpeg.Encode(writer, background, &jpeg.Options{Quality: 90})
}

func drawEssayShareDefault(essay *entities.Essay, source string, sharer int64, _url string, writer io.Writer) error {
	bgBytes, _ := asserts.GetResource("essay_default_share_bg.png")

	bgImage, _, err := image.Decode(bytes.NewReader(bgBytes))
	if err != nil {
		return err
	}

	background := image.NewRGBA(bgImage.Bounds())
	draw.Draw(background, background.Bounds(), bgImage, image.Point{}, draw.Over)

	err = img.DrawTextAutoBreakRune(background, []rune(fmt.Sprintf("《%s》", essay.Title)), 130, 944, 16, 25, 38, color.RGBA{R: 9, G: 109, B: 235, A: 255})
	if err != nil {
		return err
	}

	err = img.DrawText(background, "长按查看文章完整版", 150, 1250, 25, color.Black)
	if err != nil {
		return err
	}

	err = img.DrawQrCodeRightBottom(background, 200, star.GetEssaySharedQrCodeContent(sharer, url.QueryEscape(_url), source), 100, 130)
	if err != nil {
		return err
	}

	return jpeg.Encode(writer, background, &jpeg.Options{Quality: 100})
}

func computeQrCodeSize(bg image.Rectangle) int {
	return int(math.Round(float64(bg.Max.X-bg.Min.X) / 4))
}

func computeMargin(size int) int {
	return size / 8
}
