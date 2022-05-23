package skill

type Service interface {
	CreateEvidences() error

	// MyCertificate 用户获得认证查询
	// @Param accountId 用户ID int
	MyCertificate(accountId int64) []DtoCert

	// GenCert 生成证书
	// @Param accountId 用户ID int
	// @Param projectID 认证项目ID string
	GenCert(accountId int64, projectID string) (string, string, error)

	// ListCerts 所有证书（灰色证书）
	ListCerts() []DtoCert

	// GetCertByUid 根据证书 ID 读取证书
	// @param uid 证书 ID
	// @return DTO
	GetCertByUid(uid string) (*DtoCert, bool, error)

	// GetUserCertForProject 查询用户在某个Project 的证书
	// @param accountId 用户 ID
	// @param projectId 修炼项目 ID
	// @return DTO
	// @return 是否存在
	// @return error
	GetUserCertForProject(accountId, projectId int64) (*UserCertEntity, bool, error)
}
