package httpapi

import (
	"net/http"

	"github.com/toggler-io/toggler/external/interface/httpintf/httputils"
	"github.com/toggler-io/toggler/usecases"
)

func (sm *Handler) SetPilotEnrollmentForFeature(w http.ResponseWriter, r *http.Request) {

	pu := r.Context().Value(`ProtectedUseCases`).(*usecases.ProtectedUseCases)

	pilot, err := httputils.ParseFlagPilotFromForm(r)

	if err != nil || pilot.FlagID == `` || pilot.ExternalID == `` {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = pu.SetPilotEnrollmentForFeature(r.Context(), pilot.FlagID, pilot.ExternalID, pilot.Enrolled)

	if handleError(w, err, http.StatusInternalServerError) {
		return
	}

	serveJSON(w,  map[string]interface{}{})

}