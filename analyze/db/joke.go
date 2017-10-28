package db

import (
	"fmt"
	"github.com/num5/logger"
	"time"
)

type Joker struct {
	ID          int       `json:"id"`          // ID
	SourceUrl   string    `json:"source_url"`  // 原网址
	Category    string    `json:"category"`    // 类型
	Title       string    `json:"title"`       // 标题
	Topic       string    `json:"topic"`       // 话题
	Content     string    `json:"content"`     // 内容
	ReadNum     int       `json:"read_num"`    // 阅读量
	LikeNum     int       `json:"like_num"`    // 点赞数
	CommentNum  int       `json:"comment_num"` // 评论数
	PublishedAt time.Time `json:"published_at" gorm:"default: null"`
	CreatedAt   time.Time `json:"created_at"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at"` // 更新时间
}

func (Joker) TableName() string {
	return "spider_joker"
}

func (j *Joker) Store() error {
	db := conn()
	defer db.Close()

	if len(j.Content) <= 0 {
		return fmt.Errorf("未找到笑话，跳过本地址...")
	}

	if err := j.checkJoke(j); err != nil {
		logger.Warn(err.Error())
		return err
	}

	err := db.Create(&j).Error

	if err != nil {
		logger.Warnf("写入内容表失败：%s", err)
		return err
	}

	return nil
}

// 检测是否存在笑话
func (j *Joker) checkJoke(c *Joker) error {
	if len(c.SourceUrl) == 0 {
		return fmt.Errorf("未找到笑话，跳过本地址...")
	}

	db := conn()
	defer db.Close()

	var count int

	err := db.Model(&Joker{}).Where("source_url = ?", c.SourceUrl).Count(&count).Error

	if err != nil {
		return fmt.Errorf("笑话查重失败：%s", err.Error())
	}

	if count > 0 {
		return fmt.Errorf("笑话已经存在，跳过本地址...")
	}

	return nil
}
