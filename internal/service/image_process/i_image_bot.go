package image_process

type IImageBot interface {
	Call(tplName string, args map[string]interface{}) (string, error)
}
