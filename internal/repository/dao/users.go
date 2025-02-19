package dao

import (
	"fmt"

	"go-chat/internal/repository/model"
)

type UsersDao struct {
	*BaseDao
}

func NewUserDao(baseDao *BaseDao) *UsersDao {
	return &UsersDao{BaseDao: baseDao}
}

// Create 创建数据
func (dao *UsersDao) Create(user *model.Users) (*model.Users, error) {
	if err := dao.Db().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindById ID查询
func (dao *UsersDao) FindById(uid int) (*model.Users, error) {

	if uid == 0 {
		return nil, fmt.Errorf("uid is empty")
	}

	user := &model.Users{}

	if err := dao.Db().Where(&model.Users{Id: uid}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindByMobile 手机号查询
func (dao *UsersDao) FindByMobile(mobile string) (*model.Users, error) {

	if len(mobile) == 0 {
		return nil, fmt.Errorf("mobile is empty")
	}

	user := &model.Users{}

	if err := dao.Db().Where(&model.Users{Mobile: mobile}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// IsMobileExist 判断手机号是否存在
func (dao *UsersDao) IsMobileExist(mobile string) bool {

	if len(mobile) == 0 {
		return false
	}

	user := &model.Users{}

	rowsAffects := dao.Db().Select("id").Where(&model.Users{Mobile: mobile}).First(user).RowsAffected

	return rowsAffects != 0
}
