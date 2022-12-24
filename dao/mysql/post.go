package mysql

import (
	"bluebell/models"
)

func InsertPost(post *models.ParamReleasePost) (err error) {

	sqlStr := `insert into post(
	post_id, title, content, author_id)
	values (?, ?, ?, ?)
	`
	_, err = db.Exec(sqlStr, post.ID, post.Title, post.Content, post.AuthorID)
	return
}
