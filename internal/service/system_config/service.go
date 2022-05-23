package system_config

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/pkg/db"
	"strings"
)

func Init() {
	interfaces.S.Config = &service{}
}

type service struct{}

type ConfigDO struct {
	Key   string
	Value string
}

func (service) GetConfig(key string) (string, error) {
	var conf ConfigDO
	get, err := db.Table("system_config").Where("`key` =?", key).Get(&conf)
	if err != nil {
		return "", err
	}
	if get {
		return conf.Value, nil
	}
	return "null", nil
}

func (service) PutConfig(key string, config string) error {
	_, err := db.Exec("insert into system_config(`key`, `value`) "+
		`values(?, ?)
		ON DUPLICATE KEY UPDATE
		value = ?
	`, key, config, config)
	return err
}

func (service) GetAllConfig() (string, error) {
	dos := make([]ConfigDO, 0)

	err := db.Table("system_config").Find(&dos)
	if err != nil {
		return "", err
	}

	strList := make([]string, 0)
	for i := 0; i < len(dos); i++ {
		do := dos[i]
		strList = append(strList, `"`+do.Key+`":`+do.Value)
	}
	return "{" + strings.Join(strList, ",") + "}", nil
}
