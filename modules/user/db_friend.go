package user

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
)

// DB DB
type friendDB struct {
	db  *gorm.DB
	ctx *config.Context
}

// NewDB NewDB
func newFriendDB(ctx *config.Context) *friendDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &friendDB{
		db:  d,
		ctx: ctx,
	}
}

// InsertTx 插入好友信息
func (d *friendDB) InsertTx(m *FriendModel, tx *gorm.DB) error {
	err := tx.Table("friend").Create(m).Error
	if err != nil {
		return err
	}
	friendKey := fmt.Sprintf("%s%s", CacheKeyFriends, m.UID)
	err = d.ctx.GetRedisConn().SAdd(friendKey, m.ToUID)
	return err
}

// Insert 插入好友信息
func (d *friendDB) Insert(m *FriendModel) error {
	err := d.db.Table("friend").Create(m).Error
	if err != nil {
		return err
	}
	friendKey := fmt.Sprintf("%s%s", CacheKeyFriends, m.UID)
	err = d.ctx.GetRedisConn().SAdd(friendKey, m.ToUID)
	return err
}

// IsFriend 是否是好友
func (d *friendDB) IsFriend(uid, toUID string) (bool, error) {
	var m *FriendModel
	err := d.db.Table("friend").Where("uid=? and to_uid=?", uid, toUID).First(&m).Error
	if err != nil {
		return false, err
	}
	var isFriend = false
	if m != nil && *m.IsDeleted == 0 {
		isFriend = true
	}
	return isFriend, nil
}

// 修改好友关系
func (d *friendDB) updateRelationshipTx(uid, toUID string, isDeleted, isAlone int, sourceVercode string, version int64, tx *gorm.DB) error {
	err := tx.Table("friend").Updates(map[string]interface{}{
		"is_deleted":     isDeleted,
		"is_alone":       isAlone,
		"source_vercode": sourceVercode,
		"version":        version,
	}).Where("uid=? and to_uid=?", uid, toUID).Error
	if err != nil {
		return err
	}
	friendKey := fmt.Sprintf("%s%s", CacheKeyFriends, uid)
	if isDeleted == 1 {
		err = d.ctx.GetRedisConn().SRem(friendKey, toUID)
	} else {
		err = d.ctx.GetRedisConn().SAdd(friendKey, toUID)
	}

	return err
}

func (d *friendDB) updateRelationship2Tx(uid, toUID string, isDeleted, isAlone int, version int64, tx *gorm.DB) error {
	err := tx.Table("friend").Updates(map[string]interface{}{
		"is_deleted": isDeleted,
		"is_alone":   isAlone,
		"version":    version,
	}).Where("uid=? and to_uid=?", uid, toUID).Error
	if err != nil {
		return err
	}
	friendKey := fmt.Sprintf("%s%s", CacheKeyFriends, uid)
	if isDeleted == 1 {
		err = d.ctx.GetRedisConn().SRem(friendKey, toUID)
	} else {
		err = d.ctx.GetRedisConn().SAdd(friendKey, toUID)
	}

	return err
}

// 修改好友单项关系
func (d *friendDB) updateAloneTx(uid, toUID string, isAlone int, tx *gorm.DB) error {
	err := tx.Table("friend").Set("is_alone", isAlone).Where("uid=? and to_uid=?", uid, toUID).Error
	return err
}

// 删除好友
// func (d *friendDB) delete(uid, toUID string) error {
// 	err := d.db.DeleteFrom("friend").Where("uid=? and to_uid=?", uid, toUID).Error
// 	if err != nil {
// 		return err
// 	}
// 	friendKey := fmt.Sprintf("%s%s", CacheKeyFriends, uid)
// 	err = d.ctx.GetRedisConn().SRem(friendKey, toUID)
// 	return err
// }

// 删除好友
// func (d *friendDB) deleteTx(uid, toUID string, tx *gorm.DB) error {
// 	err := tx.Table("friend").Updates(map[string]interface{}{
// 		"is_deleted": 1,
// 		"is_alone":   1,
// 	}).Where("uid=? and to_uid=?", uid, toUID).Error

// 	//err := tx.DeleteFrom("friend").Where("uid=? and to_uid=?", uid, toUID).Error
// 	if err != nil {
// 		return err
// 	}
// 	friendKey := fmt.Sprintf("%s%s", CacheKeyFriends, uid)
// 	err = d.ctx.GetRedisConn().SRem(friendKey, toUID)
// 	return err
// }

// 通过vercode查询好友信息
func (d *friendDB) queryWithVercode(vercode string) (*FriendModel, error) {
	var m FriendModel
	err := d.db.Table("friend").Where("vercode=?", vercode).First(&m).Error
	return &m, err
}

