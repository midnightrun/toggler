// package usecases serves as a documentation purpose for the reader of the project.
// When the reader start examine the services, the reader can have a quick grasp on the situation,
// by simply listing the files in the usecases pkg.
package usecases

import (
	"github.com/adamluzsi/frameless"
	"github.com/toggler-io/toggler/domains/release"
	"github.com/toggler-io/toggler/domains/security"
)

func NewUseCases(s Storage) *UseCases {
	return &UseCases{
		RolloutManager: release.NewRolloutManager(s),
		Doorkeeper:     security.NewDoorkeeper(s),
		Issuer:         security.NewIssuer(s),

		protectedUsecases: &ProtectedUseCases{
			RolloutManager: release.NewRolloutManager(s),
			Doorkeeper:     security.NewDoorkeeper(s),
			Issuer:         security.NewIssuer(s),
		},
	}
}

type UseCases struct {
	*release.RolloutManager
	*security.Doorkeeper
	*security.Issuer

	// Deprecated
	protectedUsecases *ProtectedUseCases
}

const ErrInvalidToken frameless.Error = `invalid token error`
