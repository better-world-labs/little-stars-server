package middleware

import (
	"aed-api-server/internal/pkg"
	"aed-api-server/internal/service/user"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AccessLog(c *gin.Context) {
	remoteIP := c.GetHeader("X-Forwarded-For")
	if remoteIP == "" {
		ip, _ := c.RemoteIP()
		remoteIP = ip.String()
	}

	data, err := cloneRequestBody(c)
	if err != nil {
		log.Error("AccessLog - cloneRequestBody error:", err)
	}

	log.Infof("api-request|userId-%s %s %s %s %s %s %s\n",
		getUserId(c),
		remoteIP,
		c.Request.Method,
		c.Request.RequestURI,
		c.Request.UserAgent(),
		c.GetHeader("Referer"),
		data,
	)

	blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()

	contentType := c.Writer.Header().Get("Content-Type")
	if strings.Contains(contentType, "json") {
		log.Infof("api-response|%v %s\n", c.Writer.Status(), blw.body.String())
	} else {
		log.Infof("api-response|%v %s\n", c.Writer.Status(), contentType)
	}
}

func getUserId(c *gin.Context) string {
	authorization := c.GetHeader(pkg.AuthorizationHeaderKey)
	split := strings.Split(authorization, " ")

	if len(split) != 2 || split[0] != "Bearer" || split[1] == "" {
		return ""
	}

	token := split[1]
	claims, err := user.ParseToken(token)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%v", claims.ID)
}

func cloneRequestBody(c *gin.Context) ([]byte, error) {
	data, err := c.GetRawData()
	if err != nil {
		return nil, err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return data, nil
}
