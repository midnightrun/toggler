package httputils

import (
	"net/http"
	"strings"

	"github.com/toggler-io/toggler/domains/release"
	"github.com/pkg/errors"
)

func ParseFlagPilotFromForm(r *http.Request) (*release.ManualPilotEnrollment, error) {

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	var pilot release.ManualPilotEnrollment
	pilot.ID = r.FormValue(`pilot.id`)
	pilot.FlagID = r.FormValue(`pilot.flagID`)
	pilot.ExternalID = r.FormValue(`pilot.extID`)

	switch strings.ToLower(r.FormValue(`pilot.enrolled`)) {
	case `true`, `on`:
		pilot.IsParticipating = true
	case `false`, ``:
		pilot.IsParticipating = false
	default:
		return nil, errors.New(`unrecognised value for "pilot.enrolled" value`)
	}

	return &pilot, nil

}
