package global

import (
	"aed-api-server/internal/pkg/response"
	"net/http"
)

//TODO HTTP 包装其他Error,
var (
	ErrorUnknown              = response.NewHTTPServiceError(http.StatusInternalServerError, 3000, "系统错误")
	ErrorMessagePublishFailed = response.NewHTTPServiceError(http.StatusInternalServerError, 3001, "消息发布异常")
	ErrorInvalidAccessToken   = response.NewHTTPServiceError(http.StatusUnauthorized, 1000, "用户令牌无效")
	ErrorAccountNotFound      = response.NewHTTPServiceError(http.StatusUnauthorized, 1002, "用户不存在")
	ErrorInvalidParam         = response.NewHTTPServiceError(http.StatusBadRequest, 1001, "操作失败，请重试")
	ErrorInvalidDistance      = response.NewHTTPServiceError(http.StatusBadRequest, 2000, "操作失败，您与目的地距离超过1km")
	//ErrorWithoutArrived       = response.NewHTTPServiceError(http.StatusBadRequest, 2001, "请先确认到达现场")
	ErrorLinkInvalid = response.NewHTTPServiceError(http.StatusOK, 4000, "链接无效或已经过期")

	// exam
	ErrorExamUnSubmit     = response.NewHTTPServiceError(http.StatusOK, 5000, "请继续作答")
	ErrorExamOwnerError   = response.NewHTTPServiceError(http.StatusOK, 5001, "参考人员不一致")
	ErrorAlreadyCompleted = response.NewHTTPServiceError(http.StatusOK, 5002, "考试已经提交,请勿重复提交")
	ErrorUnknownQuestion  = response.NewHTTPServiceError(http.StatusOK, 5003, "答案与试题不匹配")
)
