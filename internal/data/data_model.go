package data

import (
	"time"
)

// User 用户
type User struct {
	Id        int64     `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Username  string    `xorm:"notnull unique"`
	Password  string    `xorm:"notnull"`
	Email     string    `xorm:"notnull unique"`
	Image     string
	Bio       string
}

func (u *User) TableName() string {
	return "user"
}

// Follow 关注用户
type Follow struct {
	Id        int64     `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	UserId    int64
	TargetId  int64
}

func (f *Follow) TableName() string {
	return "follow"
}

// Article 文章
type Article struct {
	Id          int64     `xorm:"pk autoincr"`
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
	DeletedAt   time.Time `xorm:"deleted"`
	Slug        string
	Title       string
	Description string
	Body        string
	AuthorId    int64
}

func (a *Article) TableName() string {
	return "article"
}

// Tag 标签
type Tag struct {
	Id        int64     `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	ArticleId int64
	Name      string
}

func (t *Tag) TableName() string {
	return "tag"
}

// Favorite 收藏文章
type Favorite struct {
	Id        int64     `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	ArticleId int64
	UserId    int64
}

func (f *Favorite) TableName() string {
	return "favorite"
}
