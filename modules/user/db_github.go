package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type githubDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newGithubDB(ctx *config.Context) *githubDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &githubDB{
		ctx: ctx,
		db:  d,
	}
}

func (d *githubDB) insert(m *githubUserInfoModel) error {
	err := d.db.Table("github_user").Create(m).Error
	return err
}

func (d *githubDB) insertTx(m *githubUserInfoModel, tx *gorm.DB) error {
	err := tx.Table("github_user").Create(m).Error
	return err
}
func (d *githubDB) queryWithLogin(login string) (*githubUserInfoModel, error) {
	var m githubUserInfoModel
	err := d.db.Table("github_user").Where("login=?", login).First(&m).Error
	return &m, err
}

type githubUserInfoModel struct {
	ID                      int64
	CreatedAt               db.Time
	UpdatedAt               db.Time
	Login                   string
	NodeID                  string
	AvatarURL               string
	GravatarID              string
	URL                     string
	HtmlUrl                 string
	FollowersURL            string
	FollowingURL            string
	GistsURL                string
	StarredURL              string
	SubscriptionsURL        string
	OrganizationsURL        string
	ReposURL                string
	EventsURL               string
	ReceivedEventsURL       string
	Type                    string
	SiteAdmin               bool
	Name                    string
	Company                 string
	Blog                    string
	Location                string
	Email                   string
	Hireable                bool
	Bio                     string
	TwitterUsername         string
	PublicRepos             int
	PublicGists             int
	Followers               int
	Following               int
	GithubCreatedAt         string
	GithubUpdatedAt         string
	PrivateGists            int
	TotalPrivateRepos       int
	OwnedPrivateRepos       int
	DiskUsage               int
	Collaborators           int
	TwoFactorAuthentication bool
}
