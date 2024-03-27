package user

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/chatimxxx/TangSengDaoDaoServerLib/config"
)

// DB DB
type onlineDB struct {
	db  *gorm.DB
	ctx *config.Context
}

// newOnlineDB newOnlineDB
func newOnlineDB(ctx *config.Context) *onlineDB {
	d, err := ctx.DB()
	if err != nil {
		panic(fmt.Sprintf("服务初始化失败   %v", err))
	}
	return &onlineDB{
		db:  d,
		ctx: ctx,
	}
}

// insertOrUpdateUserOnlineTx 插入或更新用户在线信息
func (o *onlineDB) insertOrUpdateUserOnlineTx(m *onlineStatusModel, tx *gorm.DB) error {
	var err error
	if *m.Online == 1 {
		err = tx.Exec("insert into user_online (uid,device_flag,last_online,online,version) values(?,?,?,1,?) ON DUPLICATE KEY UPDATE last_online=VALUES(last_online),online=VALUES(online),updated_at=NOW(),version=VALUES(version)", m.UID, m.DeviceFlag, m.LastOnline, m.Version).Error
	} else {
		err = tx.Exec("insert into user_online (uid,device_flag,last_offline,online,version) values(?,?,?,0,?) ON DUPLICATE KEY UPDATE last_offline=VALUES(last_offline),online=VALUES(online),updated_at=NOW(),version=VALUES(version)", m.UID, m.DeviceFlag, m.LastOffline, m.Version).Error
	}

	return err
}

// queryUserOnlineRecets 查询最近在线的用户(最近是指一小时内在线的,最多查询到1000条)
func (o *onlineDB) queryUserOnlineRecets(uids []string) ([]*onlineStatusWeightModel, error) {
	if len(uids) == 0 {
		return nil, nil
	}
	var ms []*onlineStatusWeightModel
	err := o.db.Select("user_online.*,IFNULL(device_flag.weight,0) weight").Table("user_online").Joins("device_flag on user_online.device_flag=device_flag.device_flag").Where("user_online.uid in ? and ( unix_timestamp(now())-user_online.last_offline<3600*24 or user_online.online=1)", uids).Order("user_online.online, user_online.last_offline DESC").Find(&ms).Error
	onlineStatusMap := map[string]*onlineStatusWeightModel{}
	if len(ms) > 0 {
		for _, m := range ms {
			oldOnline := onlineStatusMap[*m.UID]
			if oldOnline == nil {
				onlineStatusMap[*m.UID] = m
				continue
			}
			replace := false
			if *m.Online == 1 && *oldOnline.Online == 0 {
				replace = true
			}
			if *m.Online == 1 && *oldOnline.Online == 1 && m.Weight > oldOnline.Weight {
				replace = true
			}
			if *m.Online != 1 && *oldOnline.Online != 1 && *m.LastOffline > *oldOnline.LastOffline {
				replace = true
			}
			if replace {
				onlineStatusMap[*m.UID] = m
			}
		}
	}
	newModels := make([]*onlineStatusWeightModel, 0, len(onlineStatusMap))
	for _, value := range onlineStatusMap {
		newModels = append(newModels, value)
	}
	return newModels, err
}

func (o *onlineDB) queryUserLastNewOnlines(uids []string) ([]*onlineStatusWeightModel, error) {
	if len(uids) == 0 {
		return nil, nil
	}
	var ms []*onlineStatusWeightModel
	err := o.db.Select("user_online.*,IFNULL(device_flag.weight,0) weight").Table("user_online").Joins("device_flag on user_online.device_flag=device_flag.device_flag").Where("user_online.uid in ?", uids).Order("user_online.online, user_online.last_offline DESC").Find(&ms).Error
	onlineStatusMap := map[string]*onlineStatusWeightModel{}
	if len(ms) > 0 {
		for _, m := range ms {
			oldOnline := onlineStatusMap[*m.UID]
			if oldOnline == nil {
				onlineStatusMap[*m.UID] = m
				continue
			}
			replace := false
			if *m.Online == 1 && *oldOnline.Online == 0 {
				replace = true
			}
			if *m.Online == 1 && *oldOnline.Online == 1 && m.Weight > oldOnline.Weight {
				replace = true
			}
			if *m.Online != 1 && *oldOnline.Online != 1 && *m.LastOffline > *oldOnline.LastOffline {
				replace = true
			}
			if replace {
				onlineStatusMap[*m.UID] = m
			}

		}
	}
	newModels := make([]*onlineStatusWeightModel, 0, len(onlineStatusMap))
	for _, value := range onlineStatusMap {
		newModels = append(newModels, value)
	}
	return newModels, err
}

func (o *onlineDB) queryOnlinesMoreThan(t time.Duration, limit int) ([]*onlineStatusModel, error) {
	var ms []*onlineStatusModel
	err := o.db.Table("user_online").Where("`online`=1 and unix_timestamp(now()) - last_online>?", t.Seconds()).Limit(limit).Find(&ms).Error
	return ms, err
}

// 查询用户最近在线设备
func (o *onlineDB) queryLastOnlineDeviceWithUID(uid string) (*onlineStatusModel, error) {
	var m onlineStatusModel
	err := o.db.Table("user_online").Where("uid=?", uid).Order("online=1, last_offline DESC").Limit(1).First(&m).Error
	return &m, err
}

func (o *onlineDB) queryOnlineDevice(uid string, deviceFlag config.DeviceFlag) (*onlineStatusModel, error) {
	var m onlineStatusModel
	err := o.db.Table("user_online").Where("uid=? and device_flag=?", uid, deviceFlag.Uint8()).First(&m).Error
	return &m, err
}

func (o *onlineDB) exist(uid string, deviceFlag uint8, online int) (bool, error) {
	var cn int64
	err := o.db.Table("user_online").Where("uid=? and device_flag=? and `online`=?", uid, deviceFlag, online).Count(&cn).Error
	return cn > 0, err
}

// 查询用户在线设备里最大权重的在线状态
func (o *onlineDB) queryOnlineMaxWeightWithUID(uid string) (*onlineStatusModel, error) {
	var m onlineStatusModel
	err := o.db.Select("user_online.*").Table("user_online").Joins("device_flag on user_online.device_flag=device_flag.device_flag").Where("uid=? and online=1", uid).Order("device_flag.weight DESC").Limit(1).First(&m).Error
	return &m, err
}

// 查询在线用户总数量
func (o *onlineDB) queryOnlineCount() (int64, error) {
	var count int64
	err := o.db.Exec("select count(distinct uid) as count from user_online where online=1").Count(&count).Error
	return count, err
}

// OnlineStatusModel 在线状态model
type onlineStatusModel struct {
	UID         *string
	DeviceFlag  *uint8 // 设备标记 0. APP 1.web
	LastOnline  *int   // 最后一次在线时间
	LastOffline *int   // 最后一次离线时间
	Online      *int
	Version     *int64 // 数据版本
	Id          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type onlineStatusWeightModel struct {
	onlineStatusModel
	Weight int // 设备权重
}
