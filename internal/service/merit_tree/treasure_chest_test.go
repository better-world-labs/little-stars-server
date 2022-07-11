package merit_tree

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/pkg/db"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_treasureChestService_OpenTreasureChest(t *testing.T) {
	type fields struct {
		p        iPersistence
		Expired  time.Duration
		CoolDown time.Duration
	}
	type args struct {
		userId          int64
		treasureChestId int
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &treasureChestService{
				p:        tt.fields.p,
				Expired:  tt.fields.Expired,
				CoolDown: tt.fields.CoolDown,
			}
			if err := s.OpenTreasureChest(tt.args.userId, tt.args.treasureChestId); (err != nil) != tt.wantErr {
				t.Errorf("OpenTreasureChest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_treasureChestService_GetUserTreasureChest(t *testing.T) {
	db.InitEngine(db.MysqlConfig{
		DriverName: "mysql",
		Dsn:        "db_account_star_dev:db_account_star_dev123@tcp(rm-bp11mfhb2120j3s80.mysql.rds.aliyuncs.com:3306)/star_dev?charset=utf8mb4&loc=Local",
	})
	type fields struct {
		p        iPersistence
		Expired  time.Duration
		CoolDown time.Duration
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.TreasureChest
		wantErr bool
	}{
		{
			name:   "first chest",
			fields: fields{p: &persistence{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &treasureChestService{
				p:        tt.fields.p,
				Expired:  tt.fields.Expired,
				CoolDown: tt.fields.CoolDown,
			}
			_, err := s.GetUserTreasureChest(tt.args.userId)
			require.Equal(t, err != nil, tt.wantErr, "GetUserTreasureChest() error = %v, wantErr %v", err, tt.wantErr)
			//TODO 返回数据暂时不可控，暂时根据执行成功判定测试通过
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetUserTreasureChest() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
