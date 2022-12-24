package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/utils/snowflake"
)

// ReleasePost 发布帖子
func ReleasePost(post *models.ParamReleasePost) (err error) {

	post.ID = snowflake.GenID()

	// 1. 存入数据库
	if err = mysql.InsertPost(post); err != nil {
		return
	}
	// 2. 存入redis
	err = redis.CreatePost(post.ID)
	return
}
