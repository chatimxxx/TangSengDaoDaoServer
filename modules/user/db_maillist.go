package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type maillistDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newMaillistDB(ctx *config.Context) *maillistDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &maillistDB{
		ctx: ctx,
		db:  d,
	}
}

func (d *maillistDB) insertTx(m *maillistModel, tx *gorm.DB) error {
	err := tx.Exec("INSERT INTO user_maillist (zone,phone,name,vercode,uid) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `phone`=VALUES(`phone`)", m.Zone, m.Phone, m.Name, m.Vercode, m.UID).Error
	return err
}

func (d *maillistDB) queryWitchVercode(vercode string) (*maillistModel, error) {
	var m maillistModel
	err := d.db.Table("user_maillist").Where("vercode=? ", vercode).First(&m).Error
	return &m, err
}
func (d *maillistDB) query(uid string) ([]*maillistModel, error) {
	var ms []*maillistModel
	err := d.db.Table("user_maillist").Where("uid=?", uid).Find(&ms).Error
	return ms, err
}

type maillistModel struct {
	UID       *string
	Phone     *string
	Zone      *string
	Name      *string
	Vercode   *string
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
