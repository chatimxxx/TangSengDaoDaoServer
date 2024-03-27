package group

import (
	"fmt"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/common"
	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
	"github.com/gocraft/dbr/v2"
	"gorm.io/gorm"
	"time"
)

// DB DB
type DB struct {
	ctx *config.Context
	db  *gorm.DB
}

// NewDB NewDB
func NewDB(ctx *config.Context) *DB {
	db, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &DB{
		ctx: ctx,
		db:  db,
	}
}

// InsertTx 插入群信息（含事务）
func (d *DB) InsertTx(m *Model, tx *gorm.DB) error {
	err := tx.Table("group").Create(m).Error
	return err
}

// Insert 插入群信息
func (d *DB) Insert(m *Model) error {
	err := d.db.Table("group").Create(m).Error
	return err
}

// 修改群类型
func (d *DB) UpdateGroupTypeTx(groupNo string, groupType GroupType, tx *gorm.DB) error {
	err := tx.Table("group").Update("group_type", int(groupType)).Where("group_no=?", groupNo).Error
	return err
}

// 修改群类型
func (d *DB) UpdateGroupType(groupNo string, groupType GroupType) error {
	err := d.db.Table("group").Update("group_type", int(groupType)).Where("group_no=?", groupNo).Error
	return err
}

// InsertMemberTx 插入群成员信息(带事务)
func (d *DB) InsertMemberTx(m *MemberModel, tx *gorm.DB) error {
	err := tx.Table("group_member").Create(m).Error
	return err
}

// InsertMember 插入群成员信息
func (d *DB) InsertMember(m *MemberModel) error {
	err := d.db.Table("group_member").Create(m).Error
	return err
}

// DeleteMemberTx 删除群成员
func (d *DB) DeleteMemberTx(groupNo string, uid string, version int64, tx *gorm.DB) error {
	err := tx.Table("group_member").Update("is_deleted", 1).Update("version", version).Where("group_no=? and uid=?", groupNo, uid).Error
	return err
}

// DeleteMember 删除群成员
func (d *DB) DeleteMember(groupNo string, uid string, version int64) error {
	err := d.db.Table("group_member").Update("is_deleted", 1).Update("version", version).Where("group_no=? and uid=?", groupNo, uid).Error
	return err
}

// QuerySecondOldestMember 查询群里第二长老
func (d *DB) QuerySecondOldestMember(groupNo string) (*MemberModel, error) {
	var m MemberModel
	err := d.db.Table("group_member").Where("group_no=? and role<>? and is_deleted=0", groupNo, MemberRoleCreator).Order("created_at ASC").First(&m).Error
	return &m, err
}

// 通过vercode查询某个群成员
func (d *DB) queryMemberWithVercode(vercode string) (*MemberModel, error) {
	var m MemberModel
	err := d.db.Table("group_member").Where("vercode=?", vercode).First(&m).Error
	return &m, err
}

// 通过vercode查询某个群成员
func (d *DB) queryMemberWithVercodes(vercodes []string) ([]*MemberGroupDetailModel, error) {
	var ms []*MemberGroupDetailModel
	err := d.db.Select("group_member.*,IFNULL(`group`.name,'') group_name").Table("group_member").Joins("group on `group`.group_no=group_member.group_no").Where("group_member.vercode in ?", vercodes).Find(&ms).Error
	return ms, err
}

// QueryIsGroupManagerOrCreator 是否是群管理者或创建者
func (d *DB) QueryIsGroupManagerOrCreator(groupNo string, uid string) (bool, error) {
	var count int64
	err := d.db.Table("group_member").Where("group_no=? and uid=? and is_deleted=0 and (role=? or role=?)", groupNo, uid, MemberRoleCreator, MemberRoleManager).Count(&count).Error
	return count > 0, err
}

// QueryIsGroupCreator 是否是群创建者
func (d *DB) QueryIsGroupCreator(groupNo string, uid string) (bool, error) {
	var count int64
	err := d.db.Table("group_member").Where("group_no=? and uid=? and is_deleted=0 and role=?", groupNo, uid, MemberRoleCreator).Count(&count).Error
	return count > 0, err
}