// 通过vercode查询好友信息
func (d *friendDB) queryWithVercodes(vercodes []string) ([]*FriendDetailModel, error) {
	var ms []*FriendDetailModel
	err := d.db.Select("friend.*,IFNULL(user.name,'') name").Table("friend").Joins("user on friend.uid=user.uid").Where("friend.vercode in ?", vercodes).Find(&ms).Error
	return ms, err
}

// 查询某个好友
func (d *friendDB) queryWithUID(uid, toUID string) (*FriendModel, error) {
	var m FriendModel
	err := d.db.Table("friend").Where("uid=? and to_uid=?", uid, toUID).First(&m).Error
	return &m, err
}

// 查询双方好友
func (d *friendDB) queryTwoWithUID(uid, toUID string) ([]*FriendModel, error) {
	var ms []*FriendModel
	err := d.db.Table("friend").Where("(uid=? and to_uid=?) or (uid=? and to_uid=?)", uid, toUID, toUID, uid).Find(&ms).Error
	return ms, err
}

// 查询指定用户uid的在toUids范围内的好友
func (d *friendDB) queryWithToUIDsAndUID(toUids []string, uid string) ([]*FriendModel, error) {
	var ms []*FriendModel
	err := d.db.Table("friend").Where("uid=? and to_uid in ?", uid, toUids).Find(&ms).Error
	return ms, err
}

// 查询uids范围内的用户与toUID是好友的数据
func (d *friendDB) queryWithToUIDAndUIDs(toUID string, uids []string) ([]*FriendModel, error) {
	var ms []*FriendModel
	err := d.db.Table("friend").Where("to_uid=? and uid in ?", toUID, uids).Find(&ms).Error
	return ms, err
}

// QueryFriendsWithKeyword 通过关键字查询自己的好友
func (d *friendDB) QueryFriendsWithKeyword(uid string, keyword string) ([]*DetailModel, error) {
	var ms []*DetailModel
	builder := d.db.Select("friend.id,friend.to_uid,IFNULL(user.name,'') to_name,friend.is_deleted,friend.created_at,friend.updated_at,IFNULL(user_setting.mute,0) mute,IFNULL(user_setting.top,0) top,IFNULL(user_setting.version,0)+friend.version version").Table("friend").Joins("user on friend.to_uid=user.uid").Joins("user_setting on user.uid=user_setting.to_uid and user_setting.uid=friend.uid").Where("friend.uid=?", uid).Order("friend.version + IFNULL(user_setting.version,0) ASC")
	if keyword != "" {
		builder = builder.Where("user.name like ?", "%"+keyword+"%")
	}
	err := builder.Find(&ms).Error
	return ms, err
}

// SyncFriendsOfDeprecated 同步好友
// Deprecated 已废弃，用SyncFriends方法。
func (d *friendDB) SyncFriendsOfDeprecated(version int64, uid string, limit int) ([]*DetailModel, error) {
	var ms []*DetailModel
	builder := d.db.Select("friend.id,IFNULL(friend.vercode,'') vercode,friend.to_uid,IFNULL(user.name,'') to_name,IFNULL(user.category,'') to_category,IFNULL(user.robot,0) robot,IFNULL(user.short_no,'') short_no,IFNULL(friend.remark,'') remark,friend.is_deleted,friend.created_at,friend.updated_at,IFNULL(user_setting.mute,0) mute,IFNULL(user_setting.chat_pwd_on,0) chat_pwd_on,IFNULL(user_setting.blacklist,0) blacklist,IFNULL(user_setting.top,0) top,IFNULL(user_setting.receipt,0) receipt,friend.version + IFNULL(user_setting.version,0) version").Table("friend").Joins("user on friend.to_uid=user.uid").Joins("user_setting on user.uid=user_setting.to_uid and user_setting.uid=friend.uid").Where("friend.uid=?", uid).Order("friend.version + IFNULL(user_setting.version,0) ASC")
	var err error
	if version <= 0 {
		err = builder.Limit(limit).Find(&ms).Error
	} else {
		err = builder.Where("IFNULL(user_setting.version,0) + friend.version > ?", version).Limit(limit).Find(&ms).Error
	}
	return ms, err
}

func (d *friendDB) SyncFriends(version int64, uid string, limit int) ([]*FriendModel, error) {
	var ms []*FriendModel
	builder := d.db.Table("friend").Where("friend.uid=?", uid).Order("friend.version ASC")
	err := builder.Where("friend.version > ?", version).Limit(limit).Find(&ms).Error
	return ms, err
}

