package audit

import (
	"aed-api-server/internal/interfaces/entities"
	"aed-api-server/internal/interfaces/service"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/buffer"
	"strings"
	"time"
)

type Response struct {
	Code int                    `json:"code"`
	Data []entities.AuditResult `json:"data"`
	Msg  string                 `json:"msg"`
}

type Audit struct {
	AccessKey    string `conf:"aliyun-scan.access-key-id"`
	AccessSecret string `conf:"aliyun-scan.access-key-secret"`
	Endpoint     string `conf:"aliyun-scan.endpoint"`
}

//go:inject-component
func NewAudit() service.IContentAudit {
	return &Audit{}
}

func (a *Audit) ScanImage(imgUrl string) (*entities.AuditResult, error) {
	body := createScanImageBody(imgUrl)

	data, err := a.request("/green/image/scan", body)
	if err != nil {
		return nil, err
	}

	return &data[0], nil
}

func (a *Audit) ScanText(text string) (*entities.AuditResult, error) {
	body := createScanTextBody(text)

	data, err := a.request("/green/text/scan", body)
	if err != nil {
		return nil, err
	}

	return &data[0], nil
}

func md5Base64(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	sum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func createScanTextBody(text string) string {
	return fmt.Sprintf(`{"scenes":["antispam"],"tasks":[{"dataId":"%s","content":"%s"}]}`, uuid.New(), text)
}

func createScanImageBody(url string) string {
	return fmt.Sprintf(`{"scenes":["porn","terrorism","ad","live","qrcode","logo"],"tasks":[{"dataId":"%s","url":"%s"}]}`, uuid.New(), url)
}

func sign(header req.Header, accessKey, accessSecret, uri string) {
	var keys []string

	for k, _ := range header {
		if strings.HasPrefix(k, "x-acs") {
			keys = append(keys, k)
		}
	}

	toSign := buffer.Buffer{}
	toSign.AppendString("POST\n")
	toSign.AppendString(header["Content-Type"])
	toSign.AppendString("\n")
	toSign.AppendString(header["Content-MD5"])
	toSign.AppendString("\n")
	toSign.AppendString(header["Accept"])
	toSign.AppendString("\n")
	toSign.AppendString(header["Date"])
	toSign.AppendString("\n")
	toSign.AppendString("x-acs-signature-method:")
	toSign.AppendString(header["x-acs-signature-method"])
	toSign.AppendString("\n")
	toSign.AppendString("x-acs-signature-nonce:")
	toSign.AppendString(header["x-acs-signature-nonce"])
	toSign.AppendString("\n")
	toSign.AppendString("x-acs-signature-version:")
	toSign.AppendString(header["x-acs-signature-version"])
	toSign.AppendString("\n")
	toSign.AppendString("x-acs-version:")
	toSign.AppendString(header["x-acs-version"])
	toSign.AppendString("\n")
	toSign.AppendString(uri)
	fmt.Println(toSign.String())

	h := hmac.New(sha1.New, []byte(accessSecret))
	h.Write([]byte(toSign.String()))
	sum := h.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(sum)
	header["Authorization"] = fmt.Sprintf("acs %s:%s", accessKey, signature)
}

func FormatGMT(t time.Time) string {
	week := t.Weekday().String()[:3]
	day := t.Day()
	month := t.Month().String()[:3]
	year := t.Year()
	time := t.Format("15:04:05")
	return fmt.Sprintf(fmt.Sprintf("%s, %d %s %d %s GMT", week, day, month, year, time))
}

func (a *Audit) createHeader(uri, body string) req.Header {
	header := req.Header{
		"Content-Type":            "application/json",
		"Content-MD5":             md5Base64(body),
		"Accept":                  "application/json",
		"Date":                    FormatGMT(time.Now().UTC()),
		"x-acs-signature-method":  "HMAC-SHA1",
		"x-acs-signature-nonce":   uuid.NewString(),
		"x-acs-signature-version": "1",
		"x-acs-version":           "2018-05-09",
	}
	sign(header, a.AccessKey, a.AccessSecret, uri)
	return header
}

func (a *Audit) request(uri, body string) ([]entities.AuditResult, error) {
	resp, err := req.Post(fmt.Sprintf("https://%s/%s", a.Endpoint, uri), a.createHeader(uri, body), body)
	if err != nil {
		return nil, err
	}

	var res Response
	logrus.Infof("ScanRequest: uri=%s resp=%s", uri, resp.String())
	err = json.Unmarshal(resp.Bytes(), &res)
	if err != nil {
		return nil, err
	}

	if res.Code != 200 {
		return nil, errors.New(res.Msg)
	}

	return res.Data, nil
}
