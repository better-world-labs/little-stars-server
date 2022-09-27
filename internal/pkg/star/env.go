package star

const (
	Local Env = "local"
	Dev   Env = "dev"
	Test  Env = "test"
	Prod  Env = "prod"
)

type Env string

func GetDomainPrefix(env Env) string {
	switch env {
	case Local, Dev:
		return "dev."
	case Test:
		return "test."
	default:
		return ""
	}
}

func GetEnvText(env Env) string {
	switch env {
	case Local, Dev:
		return "开发"
	case Test:
		return "测试"
	default:
		return ""
	}
}
