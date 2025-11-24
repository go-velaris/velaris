package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

type listPolicy struct {
	engine *xorm.Engine
}

func NewListPolicy(engine *xorm.Engine) Policy {
	return &listPolicy{engine}
}

func userOwnsList(engine *xorm.Engine, ctxUser middleware.CtxUser, id int64) (bool, error) {
	exists, err := engine.
		Alias("l").
		Where("l.id = ?", id).
		Join("INNER", "boards b", "b.id = l.board_id").
		Where("b.user_id = ?", ctxUser.ID).
		Exist()
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (lp *listPolicy) CanView(ctxUser middleware.CtxUser, listId int64) (bool, error) {
	return userOwnsList(lp.engine, ctxUser, listId)
}

func (lp *listPolicy) CanCreate(ctxUser middleware.CtxUser, boardId int64) (bool, error) {
	return userOwnsBoard(lp.engine, ctxUser, boardId)
}

func (lp *listPolicy) CanUpdate(ctxUser middleware.CtxUser, listId int64) (bool, error) {
	return userOwnsList(lp.engine, ctxUser, listId)
}

func (lp *listPolicy) CanDelete(ctxUser middleware.CtxUser, listId int64) (bool, error) {
	return userOwnsList(lp.engine, ctxUser, listId)
}
