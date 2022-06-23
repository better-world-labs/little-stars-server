package merit_tree

import (
	"aed-api-server/internal/interfaces/entities"
	"reflect"
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &treasureChestService{
				p:        tt.fields.p,
				Expired:  tt.fields.Expired,
				CoolDown: tt.fields.CoolDown,
			}
			got, err := s.GetUserTreasureChest(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTreasureChest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserTreasureChest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
