package interactors_test

import (
	"github.com/adamluzsi/FeatureFlags/services/rollouts"
	"github.com/adamluzsi/frameless"
	"github.com/adamluzsi/frameless/iterators"
	"github.com/adamluzsi/frameless/reflects"
	"github.com/adamluzsi/frameless/resources/storages/memorystorage"
)

const (
	flagName           = `test-flag`
	PublicIDOfThePilot = `42`
)

func NewTestStorage() *TestStorage {
	return &TestStorage{
		Memory: memorystorage.NewMemory(),
	}
}

type TestStorage struct {
	*memorystorage.Memory
}

func (storage *TestStorage) FindPilotsByFeatureFlag(ff *rollouts.FeatureFlag) frameless.Iterator {
	table := storage.TableFor(rollouts.Pilot{})

	var pilots []*rollouts.Pilot

	for _, v := range table {
		pilot := v.(*rollouts.Pilot)

		if pilot.FeatureFlagID == ff.ID {
			pilots = append(pilots, pilot)
		}
	}

	return iterators.NewEmpty()
}

func (storage *TestStorage) FindPilotByFeatureFlagIDAndPublicPilotID(FeatureFlagID, ExternalPublicPilotID string) (*rollouts.Pilot, error) {
	table := storage.TableFor(rollouts.Pilot{})

	for _, v := range table {
		pilot := v.(*rollouts.Pilot)

		if pilot.FeatureFlagID == FeatureFlagID && pilot.ExternalPublicID == ExternalPublicPilotID {
			return pilot, nil
		}
	}

	return nil, nil
}

func (storage *TestStorage) FindByFlagName(name string, ptr *rollouts.FeatureFlag) (bool, error) {
	table := storage.TableFor(ptr)

	for _, v := range table {
		flag := v.(*rollouts.FeatureFlag)

		if flag.Name == name {
			return true, reflects.Link(flag, ptr)
		}
	}

	return false, nil
}