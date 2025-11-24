package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mithileshgupta12/velaris/internal/db/policy"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/middleware"
)

type StoreListRequest struct {
	Name string `json:"name"`
}

type ListHandler struct {
	listRepository repository.ListRepository
	boardPolicy    policy.Policy
	listPolicy     policy.Policy
}

func NewListHandler(
	listRepository repository.ListRepository,
	boardPolicy policy.Policy,
	listPolicy policy.Policy,
) *ListHandler {
	return &ListHandler{listRepository, boardPolicy, listPolicy}
}

func (lh *ListHandler) Index(w http.ResponseWriter, r *http.Request) {
	boardId, err := helper.ParseIntURLParam(r, "boardId")
	if err != nil || boardId < 1 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	canView, err := lh.boardPolicy.CanView(ctxUser, boardId)
	if err != nil {
		slog.Error("failed to check board view permission", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !canView {
		helper.ErrorJsonResponse(w, http.StatusNotFound, "board not found")
		return
	}

	lists, err := lh.listRepository.GetAllListsByBoardId(&repository.GetAllListsByBoardIdArgs{
		BoardId: boardId,
	})
	if err != nil {
		slog.Error("failed to get lists for board", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusOK, lists)
}

func (lh *ListHandler) Store(w http.ResponseWriter, r *http.Request) {
	var storeListRequest StoreListRequest

	if err := json.NewDecoder(r.Body).Decode(&storeListRequest); err != nil {
		slog.Error("failed to decode request", "err", err)
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	storeListRequest.Name = strings.TrimSpace(storeListRequest.Name)

	if storeListRequest.Name == "" {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name is a required field")
		return
	}

	if len(storeListRequest.Name) > 255 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "name must not be more than 255 characters long")
		return
	}

	boardId, err := helper.ParseIntURLParam(r, "boardId")
	if err != nil || boardId < 1 {
		helper.ErrorJsonResponse(w, http.StatusBadRequest, "invalid board id")
		return
	}

	ctxUser := r.Context().Value(middleware.CtxUserKey).(middleware.CtxUser)

	canCreate, err := lh.listPolicy.CanCreate(ctxUser, boardId)
	if err != nil {
		slog.Error("failed to check list create permission", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !canCreate {
		helper.ErrorJsonResponse(w, http.StatusNotFound, "board not found")
		return
	}

	list, err := lh.listRepository.CreateList(&repository.CreateListArgs{
		Name:    storeListRequest.Name,
		BoardId: boardId,
	})
	if err != nil {
		slog.Error("failed to create list", "err", err)
		helper.ErrorJsonResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	helper.JsonResponse(w, http.StatusCreated, list)
}

func (lh *ListHandler) Show(w http.ResponseWriter, r *http.Request) {
	//
}

func (lh *ListHandler) Update(w http.ResponseWriter, r *http.Request) {
	//
}

func (lh *ListHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	//
}
