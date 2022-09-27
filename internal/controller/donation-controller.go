package controller

import (
	dto2 "aed-api-server/internal/controller/dto"
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
)

type DonationController struct {
	ServerDomain string              `conf:"server.host"`
	User         service.UserService `inject:"-"`
}

//go:inject-component
func NewDonationController() *DonationController {
	return &DonationController{}
}

func (c DonationController) MountAuthRouter(r *route.Router) {
	r.GET("/donations-donated", c.ListDonationsByDonator)
	r.GET("/donations/:id", c.GetDonation)
	r.POST("/donations/:id/records", c.Donate)
	r.GET("/donations/:id/records", c.ListLatestRecords)
	r.GET("/donations/:id/records/top", c.TopRecords)
	r.GET("/donations/:id/evidence", c.GetDonationEvidence)
	r.POST("/donations/apply", c.ApplyDonation)
	r.GET("/donation-honor", c.DonationHonor)
}

func (c DonationController) MountNoAuthRouter(r *route.Router) {
	r.GET("/donations/apply/explain", c.ApplyExplain)
	r.GET("/donations", c.ListDonations)
}

func (c DonationController) MountAdminRouter(r *route.Router) {
	donationGroup := r.Group("/donations")
	donationGroup.GET("", c.AdminListDonations)
	donationGroup.GET("/:id", c.AdminGetDonation)
	donationGroup.POST("", c.AdminCreateDonation)
	r.PUT("/donations/:id/crowdfunding", c.AdminUpdateCrowdfunding)
}

func (c DonationController) AdminListDonations(ctx *gin.Context) (interface{}, error) {
	donations, err := interfaces.S.Donation.ListDonation()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"donations": dto2.DtosFromEntities(donations),
	}, nil
}
func (c DonationController) ListDonations(ctx *gin.Context) (interface{}, error) {
	var p page.Query
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	token := ctx.GetHeader(pkg.AuthorizationHeaderKey)
	var userId int64
	user, err := c.User.ParseInfoFromJwtToken(token)
	if err == nil {
		userId = user.ID
	}

	donations, err := interfaces.S.Donation.ListDonationSorted(p, userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"donations": dto2.DtosFromEntities(donations),
	}, nil
}

func (c DonationController) ListDonationsByDonator(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)

	donations, err := interfaces.S.Donation.ListDonorsDonation(userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"donations": dto2.WithDonatedDtosFromEntities(donations),
	}, nil
}

func (c DonationController) GetDonation(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid param id")
	}

	donation, _, err := interfaces.S.Donation.GetDonationDetail(id)
	if err != nil {
		return nil, err
	}

	return dto2.DtoFromEntity(donation), nil
}

func (c DonationController) Donate(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid param id")
	}

	points := struct {
		Points int `json:"points" bind:"require,min=1"`
	}{}

	err = ctx.ShouldBindJSON(&points)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	record, err := interfaces.S.Donation.Donate(&entities.DonationRecord{DonationId: id, UserId: userId, Points: points.Points})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"recodeId": record.Id,
		"points":   record.Points,
	}, nil
}

func (c DonationController) ListLatestRecords(ctx *gin.Context) (interface{}, error) {
	query := struct {
		Latest int `form:"latest"`
	}{}
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	donationIdStr := ctx.Param("id")
	donationid, err := strconv.ParseInt(donationIdStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	records, err := interfaces.S.Donation.ListRecords(donationid, query.Latest)
	if err != nil {
		return nil, err
	}

	var userIds []int64
	for _, r := range records {
		userIds = append(userIds, r.UserId)
	}
	users, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"records": dto2.RecordDtosFromEntities(records, users),
	}, nil
}

func (c DonationController) TopRecords(ctx *gin.Context) (interface{}, error) {
	query := struct {
		Size int `form:"size"`
	}{}
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	donationIdStr := ctx.Param("id")
	donationid, err := strconv.ParseInt(donationIdStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	records, err := interfaces.S.Donation.ListUsersRecordsTop(donationid, query.Size)
	if err != nil {
		return nil, err
	}

	var userIds []int64
	for _, r := range records {
		userIds = append(userIds, r.UserId)
	}
	users, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"records": dto2.RecordDtosFromEntities(records, users),
	}, nil
}

func (c DonationController) GetDonationEvidence(ctx *gin.Context) (interface{}, error) {
	donationIdStr := ctx.Param("id")
	donationId, err := strconv.ParseInt(donationIdStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError("donation not found")
	}

	records, err := interfaces.S.Donation.ListRecords(donationId, 1)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, err
	}

	latest := records[0]

	link, err := interfaces.S.Evidence.GetTransactionViewLinkByBusinessKey(strconv.FormatInt(latest.Id, 10), entities.EvidenceCategoryDonation)
	if err != nil {
		return nil, err
	}

	return dto2.EvidenceDto{
		ViewLink:         link,
		EvidenceImageUrl: fmt.Sprintf("https://%s/api/image-processing/resource/evidence?accountId=%d&category=3&businessKey=%d", c.ServerDomain, latest.UserId, latest.Id),
	}, nil
}

func (c DonationController) AdminCreateDonation(ctx *gin.Context) (interface{}, error) {
	var dto dto2.Dto
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Donation.CreateDonation(dto2.EntityFromDto(&dto))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c DonationController) AdminGetDonation(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	donation, _, err := interfaces.S.Donation.GetDonationById(id)
	return donation, nil
}

func (c DonationController) ApplyDonation(ctx *gin.Context) (interface{}, error) {
	var apply entities.DonationApply
	if err := ctx.ShouldBindJSON(&apply); err != nil {
		return nil, err
	}

	err := interfaces.S.Donation.Apply(apply, ctx.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c DonationController) ApplyExplain(ctx *gin.Context) (interface{}, error) {
	query := ctx.Query("id")
	//todo
	return query, nil
}

func (c DonationController) DonationHonor(ctx *gin.Context) (interface{}, error) {
	user := ctx.MustGet(pkg.AccountKey).(*entities.User)
	honor, err := interfaces.S.Donation.GetDonationHonor(user)
	if err != nil {
		return nil, err
	}

	return honor, nil
}

func (c DonationController) AdminUpdateCrowdfunding(ctx *gin.Context) (interface{}, error) {
	var param struct {
		ActualCrowdfunding float32 `json:"actualCrowdfunding" binding:"required"`
	}

	id, err := utils.GetContextPathParamInt64(ctx, "id")
	if err != nil {
		return nil, err
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, err
	}

	err = interfaces.S.Donation.UpdateCrowdfunding(id, param.ActualCrowdfunding)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