// QueryFriends 查询用户的所有好友
func (d *friendDB) QueryFriends(uid string) ([]*DetailModel, error) {
	var ms []*DetailModel
	err := d.db.Select("friend.*,IFNULL(user.name,'') to_name").Table("friend").Joins("user on user.uid=friend.to_uid").Where("friend.uid=? and friend.is_deleted=0", uid).Find(&ms).Error
	return ms, err
}

// QueryFriendsWithUIDs 通过用户id查询好友
func (d *friendDB) QueryFriendsWithUIDs(uid string, toUIDs []string) ([]*FriendDetailModel, error) {
	var ms []*FriendDetailModel
	err := d.db.Select("friend.*,IFNULL(user.name,'') to_name").Table("friend").Joins("user on user.uid=friend.to_uid").Where("friend.uid=? and friend.is_deleted=0 and friend.to_uid in ?", uid, toUIDs).Find(&ms).Error
	return ms, err
}

func (d *friendDB) updateVersionTx(version int64, uid string, toUID string, tx *gorm.DB) error {
	err := tx.Table("friend").Set("version", version).Where("uid=? and to_uid=?", uid, toUID).Error
	return err
}

func (d *friendDB) existBlacklist(uid string, toUID string) (bool, error) {
	var cn int64
	err := d.db.Table("user_setting").Where("((uid=? and to_uid=?) or (uid=? and to_uid=?)) and blacklist=1", uid, toUID, toUID, uid).Count(&cn).Error
	return cn > 0, err
}
func (d *friendDB) insertApplyTx(m *FriendApplyModel, tx *gorm.DB) error {
	err := tx.Table("friend_apply_record").Create(m).Error
	return err
}

func (d *friendDB) insertApply(m *FriendApplyModel) error {
	err := d.db.Table("friend_apply_record").Create(m).Error
	return err
}

func (d *friendDB) queryApplysWithPage(uid string, pageSize, page int) ([]*FriendApplyModel, error) {
	var ms []*FriendApplyModel
	err := d.db.Table("friend_apply_record").Where("uid=?", uid).Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&ms).Error
	return ms, err
}

func (d *friendDB) deleteApplyWithUidAndToUid(uid, toUid string) error {
	err := d.db.Table("friend_apply_record").Where("uid=? and to_uid=?", uid, toUid).Delete(nil).Error
	return err
}
func (d *friendDB) queryApplyWithUidAndToUid(uid, toUid string) (*FriendApplyModel, error) {
	var m FriendApplyModel
	err := d.db.Table("friend_apply_record").Where("uid=? and to_uid=?", uid, toUid).First(&m).Error
	return &m, err
}

func (d *friendDB) updateApply(m *FriendApplyModel) error {
	err := d.db.Table("friend_apply_record").Updates(map[string]interface{}{
		"status": m.Status,
	}).Where("id=?", m.Id).Error
	return err
}

func (d *friendDB) updateApplyTx(m *FriendApplyModel, tx *gorm.DB) error {
	err := tx.Table("friend_apply_record").Updates(map[string]interface{}{
		"status": m.Status,
	}).Where("id=?", m.Id).Error
	return err
}

// DetailModel 好友详情
type DetailModel struct {
	Remark     *string //好友备注
	ToUID      *string // 好友uid
	ToName     *string // 好友名字
	ToCategory *string // 用户分类
	Mute       *int    // 免打扰
	Top        *int    // 置顶
	Version    *int64  // 版本
	Vercode    *string // 验证码 加好友需要
	IsDeleted  *int    // 是否删除
	IsAlone    *int    // 是否为单项好友
	ShortNo    *string //短编号
	ChatPwdOn  *int    // 是否开启聊天密码
	Blacklist  *int    //是否在黑名单
	Receipt    *int    //消息是否回执
	Robot      *int    // 机器人0.否1.是
	Id         *int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

// FriendModel 好友对象
type FriendModel struct {
	UID           *string
	ToUID         *string
	Flag          *int
	Version       *int64
	IsDeleted     *int
	IsAlone       *int // 是否为单项好友
	Vercode       *string
	SourceVercode *string //来源验证码
	Initiator     *int    //1:发起方
	Id            *int64
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// FriendDetailModel 好友资料
type FriendDetailModel struct {
	FriendModel
	Name   *string // 用户名称
	ToName *string //对方用户名称
}

// FriendApplyModel 好友申请记录
type FriendApplyModel struct {
	UID       *string
	ToUID     *string
	Remark    *string
	Token     *string
	Status    *int // 状态 0.未处理 1.通过 2.拒绝
	Id        *int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
