package service

import (
	"github.com/prclin/alumni-circle/dao"
	_error "github.com/prclin/alumni-circle/error"
	"github.com/prclin/alumni-circle/global"
	"github.com/prclin/alumni-circle/model"
)

// GetBreakComments 获取课间评论
func GetBreakComments(acquirer, breakId, parentId uint64, pagination model.Pagination) ([]model.Comment, error) {
	commentDao := dao.NewCommentDao(global.Datasource)
	//获取基础评论
	tComments, err := commentDao.SelectByBreakId(breakId, parentId, pagination)
	if err != nil {
		global.Logger.Debug(err)
		return nil, _error.InternalServerError
	}

	//填充具体信息
	comments := make([]model.Comment, 0, len(tComments))
	for _, tComment := range tComments {
		liked := commentDao.IsLiked(acquirer, tComment.Id)
		info, err := GetAccountInfo(acquirer, tComment.AccountId)
		if err != nil {
			global.Logger.Warn(err)
		}
		comments = append(comments, model.Comment{TComment: tComment, Liked: liked, AccountInfo: info})
	}

	return comments, err
}
