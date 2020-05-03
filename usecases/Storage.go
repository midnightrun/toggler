package usecases

import (
	"io"

	"github.com/adamluzsi/frameless/resources"

	"github.com/toggler-io/toggler/domains/deployment"
	"github.com/toggler-io/toggler/domains/release"
	"github.com/toggler-io/toggler/domains/security"
)

type Storage interface {
	resources.Creator
	resources.Finder
	resources.Updater
	resources.Deleter
	release.FlagFinder
	release.PilotFinder
	release.RolloutFinder
	security.TokenFinder
	deployment.EnvironmentFinder
	io.Closer
}
