package api

import (
	"database/sql"
	db "firstprj/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD ERU"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSONP(http.StatusBadRequest, createResponseError(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSONP(http.StatusInternalServerError, createResponseError(err))
	}
	ctx.JSONP(http.StatusOK, account)
	return
}

type getAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(context *gin.Context) {
	var req getAccountRequest
	err := context.ShouldBindUri(&req)

	if err != nil {
		context.JSONP(http.StatusBadRequest, createResponseError(err))
		return
	}
	account, err := s.store.GetAccount(context, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSONP(http.StatusNotFound, createResponseError(err))
			return
		}
		context.JSONP(http.StatusInternalServerError, createResponseError(err))
		return
	}
	context.JSONP(http.StatusOK, account)
	return
}

type getListAccountRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=8"`
}

func (s *Server) GetListAccount(context *gin.Context) {
	var req getListAccountRequest
	err := context.ShouldBindQuery(&req)

	if err != nil {
		context.JSONP(http.StatusBadRequest, createResponseError(err))
		return
	}
	accounts, err := s.store.ListAccounts(context, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	})
	if err != nil {
		context.JSONP(http.StatusInternalServerError, createResponseError(err))
		return
	}
	context.JSONP(http.StatusOK, accounts)
	return
}

type deleteAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) DeleteAccount(context *gin.Context) {
	var req deleteAccountRequest
	err := context.ShouldBindUri(&req)

	if err != nil {
		context.JSONP(http.StatusBadRequest, createResponseError(err))
		return
	}
	err = s.store.DeleteAccount(context, req.Id)
	if err != nil {
		context.JSONP(http.StatusInternalServerError, createResponseError(err))
		return
	}
	createResponse(context, http.StatusOK, "", "successfully")
	return
}

type updateAccountRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

func (s *Server) UpdateAccount(context *gin.Context) {
	var req getAccountRequest
	err := context.ShouldBindUri(&req)
	if err != nil {
		context.JSONP(http.StatusBadRequest, createResponseError(err))
		return
	}
	var req2 updateAccountRequest
	err = context.BindJSON(&req2)
	if err != nil {
		context.JSONP(http.StatusBadRequest, createResponseError(err))
		return
	}
	account, err := s.store.UpdateAccountMoney(context, req.Id, req2.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSONP(http.StatusNotFound, createResponseError(err))
			return
		}
		context.JSONP(http.StatusInternalServerError, createResponseError(err))
		return
	}
	createResponse(context, http.StatusOK, account, "update successfully")
	return
}