// QueryGroupManagerOrCreatorUIDS 查询管理者或创建者的uid
func (d *DB) QueryGroupManagerOrCreatorUIDS(groupNo string) ([]string, error) {
	var uids []string
	err := d.db.Select("uid").Table("group_member").Where("group_no=? and is_deleted=0 and (role=? or role=?)", groupNo, MemberRoleCreator, MemberRoleManager).Find(&uids).Error
	return uids, err
}

func (d *DB) queryGroupMemberMaxVersion(groupNo string) (int64, error) {
	var version int64
	err := d.db.Select("IFNULL(max(version),0)").Table("group_member").Where("group_no=?", groupNo).First(&version).Error
	return version, err
}

// UpdateMemberRoleTx 更新群成员角色
func (d *DB) UpdateMemberRoleTx(groupNo string, uid string, role int, version int64, tx *gorm.DB) error {
	err := tx.Table("group_member").Update("role", role).Update("version", version).Where("group_no=? and uid=? and is_deleted=0", groupNo, uid).Error
	return err
}

// updateMemberForbiddenExpirTimeTx 修改成员禁言时长
func (d *DB) updateMemberForbiddenExpirTimeTx(groupNo string, uid string, time int, version int64, tx *gorm.DB) error {
	err := tx.Table("group_member").Set("forbidden_expir_time", time).Set("version", version).Where("group_no=? and uid=? and is_deleted=0", groupNo, uid).Error
	return err
}

// UpdateMembersToManager 更新指定群成员为管理员
func (d *DB) UpdateMembersToManager(groupNo string, members []string, version int64) error {
	if len(members) <= 0 {
		return nil
	}
	err := d.db.Table("group_member").Set("role", MemberRoleManager).Set("version", version).Where("group_no=? and uid in ? and is_deleted=0", groupNo, members).Error
	return err
}

// UpdateManagersToMember 更新指定管理员为普通成员
func (d *DB) UpdateManagersToMember(groupNo string, members []string, version int64) error {
	if len(members) <= 0 {
		return nil
	}
	err := d.db.Table("group_member").Set("role", MemberRoleCommon).Set("version", version).Where("group_no=? and uid in ? and is_deleted=0", groupNo, members).Error
	return err
}

// ExistMember 群成员是否在群内
func (d *DB) ExistMember(uid string, groupNo string) (bool, error) {
	var count int64
	err := d.db.Table("group_member").Where("group_no=? and uid=? and is_deleted=0", groupNo, uid).Count(&count).Error
	return count > 0, err
}

func (d *DB) existMembers(groupNos []string, uid string) ([]string, error) {
	var results []string
	err := d.db.Select("group_no").Table("group_member").Where("group_no in ? and uid=? and is_deleted=0", groupNos, uid).Find(&results).Error
	return results, err
}

// ExistMemberDelete 存在已删除的群成员数据
func (d *DB) ExistMemberDelete(uid string, groupNo string) (bool, error) {
	var count int64
	err := d.db.Table("group_member").Where("group_no=? and uid=? and is_deleted=1", groupNo, uid).Count(&count).Error
	return count > 0, err
}

// UpdateMemberTx 更新成员信息
func (d *DB) UpdateMemberTx(member *MemberModel, tx *gorm.DB) error {
	err := tx.Table("group_member").Updates(map[string]interface{}{
		"remark":     member.Remark,
		"role":       member.Role,
		"version":    member.Version,
		"is_deleted": member.IsDeleted,
		"invite_uid": member.InviteUID,
	}).Where("group_no=? and uid=?", member.GroupNo, member.UID).Error
	return err
}

// recoverMemberTx 恢复成员信息
func (d *DB) recoverMemberTx(member *MemberModel, tx *gorm.DB) error {
	err := tx.Table("group_member").Updates(map[string]interface{}{
		"remark":     member.Remark,
		"role":       member.Role,
		"version":    member.Version,
		"is_deleted": 0,
		"invite_uid": member.InviteUID,
		"created_at": dbr.Expr("Now()"),
	}).Where("group_no=? and uid=?", member.GroupNo, member.UID).Error
	return err
}

// UpdateMember 更新群成员
func (d *DB) UpdateMember(member *MemberModel) error {
	err := d.db.Table("group_member").Updates(map[string]interface{}{
		"remark":               member.Remark,
		"role":                 member.Role,
		"version":              member.Version,
		"is_deleted":           member.IsDeleted,
		"invite_uid":           member.InviteUID,
		"forbidden_expir_time": member.ForbiddenExpirTime,
	}).Where("group_no=? and uid=?", member.GroupNo, member.UID).Error
	return err
}

