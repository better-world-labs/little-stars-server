package exam

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImortQuestions(t *testing.T) {
	//db.InitEngine(db.MysqlConfig{
	//	DriverName: "mysql",
	//	Dsn:        "db_account_star_dev:db_account_star_dev123@tcp(rm-bp11mfhb2120j3s801o.mysql.rds.aliyuncs.com:3306)/star_dev?charset=utf8mb4&loc=Local",
	//})

	var v interface{}
	err := json.Unmarshal([]byte("2"), &v)
	require.Nil(t, err)
	//service := NewQuestionService()

	// 题库文件丢了，下次导入的时候再说
	//file, err := os.Open("/home/shenweijie/下载/题库.xlsx")
	//require.Nil(t, err)

	//err = service.ImportQuestions(file, 1, 1)
	//require.Nil(t, err)
}
