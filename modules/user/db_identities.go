package user

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"gorm.io/gorm"
	"time"
)

type identitieDB struct {
	db  *gorm.DB
	ctx *config.Context
}

func newIdentitieDB(ctx *config.Context) *identitieDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &identitieDB{
		db:  d,
		ctx: ctx,
	}
}

func (i *identitieDB) saveOrUpdateTx(m *identitiesModel, tx *gorm.DB) error {
	err := tx.Exec("insert into signal_identities(uid,identity_key,signed_prekey_id,signed_pubkey,signed_signature,registration_id) values(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE identity_key=identity_key,signed_prekey_id=signed_prekey_id,signed_pubkey=signed_pubkey,signed_signature=signed_signature,registration_id=registration_id", m.UID, m.IdentityKey, m.SignedPrekeyID, m.SignedPubkey, m.SignedSignature, m.RegistrationID).Error
	return err
}

func (i *identitieDB) deleteWithUID(uid string) error {
	err := i.db.Table("signal_identities").Where("uid=?", uid).Delete(nil).Error
	return err
}

func (i *identitieDB) queryWithUID(uid string) (*identitiesModel, error) {
	var m identitiesModel
	err := i.db.Table("signal_identities").Where("uid=?", uid).First(&m).Error
	return &m, err
}

type identitiesModel struct {
	UID             *string
	RegistrationID  *uint32
	IdentityKey     *string
	SignedPrekeyID  *int
	SignedPubkey    *string
	SignedSignature *string
	Id              *int64
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
