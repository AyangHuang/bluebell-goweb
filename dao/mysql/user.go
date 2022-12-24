package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"
)

// CheckUserExist 检查 username 是否已存在
func CheckUserExist(username string) (is bool, err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	// 注意，找了好久的 bug。Get 函数，如果结果集为空，返回 err；如果不为空，也返回 err
	err = db.Get(&count, sqlStr, username)
	if !errors.Is(err, sql.ErrNoRows) {
		if count != 0 {
			return true, nil
		} else {
			return false, ErrorUserExist
		}
	}
	return

}

// FindUserByUsername 通过用户名找用户
func FindUserByUsername(username string, user *models.MysqlUser) (err error) {
	sqlStr := `select user_id, username, password from user where username=?`
	if err = db.Get(user, sqlStr, username); err == nil {
		return nil
	}
	return
}

// InsertUser 插入用户
func InsertUser(user *models.MysqlUser) error {
	sql := `insert into user(user_id, username, password) values(?,?,?)`
	_, err := db.Exec(sql, user.UserID, user.Username, user.Password)
	return err
}

// FindUserNameByUserID 通过用户 id 找用户名
func FindUserNameByUserID(userID int64) (username string, err error) {
	sqlStr := `select username from user where user_id=?`
	err = db.Get(&username, sqlStr, userID)
	return
}
