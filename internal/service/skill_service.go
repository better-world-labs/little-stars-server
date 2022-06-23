package service

import (
	"aed-api-server/internal/interfaces"
	"aed-api-server/internal/interfaces/entities"
	service2 "aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/module/cert"
	"aed-api-server/internal/module/skill"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/global"
	"aed-api-server/internal/pkg/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ali "github.com/aliyun/aliyun-oss-go-sdk/oss"
	log "github.com/sirupsen/logrus"
	"time"

	"aed-api-server/internal/module/oss"
)

type skillService struct {
	oss            *oss.Config
	AccountService service2.UserServiceOld `inject:"-"`
}

func NewSkillService(c *oss.Config) service2.SkillService {
	return &skillService{oss: c}
}

func (s skillService) MyCertificate(accountId int64) []entities.DtoCert {
	sess := db.GetSession()
	defer sess.Close()

	certs := make([]entities.UserCert, 0)
	sess.Where("account_id = ?", accountId).Find(&certs)
	res := make([]entities.DtoCert, 0)
	for _, v := range certs {
		imgs := make(map[string]string)
		json.Unmarshal([]byte(v.Img), &imgs)
		res = append(res, entities.DtoCert{
			Uid:         v.Uid,
			ProjectId:   v.ProjectId,
			ProjectName: v.ProjectName,
			Origin:      imgs["origin"],
			Thumbnail:   imgs["thumbnail"],
			Created:     v.Created,
		})
	}

	return res
}

func (s skillService) CreateEvidences() error {
	certs, err := s.ListAllCerts()
	if err != nil {
		return err
	}

	for _, cert := range certs {
		_, exists, err := interfaces.S.Evidence.GetEvidenceByBusinessKey(cert.Uid, entities.EvidenceCategoryCert)
		if err != nil {
			return err
		}

		if !exists {
			log.Infof("evidence for cert %s is not fond, generate one", cert.Uid)
			errChan := interfaces.S.Evidence.CreateCertEvidenceAsync(cert.AccountId, cert.Uid, "\"茫茫人海之中，去挽救下一个倒地昏迷的人吧\"")
			err = <-errChan
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func (s skillService) GenCert(accountID int64, projectID string) (string, string, error) {
	sess := db.GetSession()
	defer sess.Close()

	desc := "\"茫茫人海之中，去挽救下一个倒地昏迷的人吧\""

	certID := fmt.Sprintf("%d%v", accountID, time.Now().Nanosecond())
	evidenceErrChan := interfaces.S.Evidence.CreateCertEvidenceAsync(accountID, certID, desc)
	creator, err := cert.NewImageCreatorDefaultAssert()
	if err != nil {
		return "", "", err
	}

	acc := new(entities.User)
	acc.ID = accountID
	exists, err := sess.Table("account").Get(acc) //TODO 不要直接查库,清理掉!!!
	if err != nil {
		return "", "", err
	}

	if !exists {
		return "", "", errors.New("not found")
	}

	var writer bytes.Buffer
	err = creator.Create(acc.Avatar, acc.Nickname, desc, time.Now(), &writer)
	if err != nil {
		return "", "", err
	}

	client, err := ali.New(s.oss.Endpoint, s.oss.AccesskeyId, s.oss.AccesskeySecret)
	if err != nil {
		return "", "", err
	}

	bucket, err := client.Bucket(s.oss.BucketName)
	if err != nil {
		return "", "", err
	}

	certImg := fmt.Sprintf("%v/%v/cert_%v", s.oss.UploadDir, accountID, projectID)
	err = bucket.PutObject(certImg, bytes.NewReader(writer.Bytes()))
	if err != nil {
		return "", "", err
	}

	cert := new(entities.UserCert)
	cert.AccountId = accountID
	cert.ProjectId = utils.ToInt(projectID)
	imgs := make(map[string]string)
	imgs["origin"] = fmt.Sprintf("https://%s/%s", s.oss.Domain, certImg)
	imgs["thumbnail"] = fmt.Sprintf("https://%s/%s", s.oss.Domain, certImg)
	b, _ := json.Marshal(imgs)
	cert.Img = string(b)
	cert.Uid = certID
	cert.Created = global.FormattedTime(time.Now())

	err = db.WithTransaction(sess, func() error {
		now := time.Now()
		_, err := sess.Exec(`insert user_cert (uid, account_id, project_id, img, created) 
        values (?, ?, ?, ?, ?) on duplicate key update uid = ?, img = ?, created = ?`, cert.Uid,
			cert.AccountId,
			cert.ProjectId,
			cert.Img,
			now,
			cert.Uid,
			cert.Img,
			now)
		if err != nil {
			return err
		}

		err = <-evidenceErrChan
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return imgs["origin"], certID, nil
}

func (s skillService) ListCerts() []entities.DtoCert {
	sess := db.GetSession()
	defer sess.Close()

	prjs := make([]skill.Project, 0)
	sess.Find(&prjs)
	res := make([]entities.DtoCert, 0)
	for _, v := range prjs {
		imgs := make(map[string]string)
		json.Unmarshal([]byte(v.GrayImg), &imgs)
		res = append(res, entities.DtoCert{
			ProjectId:   v.Id,
			ProjectName: v.Name,
			Origin:      imgs["origin"],
			Thumbnail:   imgs["thumbnail"],
		})
	}

	return res
}

func (s skillService) GetCertByUid(uid string) (*entities.DtoCert, bool, error) {
	session := db.GetSession()
	defer session.Close()

	var res entities.DtoCert
	exists, err := session.Table("user_cert").Where("uid = ?", uid).Get(&res)
	if err != nil {
		return nil, false, err
	}

	return &res, exists, nil
}

func (s *skillService) GetUserCertForProject(accountId, projectId int64) (*entities.UserCertEntity, bool, error) {
	session := db.GetSession()
	defer session.Close()

	var res entities.UserCertEntity
	exists, err := session.Table("user_cert").Where("project_id = ? and account_id = ?", projectId, accountId).Get(&res)
	if err != nil {
		return nil, false, err
	}

	return &res, exists, nil
}

func (s *skillService) ListAllCerts() ([]*entities.UserCertEntity, error) {
	session := db.GetSession()
	defer session.Close()

	var res []*entities.UserCertEntity
	err := session.Table("user_cert").Find(&res)

	return res, err
}
