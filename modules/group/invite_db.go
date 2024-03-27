package group

import (
	"gorm.io/gorm"
	"time"
)

// InsertInviteTx 添加邀请信息
func (d *DB) InsertInviteTx(model *InviteModel, tx *gorm.DB) error {
	err := tx.Table("group_invite").Create(model).Error
	return err
}

// InsertInviteItemTx 添加邀请项
func (d *DB) InsertInviteItemTx(model *InviteItemModel, tx *gorm.DB) error {
	err := tx.Table("invite_item").Create(model).Error
	return err
}

// QueryInviteDetail 查询邀请详情
func (d *DB) QueryInviteDetail(inviteNo string) (*InviteDetailModel, error) {
	var m InviteDetailModel
	err := d.db.Select("group_invite.*,IFNULL(user.name,'') inviter_name").Table("group_invite").Joins("user on group_invite.inviter=user.uid").Where("invite_no=?", inviteNo).First(&m).Error
	return &m, err
}

// QueryInviteItemDetail 查询邀请item详情
func (d *DB) QueryInviteItemDetail(inviteNo string) ([]*InviteItemDetailModel, error) {
	var ms []*InviteItemDetailModel
	err := d.db.Select("invite_item.*,IFNULL(user.name,'') name").Table("invite_item").Joins("user on invite_item.uid=user.uid").Where("invite_item.invite_no=?", inviteNo).Find(&ms).Error
	return ms, err

}

// UpdateInviteStatusTx 更新邀请信息状态
func (d *DB) UpdateInviteStatusTx(allower string, status int, inviteNo string, tx *gorm.DB) error {
	err := tx.Table("group_invite").Update("allower", allower).Update("status", status).Where("invite_no=?", inviteNo).Error
	return err
}

// UpdateInviteItemStatusTx 更新邀请信息状态
func (d *DB) UpdateInviteItemStatusTx(status int, inviteNo string, tx *gorm.DB) error {
	err := tx.Table("invite_item").Update("status", status).Where("invite_no=?", inviteNo).Error
	return err
}

// InviteModel InviteModel
type InviteModel struct {
	InviteNo  *string `json:"invite_no"` // 邀请唯一编号
	GroupNo   *string `json:"group_no"`  // 群唯一编号
	Inviter   *string `json:"inviter"`   // 邀请者
	Remark    *string `json:"remark"`    // 邀请备注
	Status    *int    `json:"status"`    // 状态 0.未确认 1.已确认
	Allower   *string `json:"allower"`   // 确认者
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// InviteDetailModel 邀请者详情
type InviteDetailModel struct {
	InviteModel
	InviterName string `json:"inviter_name"` // 邀请者名称

}

// InviteItemDetailModel item详情
type InviteItemDetailModel struct {
	InviteItemModel
	Name string // 被邀请者名称
}

// InviteItemModel InviteItemModel
type InviteItemModel struct {
	InviteNo  *string `json:"invite_no"` // 邀请唯一编号
	GroupNo   *string `json:"group_no"`  // 群唯一编号
	Inviter   *string `json:"inviter"`   // 邀请者
	UID       *string `json:"uid"`       // 被邀请uid
	Status    *int    `json:"status"`    // 状态 0.未确认 1.已确认
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
