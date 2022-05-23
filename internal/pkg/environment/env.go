package environment

const (
	Testing    = "testing"
	Preview    = "preview"
	Production = "production"
)

func GetDomainPrefix(env string) string {
	switch env {
	case Testing:
		return "dev."
	case Preview:
		return "test."
	default:
		return ""
	}
}
