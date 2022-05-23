package speech

import (
	"fmt"
)

//TODO 重新抽象
type PathGenerator interface {
	GeneratePath(aid int64, value string) (string, error)
}

type pathGenerator struct {
	tokenService TokenService
}

func NewPathGenerator(tokenService TokenService) PathGenerator {
	return &pathGenerator{
		tokenService: tokenService,
	}
}

func (p pathGenerator) GeneratePath(aid int64, value string) (string, error) {
	token := p.tokenService.Generate()
	err := p.tokenService.PutToken(token, value)
	return fmt.Sprintf("p/c/%s/%d", token, aid), err
}
