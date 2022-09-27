package user

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"aed-api-server/internal/pkg/db"
	"aed-api-server/internal/pkg/domain/emitter"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

type traceService struct {
	User   service.UserService `inject:"-"`
	Wechat service.IWechat     `inject:"-"`
	Oss    service.OssService  `inject:"-"`
}

//go:inject-component
func NewTraceService() service.TraceService {
	return &traceService{}
}

func (s traceService) Create(code string, trace entities.Trace) (*entities.Trace, error) {
	session, err := s.User.Code2Session(code)
	if err != nil {
		return nil, err
	}

	trace.OpenID = session.Openid
	if trace.OpenID == "" {
		return nil, err
	}

	return &trace, db.Begin(func(session *xorm.Session) error {
		_, err = db.GetEngine().Table("generalize_trace").Insert(trace)
		if err != nil {
			return err
		}
		return emitter.Emit(&trace)
	})
}

func (s traceService) GetEarliestSharerTrace(openid string) (*entities.Trace, bool, error) {
	var t entities.Trace
	exists, err := db.SQL(`select * from generalize_trace where open_id = ?
		and sharer regexp '^[0-9]+$' limit 1`, openid).Get(&t)
	if err != nil {
		return nil, exists, err
	}

	return &t, exists, nil
}

func (s traceService) CreateQrCode(req *entities.CreateQrCodeReq) (*entities.CreateQrCodeRes, error) {
	if req.Sharer == "" || req.Source == "" || req.PagePath == "" {
		return nil, errors.New("参数不全")
	}
	path, err := url.Parse(req.PagePath)
	if err != nil {
		return nil, err
	}

	query := path.Query()

	query.Set("source", req.Source)
	query.Set("sharer", req.Sharer)
	path.RawQuery = query.Encode()

	log.Info("path:", path.String())

	r, t, err := s.Wechat.GenMinaCode(path.String(), 1280, false, false, "")
	if strings.HasPrefix(t, "image/") {
		t = strings.Replace(t, "image/", "", -1)
	} else {
		log.Errorf("gen image content-type is err")
		return nil, errors.New("gen image content-type is err")
	}

	if err != nil {
		return nil, err
	}
	id := utils.GetUUID()

	imgUrl, err := s.Oss.Upload(fmt.Sprintf("share-trace/%s.%s", id, t), r)

	if err != nil {
		return nil, err
	}
	return &entities.CreateQrCodeRes{
		Image: imgUrl,
	}, nil
}