// 修改群成员状态
func (d *DB) updateMembersStatus(version int64, groupNo string, status int, uids []string) error {
	err := d.db.Table("group_member").Updates(map[string]interface{}{
		"status":  status,
		"version": version,
	}).Where("group_no=? and uid in ?", groupNo, uids).Error
	return err
}

// QueryWithGroupNo 根据群编号查询群信息
func (d *DB) QueryWithGroupNo(groupNo string) (*Model, error) {
	var model *Model
	err := d.db.Table("`group`").Where("group_no=?", groupNo).Find(&model).Error
	return model, err
}

// QueryWithGroupNo 根据群编号查询群信息
func (d *DB) QueryWithGroupNos(groupNos []string) ([]*Model, error) {
	var models []*Model
	err := d.db.Table("`group`").Where("group_no in ?", groupNos).Find(&models).Error
	return models, err
}

func (d *DB) queryUserSupers(uid string) ([]*Model, error) {
	var models []*Model
	err := d.db.Select("`group`.*").Table("group_member").Joins("group on group.group_no=group_member.group_no").Where("group.group_type=? and group.status=? and group_member.is_deleted=0 and group_member.uid=?", GroupTypeSuper, GroupStatusNormal, uid).Find(&models).Error
	return models, err
}

// UpdateTx 更新群信息（带事务）
func (d *DB) UpdateTx(model *Model, tx *gorm.DB) error {
	err := tx.Table("group").Updates(map[string]interface{}{
		"name":      model.Name,
		"notice":    model.Notice,
		"creator":   model.Creator,
		"status":    model.Status,
		"version":   model.Version,
		"forbidden": model.Forbidden,
		"invite":    model.Invite,
	}).Where("id=?", model.Id).Error
	return err
}

// Update 更新群信息
func (d *DB) Update(model *Model) error {
	err := d.db.Table("group").Updates(map[string]interface{}{
		"name":                   model.Name,
		"notice":                 model.Notice,
		"creator":                model.Creator,
		"status":                 model.Status,
		"version":                model.Version,
		"forbidden":              model.Forbidden,
		"invite":                 model.Invite,
		"forbidden_add_friend":   model.ForbiddenAddFriend,
		"allow_view_history_msg": model.AllowViewHistoryMsg,
	}).Where("id=?", model.Id).Error
	return err
}

func (d *DB) updateAvatar(avatar string, groupNo string) error {
	err := d.db.Table("group").Updates(map[string]interface{}{
		"avatar":           avatar,
		"is_upload_avatar": 1,
	}).Where("group_no=?", groupNo).Error
	return err
}

// QueryDetailWithGroupNo 查询群详情
func (d *DB) QueryDetailWithGroupNo(groupNo string, uid string) (*DetailModel, error) {
	var ms *DetailModel
	err := d.db.Select("`group`.*,IFNULL(group_setting.version,0) + `group`.version  version,IFNULL(group_setting.chat_pwd_on,0) chat_pwd_on,IFNULL(group_setting.mute,0) mute,IFNULL(group_setting.top,0) top,IFNULL(group_setting.show_nick,0) show_nick,IFNULL(group_setting.save,0) save,IFNULL(group_setting.revoke_remind,0) revoke_remind,IFNULL(group_setting.revoke_remind,1) revoke_remind,IFNULL(group_setting.join_group_remind,0) join_group_remind,IFNULL(group_setting.screenshot,1) screenshot,IFNULL(group_setting.receipt,1) receipt,IFNULL(group_setting.flame,0) flame,IFNULL(group_setting.flame_second,0) flame_second,IFNULL(group_setting.remark,'') remark").Table("`group`").Joins("`group_setting` on `group`.group_no=group_setting.group_no and group_setting.uid=?").Where("`group`.group_no=?", uid, groupNo).Find(&ms).Error
	return ms, err
}

