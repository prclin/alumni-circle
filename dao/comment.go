package dao

import (
	"github.com/prclin/alumni-circle/model"
	"gorm.io/gorm"
)

type CommentDao struct {
	Tx *gorm.DB
}

func NewCommentDao(tx *gorm.DB) *CommentDao {
	return &CommentDao{Tx: tx}
}

func (dao *CommentDao) SelectByBreakId(breakId, parentId uint64, pagination model.Pagination) ([]model.TComment, error) {
	var comments []model.TComment
	sql := "select id, parent_id, account_id, break_id, content, reply_count, like_count, extra, create_time, update_time from comment where break_id=? and parent_id=? limit ?,?"
	err := dao.Tx.Raw(sql, breakId, parentId, (pagination.Page-1)*pagination.Size, pagination.Size).Scan(&comments).Error
	if comments == nil {
		comments = make([]model.TComment, 0, 0)
	}
	return comments, err
}

// IsLiked 是否点赞
func (dao *CommentDao) IsLiked(accountId, commentId uint64) bool {
	var liked bool
	sql := "select count(*) from comment_like where account_id=? and comment_id=?"
	dao.Tx.Raw(sql, accountId, commentId).First(&liked)
	return liked
}
