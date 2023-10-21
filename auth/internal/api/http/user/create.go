package user

import (
	"encoding/json"
	"github.com/patyukin/banking-system/auth/internal/converter"
	desc "github.com/patyukin/banking-system/auth/pkg/user_v1"
	"io"
	"log"
	"net/http"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(r.Body)

	var req desc.CreateUserRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Невалидный тело запроса", http.StatusBadRequest)
		return
	}

	if req.GetPassword() != req.GetPasswordConfirm() {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	user := converter.ToUserFromDesc(&req)
	id, err := h.userService.Create(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	resJson, err := json.Marshal(&desc.CreateUserResponse{Id: id})
	w.Write(resJson)
}
