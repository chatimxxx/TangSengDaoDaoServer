package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/pkg/db"
	"gorm.io/gorm"
)

type giteeDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newGiteeDB(ctx *config.Context) *giteeDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &giteeDB{
		ctx: ctx,
		db:  d,
	}
}

func (d *giteeDB) insert(m *gitUserInfoModel) error {
	err := d.db.Table("gitee_user").Create(m).Error
	return err
}

func (d *giteeDB) insertTx(m *gitUserInfoModel, tx *gorm.DB) error {
	err := tx.Table("gitee_user").Create(m).Error
	return err
}
func (d *giteeDB) queryWithLogin(login string) (*gitUserInfoModel, error) {
	var m gitUserInfoModel
	err := d.db.Table("gitee_user").Where("login=?", login).First(&m).Error
	return &m, err
}

type gitUserInfoModel struct {
	Id                int64
	CreatedAt         db.Time
	UpdatedAt         db.Time
	Login             string
	Name              string
	Email             string
	Bio               string
	AvatarURL         string
	Blog              string
	EventsURL         string
	Followers         int
	FollowersURL      string
	Following         int
	FollowingURL      string
	GistsURL          string
	HtmlURL           string
	MemberRole        string
	OrganizationsURL  string
	PublicGists       int
	PublicRepos       int
	ReceivedEventsURL string
	Remark            string
	ReposURL          string
	Stared            int
	StarredURL        string
	SubscriptionsURL  string
	URL               string
	Watched           int
	Weibo             string
	Type              string
	GiteeCreatedAt    string
	GiteeUpdatedAt    string
}
