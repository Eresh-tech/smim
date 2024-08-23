package repo

import (
	"errors"
	"gim/internal/business/domain/user/model"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const (
	UserKey      = "user:"
	UserExpire   = 48 * time.Hour //新用户两天没登录会回收账户
	UserPhone    = "phone:"
	UserNickName = "nickname:"
)

type userDao struct{}

var UserDao = new(userDao)

// Add 插入一条用户信息
func (u *userDao) Add(user *model.User) (int64, error) {
	t := time.Now()
	user.CreateTime = t
	user.UpdateTime = t
	err := u.Save(user)
	if err != nil {
		return user.Id, gerrors.WrapError(err)
	}
	return user.Id, nil
}

// Get 获取用户信息time.Now()
func (u *userDao) Get(userId int64) (*model.User, error) {
	var user = model.User{Id: userId}
	err := db.RedisUtil.Get(UserKey+strconv.FormatInt(userId, 10), &user)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, gerrors.WrapError(err)
	}
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return &user, err
}

// Save 保存, user id->user, phone -> user id(if phone exists), nickname->user id(if nickname exists)
func (u *userDao) Save(user *model.User) error {
	err := db.RedisUtil.Set(UserKey+strconv.FormatInt(user.Id, 10), user, 0)
	if err != nil {
		return gerrors.WrapError(err)
	}

	//如果phone有效，记录phone -> user id
	if user.PhoneNumber == "" {
		err := db.RedisUtil.Set(UserPhone+user.PhoneNumber, user.Id, 0)
		if err != nil {
			return gerrors.WrapError(err)
		}
	}

	//如果nickname有效，记录nickname -> user id
	if user.PhoneNumber == "" {
		err := db.RedisUtil.Set(UserNickName+user.Nickname, user.Id, 0)
		if err != nil {
			return gerrors.WrapError(err)
		}
	}
	return nil
}

// GetByPhoneNumber 根据手机号获取用户信息
func (u *userDao) GetByPhoneNumber(phoneNumber string) (*model.User, error) {
	uid_str := ""
	err := db.RedisUtil.Get(UserPhone+phoneNumber, uid_str)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	uid, err := strconv.ParseInt(uid_str, 10, 64)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return u.Get(uid)
}

// GetByPhoneNumber 根据手机号获取用户信息
func (u *userDao) GetByNickName(nickname string) (*model.User, error) {
	uid_str := ""
	err := db.RedisUtil.Get(UserNickName+nickname, uid_str)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}

	uid, err := strconv.ParseInt(uid_str, 10, 64)
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return u.Get(uid)
}

// GetByIds 获取用户信息, TODO:mget
func (u *userDao) GetByIds(userIds []int64) ([]model.User, error) {
	var users []model.User
	for _, userId := range userIds {
		user, err := u.Get(userId)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

// Search 查询用户
func (u *userDao) Search(key string) ([]model.User, error) {
	var users []model.User
	user, err := u.GetByPhoneNumber(key)
	if err != nil {
		return nil, err
	}
	users = append(users, *user)

	user, err = u.GetByNickName(key)
	if err != nil {
		return nil, err
	}
	users = append(users, *user)
	return users, nil
}
