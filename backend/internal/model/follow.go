package model

import "time"

// Follow 关注关系：唯一索引保证同一关注对被关注两次
type Follow struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FollowerID uint      `gorm:"uniqueIndex:uk_follow_pair;not null" json:"follower_id"`
	FollowedID uint      `gorm:"uniqueIndex:uk_follow_pair;not null" json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}
