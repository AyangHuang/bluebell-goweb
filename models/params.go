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
	ID       int64  `json:"id,string"`
	AuthorID int64  `json:"author_id"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required"`              // 贴子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1)还是反对票(-1)取消投票（0）
}
