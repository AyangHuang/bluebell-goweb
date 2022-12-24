package models

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamReleasePost 发布帖子请求参数
type ParamReleasePost struct {
	ID       int64  `json:"id,string" db:"post_id"`
	AuthorID int64  `json:"author_id" db:"author_id"`
	Title    string `json:"title" db:"title" binding:"required"`
	Content  string `json:"content" db:"content" binding:"required"`
}
