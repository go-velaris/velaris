package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

type boardPolicy struct {
	engine *xorm.Engine
}

func NewBoardPolicy(engine *xorm.Engine) Policy {
	return &boardPolicy{engine}
}

func userOwnsBoard(engine *xorm.Engine, ctxUser middleware.CtxUser, id int64) (bool, error) {
	exists, err := engine.
		Where("id = ? AND user_id = ?", id, ctxUser.ID).
		Exist()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (bp *boardPolicy) CanView(ctxUser middleware.CtxUser, boardId int64) (bool, error) {
	return userOwnsBoard(bp.engine, ctxUser, boardId)
}
func (bp *boardPolicy) CanCreate(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return true, nil
}
func (bp *boardPolicy) CanUpdate(ctxUser middleware.CtxUser, boardId int64) (bool, error) {
	return userOwnsBoard(bp.engine, ctxUser, boardId)
}

func (bp *boardPolicy) CanDelete(ctxUser middleware.CtxUser, boardId int64) (bool, error) {
	return userOwnsBoard(bp.engine, ctxUser, boardId)
}
