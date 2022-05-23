package donation

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg"
	page "aed-api-server/internal/pkg/query"
	"aed-api-server/internal/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Controller struct {
	serverDomain string
}

func NewController(domain string) *Controller {
	return &Controller{serverDomain: domain}
}

func (c Controller) ListDonations(ctx *gin.Context) {
	var p page.Query
	err := ctx.BindQuery(&p)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, err.Error())
		return
	}

	donations, err := interfaces.S.Donation.ListDonation(p)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"donations": DtosFromEntities(donations),
	})
}

func (c Controller) ListDonationsByDonator(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)

	donations, err := interfaces.S.Donation.ListDonatorsDonation(userId)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"donations": WithDonatedDtosFromEntities(donations),
	})
}

func (c Controller) GetDonation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, "invalid param id")
		return
	}

	donation, _, err := interfaces.S.Donation.GetDonationDetail(id)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}
	response.ReplyOK(ctx, DtoFromEntity(donation))
}

func (c Controller) Donate(ctx *gin.Context) {
	userId := ctx.MustGet(pkg.AccountIDKey).(int64)
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, "invalid param id")
		return
	}

	points := struct {
		Points int `json:"points" bind:"require,min=1"`
	}{}

	err = ctx.BindJSON(&points)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, "invalid param id")
		return
	}

	donation, err := interfaces.S.Donation.Donate(&entities.DonationRecord{DonationId: id, UserId: userId, Points: points.Points})
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, DtoFromEntity(donation))
}

func (c Controller) ListLatestRecords(ctx *gin.Context) {
	query := struct {
		Latest int `form:"latest"`
	}{}
	err := ctx.BindQuery(&query)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, err.Error())
		return
	}

	donationIdStr := ctx.Param("id")
	donationid, err := strconv.ParseInt(donationIdStr, 10, 64)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, err.Error())
		return
	}

	records, err := interfaces.S.Donation.ListRecords(donationid, query.Latest)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}
	var userIds []int64
	for _, r := range records {
		userIds = append(userIds, r.UserId)
	}
	users, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"records": RecordDtosFromEntities(records, users),
	})
}

func (c Controller) TopRecords(ctx *gin.Context) {
	query := struct {
		Size int `form:"size"`
	}{}
	err := ctx.BindQuery(&query)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, err.Error())
		return
	}

	donationIdStr := ctx.Param("id")
	donationid, err := strconv.ParseInt(donationIdStr, 10, 64)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, err.Error())
		return
	}

	records, err := interfaces.S.Donation.ListUsersRecordsTop(donationid, query.Size)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	var userIds []int64
	for _, r := range records {
		userIds = append(userIds, r.UserId)
	}
	users, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, map[string]interface{}{
		"records": RecordDtosFromEntities(records, users),
	})
}

func (c Controller) GetDonationEvidence(ctx *gin.Context) {
	donationIdStr := ctx.Param("id")
	donationId, err := strconv.ParseInt(donationIdStr, 10, 64)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, "donation not found")
		return
	}

	records, err := interfaces.S.Donation.ListRecords(donationId, 1)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}
	if len(records) == 0 {
		response.ReplyError(ctx, err)
		return
	}

	latest := records[0]

	link, err := interfaces.S.Evidence.GetTransactionViewLinkByBusinessKey(strconv.FormatInt(latest.Id, 10), entities.EvidenceCategoryDonation)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, EvidenceDto{
		ViewLink:         link,
		EvidenceImageUrl: fmt.Sprintf("%s/api/image-processing/resource/evidence?accountId=%d&category=3&businessKey=%d", c.serverDomain, latest.UserId, latest.Id),
	})
}

func (c Controller) AdminCreateDonation(ctx *gin.Context) {
	var dto Dto
	err := ctx.BindJSON(&dto)
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err = interfaces.S.Donation.CreateDonation(EntityFromDto(&dto))
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}

	response.ReplyOK(ctx, nil)
}

func (c Controller) AdminGetDonation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ReplyIllegalArgumentError(ctx, "invalid param id")
		return
	}

	donation, _, err := interfaces.S.Donation.GetDonationById(id)
	response.ReplyOK(ctx, donation)
}

func (c Controller) ApplyDonation(ctx *gin.Context) {
	var apply entities.DonationApply
	if err := ctx.BindJSON(&apply); err != nil {
		response.ReplyError(ctx, err)
		return
	}

	err := interfaces.S.Donation.Apply(apply, ctx.MustGet(pkg.AccountIDKey).(int64))
	if err != nil {
		response.ReplyError(ctx, err)
		return
	}
	response.ReplyOK(ctx, nil)
}

func (c Controller) ApplyExplain(ctx *gin.Context) {
	query := ctx.Query("id")
	//todo
	response.ReplyOK(ctx, query)
}

func SplitToIntArray(s string, split string) ([]int, error) {
	var arr []int
	stringsArray := strings.Split(s, split)
	for _, item := range stringsArray {
		i, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		arr = append(arr, i)
	}

	return arr, nil
}
