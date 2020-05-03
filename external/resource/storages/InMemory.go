package storages

import (
	"context"

	"github.com/adamluzsi/frameless"
	"github.com/adamluzsi/frameless/iterators"
	"github.com/adamluzsi/frameless/resources/storages/memorystorage"

	"github.com/toggler-io/toggler/domains/deployment"
	"github.com/toggler-io/toggler/domains/release"
	"github.com/toggler-io/toggler/domains/security"
)

func NewInMemory() *InMemory {
	return &InMemory{Memory: memorystorage.NewMemory()}
}

type InMemory struct{ *memorystorage.Memory }

func (s *InMemory) FindReleasePilotsByExternalID(ctx context.Context, pilotExtID string) release.PilotEntries {
	var pilots []release.ManualPilotEnrollment

	for _, e := range s.TableFor(release.ManualPilotEnrollment{}) {
		p := e.(*release.ManualPilotEnrollment)

		if p.ExternalID == pilotExtID {
			pilots = append(pilots, *p)
		}
	}

	return iterators.NewSlice(pilots)
}

func (s *InMemory) FindReleaseFlagsByName(ctx context.Context, names ...string) frameless.Iterator {
	var flags []release.Flag

	nameIndex := make(map[string]struct{})

	for _, name := range names {
		nameIndex[name] = struct{}{}
	}

	for _, e := range s.TableFor(release.Flag{}) {
		flag := e.(*release.Flag)

		if _, ok := nameIndex[flag.Name]; ok {
			flags = append(flags, *flag)
		}
	}

	return iterators.NewSlice(flags)
}

func (s *InMemory) FindReleasePilotsByReleaseFlag(ctx context.Context, flag release.Flag) release.PilotEntries {
	table := s.TableFor(release.ManualPilotEnrollment{})

	var pilots []release.ManualPilotEnrollment

	for _, v := range table {
		pilot := v.(*release.ManualPilotEnrollment)

		if pilot.FlagID == flag.ID {
			pilots = append(pilots, *pilot)
		}
	}

	return iterators.NewSlice(pilots)
}

func (s *InMemory) FindReleasePilotByReleaseFlagAndDeploymentEnvironmentAndExternalID(ctx context.Context, flagID, envID, pilotExtID string) (*release.ManualPilotEnrollment, error) {
	table := s.TableFor(release.ManualPilotEnrollment{})

	for _, v := range table {
		pilot := v.(*release.ManualPilotEnrollment)

		if pilot.FlagID == flagID && pilot.DeploymentEnvironmentID == envID && pilot.ExternalID == pilotExtID {
			p := *pilot
			return &p, nil
		}
	}

	return nil, nil
}

func (s *InMemory) FindDeploymentEnvironmentByAlias(ctx context.Context, idOrName string, env *deployment.Environment) (bool, error) {
	for _, v := range s.TableFor(deployment.Environment{}) {
		record := v.(*deployment.Environment)

		if record.ID == idOrName || record.Name == idOrName {
			*env = *record
			return true, nil
		}
	}
	return false, nil
}

func (s *InMemory) FindReleaseFlagByName(ctx context.Context, name string) (*release.Flag, error) {
	for _, v := range s.TableFor(release.Flag{}) {
		flagRecord := v.(*release.Flag)

		if flagRecord.Name == name {
			f := *flagRecord
			return &f, nil
		}
	}
	return nil, nil
}

func (s *InMemory) FindTokenBySHA512Hex(ctx context.Context, t string) (*security.Token, error) {
	table := s.TableFor(security.Token{})

	for _, token := range table {
		token := token.(*security.Token)

		if token.SHA512 == t {
			t := *token
			return &t, nil
		}
	}

	return nil, nil
}

func (s *InMemory) FindReleaseRolloutByReleaseFlagAndDeploymentEnvironment(ctx context.Context, flag release.Flag, env deployment.Environment, ptr *release.Rollout) (bool, error) {
	for _, rollout := range s.TableFor(*ptr) {
		rollout := rollout.(*release.Rollout)

		if rollout.FlagID == flag.ID && rollout.DeploymentEnvironmentID == env.ID {
			*ptr = *rollout
			return true, nil
		}
	}
	return false, nil
}