// QueryDetailWithGroupNos 查询群集合
func (d *DB) QueryDetailWithGroupNos(groupNos []string, uid string) ([]*DetailModel, error) {
	if len(groupNos) <= 0 {
		return nil, nil
	}
	var ms []*DetailModel
	err := d.db.Select("`group`.*,IFNULL(group_setting.version,0) + `group`.version  version,IFNULL(group_setting.chat_pwd_on,0) chat_pwd_on,IFNULL(group_setting.mute,0) mute,IFNULL(group_setting.top,0) top,IFNULL(group_setting.show_nick,0) show_nick,IFNULL(group_setting.save,0) save,IFNULL(group_setting.revoke_remind,0) revoke_remind,IFNULL(group_setting.revoke_remind,1) revoke_remind,IFNULL(group_setting.join_group_remind,0) join_group_remind,IFNULL(group_setting.screenshot,1) screenshot,IFNULL(group_setting.receipt,1) receipt,IFNULL(group_setting.flame,0) flame,IFNULL(group_setting.flame_second,0) flame_second,IFNULL(group_setting.remark,'') remark").Table("`group`").Joins("`group_setting` on `group`.group_no=group_setting.group_no and group_setting.uid=?").Where("`group`.group_no in ?", uid, groupNos).Find(&ms).Error
	return ms, err
}

// QueryGroupsWithGroupNos 通过群ID查询一批群信息
func (d *DB) QueryGroupsWithGroupNos(groupNos []string) ([]*Model, error) {
	if len(groupNos) <= 0 {
		return nil, nil
	}
	var ms []*Model
	err := d.db.Table("`group`").Where("group_no in ?", groupNos).Find(&ms).Error
	return ms, err
}

// QueryMemberWithUID 查询群成员
func (d *DB) QueryMemberWithUID(uid string, groupNo string) (*MemberModel, error) {
	var ms *MemberModel
	err := d.db.Table("group_member").Where("uid=? and group_no=? and is_deleted=0", uid, groupNo).Find(&ms).Error
	return ms, err
}

// QueryMembersWithUids 查询群内的指定成员
func (d *DB) QueryMembersWithUids(uids []string, groupNo string) ([]*MemberModel, error) {
	if len(uids) == 0 {
		return nil, nil
	}
	var ms []*MemberModel
	err := d.db.Table("group_member").Where("uid in ? and group_no=? and is_deleted=0", uids, groupNo).Find(&ms).Error
	return ms, err
}

// QueryMembersWithStatus 通过成员状态查询成员
func (d *DB) QueryMembersWithStatus(groupNo string, status int) ([]*MemberModel, error) {
	var ms []*MemberModel
	err := d.db.Table("group_member").Where("group_no=? and status=?", groupNo, status).Find(&ms).Error
	return ms, err
}

// SyncMembers 同步群成员
func (d *DB) SyncMembers(groupNo string, version int, limit int) ([]*MemberDetailModel, error) {
	var ms []*MemberDetailModel
	builder := d.db.Select("group_member.id,group_member.vercode,group_member.uid,group_member.status,group_member.group_no,group_member.remark,group_member.role,IFNULL(user.name,'') name,IFNULL(user.username,'') username,group_member.is_deleted,group_member.robot,group_member.version,group_member.invite_uid,group_member.forbidden_expir_time,group_member.created_at,group_member.updated_at").Table("group_member").Joins("user on group_member.uid=user.uid").Where("group_member.group_no=?", groupNo).Order("group_member.version DESC")
	var err error
	if version <= 0 {
		err = builder.Limit(limit).Find(&ms).Error
	} else {
		err = builder.Where("group_member.version > ?", version).Limit(limit).Find(&ms).Error
	}
	return ms, err
}

