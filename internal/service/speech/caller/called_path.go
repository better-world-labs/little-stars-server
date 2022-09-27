package caller

import (
	"aed-api-server/internal/interfaces/service"
	"fmt"
)

//TODO 重新抽象
type PathGenerator interface {
	GeneratePath(aid int64, value string) (string, error)
}

type pathGenerator struct {
	tokenService service.TokenService
}

func NewPathGenerator(tokenService service.TokenService) PathGenerator {
	return &pathGenerator{
		tokenService: tokenService,
	}
}

func (p pathGenerator) GeneratePath(aid int64, value string) (string, error) {
	token := p.tokenService.Generate()
	err := p.tokenService.PutToken(token, value)
	return fmt.Sprintf("api/p/c/%s/%d", token, aid), err
}
