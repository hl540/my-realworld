package data

import "gorm.io/gorm"

// User 用户
type User struct {
	gorm.Model
	Username string `gorm:"column:username;unique;not null;"`
	Password string `gorm:"column:password;not null;"`
	Email    string `gorm:"column:email;unique;not null;"`
	Image    string `gorm:"column:image;"`
	Bio      string `gorm:"column:bio;"`
}

func (u *User) TableName() string {
	return "user"
}

// Follow 关注用户
type Follow struct {
	gorm.Model
	UserID   uint `gorm:"column:user_id;"`
	TargetID uint `gorm:"column:target_id;"`
}

func (f *Follow) TableName() string {
	return "follow"
}

// Article 文章
type Article struct {
	gorm.Model
	Slug        string `gorm:"column:slug;unique"`
	Title       string `gorm:"column:title;"`
	Description string `gorm:"column:description;"`
	Body        string `gorm:"column:body;"`
	AuthorID    uint   `gorm:"column:author_id;"`
}

func (a *Article) TableName() string {
	return "article"
}

// Tag 标签
type Tag struct {
	gorm.Model
	ArticleID uint   `gorm:"column:article_id;"`
	Name      string `gorm:"column:name;"`
}

func (t *Tag) TableName() string {
	return "tag"
}

// Favorite 收藏文章
type Favorite struct {
	gorm.Model
	ArticleID uint `gorm:"column:article_id;"`
	UserID    uint `gorm:"column:user_id;"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}
