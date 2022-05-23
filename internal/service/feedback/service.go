package feedback

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"io"
	"time"
)

func Init() {
	interfaces.S.Feedback = &service{}
}

type service struct{}

type FeedbackDO struct {
	UserId    int64
	CreatedAt time.Time
	Type      int
	Content   string
	Images    []string
}

func (service) SubmitFeedback(userId int64, feedback *entities.Feedback) error {
	utils.Go(func() {
		sendDingTalk(feedback.Type, feedback.Content, feedback.Images, userId)
	})

	_, err := db.Table("feedback").Insert(FeedbackDO{
		UserId:    userId,
		CreatedAt: time.Now(),
		Type:      feedback.Type,
		Content:   feedback.Content,
		Images:    feedback.Images,
	})

	if err != nil {
		return err
	}
	return nil
}

func sendDingTalk(t int, content string, images []string, userId int64) {
	config := interfaces.GetConfig()
	typeStr := "功能异常"
	if t == 2 {
		typeStr = "产品建议"
	}

	var imgStr string
	for i := range images {
		imgStr += fmt.Sprintf("![](%s)\n", images[i])
	}

	user, _, _ := interfaces.S.User.GetUserById(userId)

	var msgStr = fmt.Sprintf(`
### 用户反馈-%s 
%s  
%s
用户：%s(%s)
`, typeStr, content, imgStr, user.Nickname, user.Uid)

	err := utils.SendDingTalkBot(config.DonationApplyNotify, &utils.DingTalkMsg{
		Msgtype: "markdown",
		Markdown: utils.Markdown{
			Title: "用户反馈",
			Text:  msgStr,
		},
	})
	if err != nil {
		log.Error("sendDingTalk error", err)
	}
}

func (service) ExportFeedback(beginDate time.Time, endDate time.Time, writer *io.PipeWriter) {
	go doExport(beginDate, endDate, writer)
}

var TableTitle = []string{
	"用户ID",
	"用户昵称",
	"反馈类型",
	"反馈内容",
	"图片1",
	"图片2",
	"图片3",
}

type TableRow struct {
	UserId          int64
	UserName        string
	FeedBackType    string
	FeedbackContent string
	ImageUrl1       string
	ImageUrl2       string
	ImageUrl3       string
}

func doExport(beginDate time.Time, endDate time.Time, writer *io.PipeWriter) {
	defer func() {
		info := recover()
		if info != nil {
			log.Error("recover:", info)
		}
		_ = writer.Close()
	}()

	file := excelize.NewFile()
	excelWriter, _ := file.NewStreamWriter("Sheet1")
	writeTableTitle(excelWriter, TableTitle)
	rows, err := fetchFeedbackData(beginDate, endDate)
	if err != nil {
		writeTableError(excelWriter, err)
		return
	}

	writeTableData(excelWriter, rows)

	if err := excelWriter.Flush(); err != nil {
		log.Error("writeTableTitle flush err", err)
	}
	if err := file.Write(writer); err != nil {
		log.Error("file.Write err:", err)
	}
}

func writeTableData(excelWriter *excelize.StreamWriter, rows []*TableRow) {
	rowsCount := len(rows)
	for i := 0; i < rowsCount; i++ {
		rowID := i + 2
		item := rows[i]
		row := make([]interface{}, 0)
		row = append(row,
			excelize.Cell{Value: item.UserId},
			excelize.Cell{Value: item.UserName},
			excelize.Cell{Value: item.FeedBackType},
			excelize.Cell{Value: item.FeedbackContent},
			excelize.Cell{Value: item.ImageUrl1},
			excelize.Cell{Value: item.ImageUrl2},
			excelize.Cell{Value: item.ImageUrl3},
		)
		cell, err := excelize.CoordinatesToCellName(1, rowID)
		if err != nil {
			log.Error("excelize.CoordinatesToCellName", err)
		}

		if err := excelWriter.SetRow(cell, row); err != nil {
			log.Error("writeTableData err:", err)
		}
	}
}

func writeTableError(excelWriter *excelize.StreamWriter, e error) {
	if err := excelWriter.SetRow("A2", []interface{}{e.Error()}); err != nil {
		log.Error("writeTableError err:", err)
	}
	if err := excelWriter.Flush(); err != nil {
		log.Error("writeTableError flush err", err)
	}
}

func writeTableTitle(excelWriter *excelize.StreamWriter, titles []string) {
	data := make([]interface{}, 0, len(titles))
	for i := 0; i < len(titles); i++ {
		data = append(data, excelize.Cell{Value: titles[i]})
	}
	if err := excelWriter.SetRow("A1", data); err != nil {
		log.Error("writeTableTitle err", err)
	}
}

func fetchFeedbackData(beginDate time.Time, endDate time.Time) ([]*TableRow, error) {
	list := make([]*FeedbackDO, 0)

	err := db.Table("feedback").Where("created_at between ? and ?", beginDate, endDate).Find(&list)
	if err != nil {
		log.Error("err info", err)
		return nil, err
	}

	userIdMap := make(map[int64]bool)
	userIds := make([]int64, 0)
	for i := range list {
		id := list[i].UserId

		_, ok := userIdMap[id]
		if !ok {
			userIdMap[id] = true
			userIds = append(userIds, id)
		}
	}

	ds, err := interfaces.S.User.GetListUserByIDs(userIds)
	if err != nil {
		log.Error("err info", err)
		return nil, err
	}

	userMap := make(map[int64]*entities.SimpleUser)
	for i := range ds {
		user := ds[i]
		userMap[user.ID] = user
	}

	rows := make([]*TableRow, 0)
	for i := range list {
		feedback := list[i]
		user := userMap[feedback.UserId]
		row := TableRow{
			UserId:          user.ID,
			UserName:        user.Nickname,
			FeedBackType:    getFeedbackTypeName(feedback.Type),
			FeedbackContent: feedback.Content,
		}

		images := feedback.Images
		if len(images) > 0 {
			row.ImageUrl1 = images[0]
		}
		if len(images) > 1 {
			row.ImageUrl2 = images[1]
		}
		if len(images) > 2 {
			row.ImageUrl3 = images[2]
		}
		rows = append(rows, &row)
	}
	return rows, nil
}

func getFeedbackTypeName(t int) string {
	if t == 1 {
		return "功能异常"
	}
	if t == 2 {
		return "产品建议"
	}
	return ""
}
