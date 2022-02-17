package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/otter-im/profile/internal/app/model"
	"log"
	"net/http"
)

func ProfileDefaultHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-Otter-Login-User-Id")
	if id == "" {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "please login")
		return
	}
	profileByIdHandler(w, r, id)
}

func ProfileGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileByIdHandler(w, r, vars["id"])
}

func profileByIdHandler(w http.ResponseWriter, r *http.Request, userId string) {
	profiles := model.ProfileDTO()

	id, err := uuid.Parse(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid uuid: %v", userId)
		return
	}

	profile, err := profiles.SelectProfile(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "profile error for id: %v", userId)
		log.Println(err)
		return
	}

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "json marshal error for id: %v", userId)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
