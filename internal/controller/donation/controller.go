package donation

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.openviewtech.com/openview-pub/gopkg/route"
	"strconv"
)

type Controller struct {
	serverDomain string
}

func NewController(domain string) *Controller {
	return &Controller{serverDomain: domain}
}

func (c Controller) MountAuthRouter(r *route.Router) {
	r.GET("/donations-donated", c.ListDonationsByDonator)
	r.GET("/donations/:id", c.GetDonation)
	r.POST("/donations/:id/records", c.Donate)
	r.GET("/donations/:id/records", c.ListLatestRecords)
	r.GET("/donations/:id/records/top", c.TopRecords)
	r.GET("/donations/:id/evidence", c.GetDonationEvidence)
	r.POST("/donations/apply", c.ApplyDonation)
	r.GET("/donation-honor", c.DonationHonor)
}

func (c Controller) MountNoAuthRouter(r *route.Router) {
	r.GET("/donations/apply/explain", c.ApplyExplain)
	r.GET("/donations", c.ListDonations)
}

func (c Controller) MountAdminRouter(r *route.Router) {
	donationGroup := r.Group("/donations")
	donationGroup.POST("", c.AdminCreateDonation)
	donationGroup.GET("/:id", c.AdminGetDonation)
}

func (c Controller) ListDonations(ctx *gin.Context) (interface{}, error) {
	var p page.Query
	err := ctx.ShouldBindQuery(&p)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	donations, err := interfaces.S.Donation.ListDonation(p)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"donations": DtosFromEntities(donations),
	}, nil
}

func (c Controller) ListDonationsByDonator(ctx *gin.Context) (interface{}, error) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)

	donations, err := interfaces.S.Donation.ListDonatorsDonation(userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"donations": WithDonatedDtosFromEntities(donations),
	}, nil
}

func (c Controller) GetDonation(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError("invalid param id")
	}

	donation, _, err := interfaces.S.Donation.GetDonationDetail(id)
	if err != nil {
		return nil, err
	}

	return DtoFromEntity(donation), nil
}

func (c Controller) Donate(ctx *gin.Context) (interface{}, error) {
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

func (c Controller) ListLatestRecords(ctx *gin.Context) (interface{}, error) {
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
		"records": RecordDtosFromEntities(records, users),
	}, nil
}

func (c Controller) TopRecords(ctx *gin.Context) (interface{}, error) {
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
		"records": RecordDtosFromEntities(records, users),
	}, nil
}

func (c Controller) GetDonationEvidence(ctx *gin.Context) (interface{}, error) {
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

	return EvidenceDto{
		ViewLink:         link,
		EvidenceImageUrl: fmt.Sprintf("%s/api/image-processing/resource/evidence?accountId=%d&category=3&businessKey=%d", c.serverDomain, latest.UserId, latest.Id),
	}, nil
}

func (c Controller) AdminCreateDonation(ctx *gin.Context) (interface{}, error) {
	var dto Dto
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}

	err = interfaces.S.Donation.CreateDonation(EntityFromDto(&dto))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c Controller) AdminGetDonation(ctx *gin.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, response.NewIllegalArgumentError(err.Error())
	}

	donation, _, err := interfaces.S.Donation.GetDonationById(id)
	return donation, nil
}

func (c Controller) ApplyDonation(ctx *gin.Context) (interface{}, error) {
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

func (c Controller) ApplyExplain(ctx *gin.Context) (interface{}, error) {
	query := ctx.Query("id")
	//todo
	return query, nil
}

func (c Controller) DonationHonor(ctx *gin.Context) (interface{}, error) {
	user := ctx.MustGet(pkg.AccountKey).(*entities.User)
	honor, err := interfaces.S.Donation.GetDonationHonor(user)
	if err != nil {
		return nil, err
	}

	return honor, nil
}
