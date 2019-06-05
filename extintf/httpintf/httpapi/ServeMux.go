package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/adamluzsi/FeatureFlags/usecases"
	"log"
	"net/http"
)

func NewServeMux(uc *usecases.UseCases) *ServeMux {
	mux := &ServeMux{ServeMux: http.NewServeMux(), UseCases: uc,}
	featureAPI := buildFeatureAPI(mux)
	flagsAPI := buildFlagAPI(mux)

	mux.Handle(`/feature/`, http.StripPrefix(`/feature`, featureAPI))
	featureAPI.Handle(`/flag/`, http.StripPrefix(`/flag`, flagsAPI))

	return mux
}

func buildFeatureAPI(handlers *ServeMux) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(`/is-enabled.json`, http.HandlerFunc(handlers.IsFeatureEnabledFor))
	mux.Handle(`/is-globally-enabled.json`, http.HandlerFunc(handlers.IsFeatureGloballyEnabled))
	return mux
}

func buildFlagAPI(handlers *ServeMux) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(`/create.form`, http.HandlerFunc(handlers.CreateFeatureFlagFORM))
	mux.Handle(`/create.json`, http.HandlerFunc(handlers.CreateFeatureFlagJSON))
	mux.Handle(`/update.form`, http.HandlerFunc(handlers.UpdateFeatureFlagFORM))
	mux.Handle(`/update.json`, http.HandlerFunc(handlers.UpdateFeatureFlagJSON))
	mux.Handle(`/list.json`, http.HandlerFunc(handlers.ListFeatureFlags))
	mux.Handle(`/set-enrollment-manually.json`, http.HandlerFunc(handlers.SetPilotEnrollmentForFeature))
	return authMiddleware(handlers.UseCases, mux)
}

type ServeMux struct {
	*http.ServeMux
	*usecases.UseCases
}

func serveJSON(w http.ResponseWriter, status int, data interface{}) {
	buf := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set(`Content-Type`, `application/json`)
	w.WriteHeader(status)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Println(err)
	}
}

func handleError(w http.ResponseWriter, err error, errCode int) (errorWasHandled bool) {
	if err == usecases.ErrInvalidToken {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return true
	}

	if err != nil {
		http.Error(w, http.StatusText(errCode), errCode)
		return true
	}

	return false
}

func authMiddleware(uc *usecases.UseCases, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.URL.Query().Get(`token`)

		if token == `` {
			token = r.Header.Get(`X-Auth-Token`)
		}

		if token == `` {
			cookie, err := r.Cookie(`token`)
			if err != nil && err != http.ErrNoCookie {
				handleError(w, err, http.StatusInternalServerError)
				return
			}

			if err != http.ErrNoCookie && cookie != nil {
				token = cookie.Value
			}
		}

		pu, err := uc.ProtectedUsecases(token)

		if err == usecases.ErrInvalidToken {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if handleError(w, err, http.StatusInternalServerError) {
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), `ProtectedUsecases`, pu)))

	})
}
