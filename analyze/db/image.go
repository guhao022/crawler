package db

import (
	"time"
	"fmt"
	"github.com/num5/logger"
)

type Image struct {
	ID          int       `json:"id"`         //ID
	SourceUrl   string    `json:"source_url"` // 原网址
	Path        string    `json:"path"`       // 图片储存地址
	ReadNum     int       `json:"read_num"`    // 阅读量
	LikeNum     int       `json:"like_num"`    // 点赞数
	CommentNum  int       `json:"comment_num"` // 评论数
	PublishedAt time.Time `json:"published_at" gorm:"default: null"`
	CreatedAt   time.Time `json:"created_at"` // 创建时间
	UpdatedAt   time.Time `json:"updated_at"` // 更新时间
}

func (Image) TableName() string {
	return "spider_image"
}

func (i *Image) Store() error {
	db := conn()
	defer db.Close()


	if err := i.checkJoke(i); err != nil {
		logger.Warn(err.Error())
		return err
	}

	err := db.Create(&i).Error

	if err != nil {
		logger.Warnf("写入内容表失败：%s", err)
		return err
	}

	return nil
}

// 检测是否存在笑话
func (i *Image) checkJoke(c *Image) error {
	if len(c.SourceUrl) == 0 {
		return fmt.Errorf("未找到图片，跳过本地址...")
	}

	db := conn()
	defer db.Close()

	var count int

	err := db.Model(&Image{}).Where("source_url = ?", c.SourceUrl).Count(&count).Error

	if err != nil {
		return fmt.Errorf("图片查重失败：%s", err.Error())
	}

	if count > 0 {
		return fmt.Errorf("图片已经存在，跳过本地址...")
	}

	return nil
}