// 通过名字关键字查询成员列表
func (d *DB) queryMembersWithKeyword(groupNo string, loginUID string, keyword string, page int, limit int) ([]*MemberDetailModel, error) {
	var builder *gorm.DB
	if keyword != "" {
		builder = d.db.Select("group_member.id,group_member.vercode,group_member.uid,group_member.status,group_member.group_no,group_member.remark,group_member.role,IFNULL(user.name,'') name,IFNULL(user.username,'') username,group_member.is_deleted,group_member.robot,group_member.version,group_member.invite_uid,group_member.forbidden_expir_time,group_member.created_at,group_member.updated_at").Table("group_member").Joins("user on group_member.uid=user.uid").Joins("user_setting on user_setting.uid='?' and user_setting.to_uid=group_member.uid", loginUID).Where("group_member.group_no=? and group_member.is_deleted=0 and (group_member.remark like ? or user.name like ? or user_setting.remark like ?)", groupNo, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Order("group_member.created_at DESC")
	} else {
		builder = d.db.Select("group_member.id,group_member.vercode,group_member.uid,group_member.status,group_member.group_no,group_member.remark,group_member.role,IFNULL(user.name,'') name,IFNULL(user.username,'') username,group_member.is_deleted,group_member.robot,group_member.version,group_member.invite_uid,group_member.forbidden_expir_time,group_member.created_at,group_member.updated_at").Table("group_member").Joins("user on group_member.uid=user.uid").Where("group_member.group_no=? and group_member.is_deleted=0", groupNo).Order(fmt.Sprintf("group_member.role=%d DESC", MemberRoleCreator)).Order(fmt.Sprintf("group_member.role=%d DESC", MemberRoleManager)).Order("group_member.created_at DESC")
	}
	var ms []*MemberDetailModel
	var err error
	err = builder.Offset((page - 1) * limit).Limit(limit).Find(&ms).Error
	return ms, err
}

func (d *DB) queryMembersWithGroupNo(groupNo string) ([]*MemberDetailModel, error) {
	var ms []*MemberDetailModel
	err := d.db.Table("group_member").Select("group_member.id,group_member.vercode,group_member.uid,group_member.status,group_member.group_no,group_member.remark,group_member.role,IFNULL(user.name,'') name,group_member.is_deleted,group_member.version,group_member.created_at,group_member.updated_at").Joins("user on group_member.uid=user.uid").Where("group_member.group_no=? and group_member.is_deleted=0", groupNo).Find(&ms).Error
	return ms, err
}

func (d *DB) queryBlacklistMemberUIDsWithGroupNo(groupNo string) ([]string, error) {
	var uids []string
	err := d.db.Table("group_member").Select("group_member.uid").Where("group_member.group_no=? and group_member.is_deleted=0 and status=?", groupNo, common.GroupMemberStatusBlacklist).Find(&uids).Error
	return uids, err
}

// 查询在线成员数量
func (d *DB) queryMemberOnlineCount(groupNo string) (int64, error) {
	var count int64
	err := d.db.Table("group_member").Select("count(DISTINCT user_online.uid) left join user_online on group_member.uid=user_online.uid").Where("group_no=? and group_member.is_deleted=0 and user_online.online=1", groupNo).Count(&count).Error
	return count, err
}

// QueryMembersFirstNine 查询最先加入群聊的九为群成员
func (d *DB) QueryMembersFirstNine(groupNo string) ([]*MemberModel, error) {
	var ms []*MemberModel
	err := d.db.Table("group_member").Where("group_no=? and is_deleted=0", groupNo).Order("created_at").Limit(9).Find(&ms).Error
	return ms, err
}

// QueryMembersFirstNineExclude 查询最先加入群聊的九位群成员 【excludeUIDs】为排除的用户
func (d *DB) QueryMembersFirstNineExclude(groupNo string, excludeUIDs []string) ([]*MemberModel, error) {
	if len(excludeUIDs) <= 0 {
		return d.QueryMembersFirstNine(groupNo)
	}
	var ms []*MemberModel
	err := d.db.Table("group_member").Where("group_no=? and is_deleted=0 and uid not in ?", groupNo, excludeUIDs).Order("created_at").Limit(9).Find(&ms).Error
	return ms, err
}

// 成员是否在最先加入的9位成员内
func (d *DB) membersInFirstNine(groupNo string, uids []string) (bool, error) {
	if len(uids) == 0 {
		return false, nil
	}
	var count int64
	err := d.db.Exec("select count(*) from (select uid from group_member where group_no=? and is_deleted=0 order by created_at asc limit 9) t where t.uid in ?", groupNo, uids).Count(&count).Error
	return count > 0, err
}

// QueryMemberCount 查询群成员数量
func (d *DB) QueryMemberCount(groupNo string) (int64, error) {
	var count int64
	err := d.db.Table("group_member").Where("group_no=? and is_deleted=0", groupNo).Count(&count).Error
	return count, err
}

// 查询群总数
func (d *DB) queryGroupCount() (int64, error) {
	var count int64
	err := d.db.Table("`group`").Count(&count).Error
	return count, err
}

// 查询某天的新建群数量
func (d *DB) queryCreatedCountWithDate(date string) (int64, error) {
	var count int64
	err := d.db.Table("`group`").Where("date_format(created_at,'%Y-%m-%d')=?", date).Count(&count).Error
	return count, err
}

// querySavedGroups 查询我保存的群
func (d *DB) querySavedGroups(uid string) ([]*DetailModel, error) {
	var ms []*DetailModel
	err := d.db.Table("`group`").Select("`group`.*,IFNULL(group_setting.version,0) + `group`.version  version,IFNULL(group_setting.chat_pwd_on,0) chat_pwd_on,IFNULL(group_setting.mute,0) mute,IFNULL(group_setting.top,0) top,IFNULL(group_setting.show_nick,0) show_nick,IFNULL(group_setting.save,0) save,IFNULL(group_setting.remark,'') remark left join `group_setting` on `group`.group_no=group_setting.group_no").Where("`group_setting`.save=1 and `group_setting`.uid=?", uid).First(&ms).Error
	return ms, err
}

// 查询某个用户参与的所有群
func (d *DB) queryGroupsWithMemberUID(memberUID string) ([]*Model, error) {
	var ms []*Model
	err := d.db.Table("`group_member`").Select("distinct `group`.* left join `group` on `group`.group_no=group_member.group_no").Where("group_member.uid=? and group_member.is_deleted=0", memberUID).Find(&ms).Error
	return ms, err
}

// 查询禁言时长到期成员
func (d *DB) queryForbiddenExpirationTimeMembers(limit int) ([]*MemberModel, error) {
	var models []*MemberModel
	err := d.db.Table("group_member").Where("forbidden_expir_time <>0 and unix_timestamp(now())-forbidden_expir_time>0").Limit(limit).Find(&models).Error
	return models, err
}

// 查询群头像是否已被群主更新过
func (d *DB) queryGroupAvatarIsUpload(groupNo string) (int, error) {
	var result int
	err := d.db.Select("is_upload_avatar").Table("`group`").Where("group_no=?", groupNo).First(&result).Error
	return result, err
}

// ---------- model ----------

// DetailModel 群详情
type DetailModel struct {
	Model
	Mute            *int    // 免打扰
	Top             *int    // 置顶
	ShowNick        *int    // 显示昵称
	Save            *int    // 是否保存
	ChatPwdOn       *int    //是否开启聊天密码
	RevokeRemind    *int    //撤回提醒
	JoinGroupRemind *int    // 进群提醒
	Screenshot      *int    //截屏通知
	Receipt         *int    //消息是否回执
	Flame           *int    // 是否开启阅后即焚
	FlameSecond     *int    // 阅后即焚秒数
	Remark          *string // 群备注
}

// Model 群db model
type Model struct {
	GroupNo             *string // 群编号
	GroupType           *int    // 群类型 0.普通群 1.超大群
	Name                *string // 群名称
	Avatar              *string // 群头像
	Notice              *string // 群公告
	Creator             *string // 创建者uid
	Status              *int    // 群状态
	Version             *int64  // 版本号
	Forbidden           *int    // 是否全员禁言
	Invite              *int    // 是否开启邀请确认 0.否 1.是
	ForbiddenAddFriend  *int    //群内禁止加好友
	AllowViewHistoryMsg *int    // 是否允许新成员查看历史消息
	Category            *string // 群分类
	Id                  *int64
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}

// MemberModel 成员model
type MemberModel struct {
	GroupNo            *string // 群编号
	UID                *string // 成员uid
	Remark             *string // 成员备注
	Role               *int    // 成员角色 1. 创建者	 2.管理员
	Version            *int64
	Status             *int    // 1.正常 2.黑名单
	Vercode            *string //验证码
	IsDeleted          *int    // 是否删除
	InviteUID          *string // 邀请者
	Robot              *int    // 机器人
	ForbiddenExpirTime *int64  // 禁言时长
	Id                 *int64
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

// MemberDetailModel 成员详情model
type MemberDetailModel struct {
	UID                *string // 成员uid
	GroupNo            *string // 群编号
	Name               *string // 群成员名称
	Remark             *string // 成员备注
	Role               *int    // 成员角色
	Version            *int64
	Vercode            *string //验证码
	InviteUID          *string // 邀请人
	IsDeleted          *int    // 是否删除
	Status             *int    // 1.正常 2.黑名单
	Username           *string
	Robot              *int   // 机器人标识0.否1.是
	ForbiddenExpirTime *int64 // 禁言时长
	Id                 *int64
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}

type MemberGroupDetailModel struct {
	GroupName *string // 群名称
	MemberModel
}
