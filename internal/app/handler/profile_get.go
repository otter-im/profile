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

func ProfileGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	profiles := model.ProfileDTO()

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid uuid: %v", vars["id"])
		return
	}

	profile, err := profiles.SelectProfile(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "profile error for id: %v", vars["id"])
		log.Println(err)
		return
	}

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "json marshal error for id: %v", vars["id"])
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
