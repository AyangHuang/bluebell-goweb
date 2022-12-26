package mysql

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
)

func InsertPost(post *models.ParamReleasePost) (err error) {

	sqlStr := `insert into post(
	post_id, title, content, author_id)
	values (?, ?, ?, ?)
	`
	_, err = db.Exec(sqlStr, post.ID, post.Title, post.Content, post.AuthorID)
	return
}

func UpdatePostVoteNum(s *[]redis.Z, postVoteNum *[]float64, postIsExpired *[]bool) error {
	sqlStr := `update post set votes_num=? where post_id=?`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}

	for i, length := 0, len(*s); i < length; i++ {
		if (*postIsExpired)[i] {
			id, _ := strconv.ParseFloat((*s)[i].Member.(string), 64)
			if _, err = stmt.Exec(int64((*postVoteNum)[i]), id); err != nil {
				return err
			}
		}
	}

	return nil
}
