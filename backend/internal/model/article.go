package model

import "time"

// Article 文章实体：状态 Status（0 草稿 / 1 已发布 / 2 已下架，constants.ArticleStatus）
type Article struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	AuthorID      uint       `gorm:"index;not null" json:"author_id"`
	Title         string     `gorm:"size:256;not null" json:"title"`
	Content       string     `gorm:"type:longtext;not null" json:"content"`
	Summary       string     `gorm:"size:512" json:"summary"`
	CoverURL      string     `gorm:"size:512" json:"cover_url"`
	Status        int        `gorm:"not null;default:0" json:"status"`
	LikeCount     int        `gorm:"not null;default:0" json:"like_count"`
	ViewCount     int        `gorm:"not null;default:0" json:"view_count"`
	FavoriteCount int        `gorm:"not null;default:0" json:"favorite_count"`
	CommentCount  int        `gorm:"not null;default:0" json:"comment_count"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Author User    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Topics []Topic `gorm:"many2many:article_topics;" json:"topics,omitempty"`
}
