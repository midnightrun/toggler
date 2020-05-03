package specs

import (
	"testing"

	"github.com/adamluzsi/frameless/resources"
	"github.com/adamluzsi/frameless/resources/specs"

	"github.com/toggler-io/toggler/domains/release"
)

type RolloutStorageSpec struct {
	Subject interface {
		release.RolloutFinder
		resources.Creator
		resources.Finder
		resources.Deleter
		resources.Updater
	}
	specs.FixtureFactory
}

func (spec RolloutStorageSpec) Test(t *testing.T) {
	t.Run(`Root`, func(t *testing.T) {
		specs.CommonSpec{
			EntityType:     release.RolloutDecisionByPercentage{},
			FixtureFactory: spec.FixtureFactory,
			Subject:        spec.Subject,
		}.Test(t)

		specs.CommonSpec{
			EntityType:     release.RolloutDecisionByAPI{},
			FixtureFactory: spec.FixtureFactory,
			Subject:        spec.Subject,
		}.Test(t)
	})
}

func (spec RolloutStorageSpec) Benchmark(b *testing.B) {
	b.Skip()
}

//FindReleaseRolloutByReleaseFlagAndDeploymentEnvironment
