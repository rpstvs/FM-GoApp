package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rpstvs/fm-goapp/internal/store"
	"github.com/rpstvs/fm-goapp/internal/tokens"
	"github.com/rpstvs/fm-goapp/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct{
	Username string
	Password string

}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request){
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)~

	if err!= nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{})
		return
	}

	user , err := h.userStore.GetUserByUsername(req.Username)

	if err != nil || user == nil {
		//handle error
		return
	}

	passwordsMatch , err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		return 
	}

	if !passwordsMatch{
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
