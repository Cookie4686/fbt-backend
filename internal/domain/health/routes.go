package health

import (
	"fbt/backend/internal/dependency"
	"fbt/backend/internal/util"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthResponsePayload struct {
	Status string `json:"status"`
}

func Routes(d *dependency.Dependency, r *mux.Router) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		util.SendJson(w, &util.Response[HealthResponsePayload]{
			StatusCode: http.StatusOK,
			Payload:    &HealthResponsePayload{Status: "ok"},
		})
	})
}
