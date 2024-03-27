package webhook

import "gorm.io/gorm"

// DB DB
type DB struct {
	db *gorm.DB
}

// NewDB NewDB
func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// GetThirdName 获取三个名字 （常用名字，好友备注，群内名字） （TODO: 此方法不应该直接写sql 应该调用各模块的server来获取数据）
func (db *DB) GetThirdName(fromUID string, toUID string, groupNo string) (string, string, string, error) {
	if fromUID == "" {
		return "", "", "", nil
	}
	var name string        // 常用名
	var remark string      // 好友备注
	var nameInGroup string // 群内备注

	if toUID == "" && groupNo == "" {
		err := db.db.Select("name").Table("`user`").Where("uid=?", fromUID).First(&name).Error
		if err != nil {
			return "", "", "", nil
		}
	} else if toUID != "" && groupNo == "" {
		var nameStruct struct {
			Name   string
			Remark string
		}
		builder := db.db.Exec("select `user`.name,IFNULL(friend.remark,'') remark from `user` left join friend on `user`.uid=friend.to_uid and friend.uid=? where `user`.uid=? ", toUID, fromUID)
		err := builder.First(&nameStruct).Error
		if err != nil {
			return "", "", "", err
		}
		name = nameStruct.Name
		remark = nameStruct.Remark
	} else if toUID == "" && groupNo != "" {
		var nameStruct struct {
			Name        string
			NameInGroup string
		}
		err := db.db.Exec("select `user`.name,IFNULL(group_member.remark,'') name_in_group from `user` left join group_member on group_member.group_no=?  and `user`.uid=group_member.uid and group_member.is_deleted=0 where `user`.uid=? ", groupNo, fromUID).First(&nameStruct).Error
		if err != nil {
			return "", "", "", err
		}
		name = nameStruct.Name
		nameInGroup = nameStruct.NameInGroup
	} else if toUID != "" && groupNo != "" {
		var nameStruct struct {
			Name        string
			Remark      string
			NameInGroup string
		}
		err := db.db.Exec("select `user`.name,IFNULL(group_member.remark,'') name_in_group,IFNULL(friend.remark ,'') remark from `user` left join group_member on  group_member.group_no=?  and `user`.uid=group_member.uid and group_member.is_deleted=0 left join friend on `user`.uid=friend.to_uid and `user`.uid=? and friend.uid=? where `user`.uid=?", groupNo, fromUID, toUID, fromUID).First(&nameStruct).Error
		if err != nil {
			return "", "", "", err
		}
		name = nameStruct.Name
		nameInGroup = nameStruct.NameInGroup
		remark = nameStruct.Remark
	}
	return name, remark, nameInGroup, nil
}

// GetGroupName 获取群名
func (db *DB) GetGroupName(groupNo string) (string, error) {
	var name string
	err := db.db.Select("name").Table("`group`").Where("group_no=?", groupNo).First(&name).Error
	return name, err
}
