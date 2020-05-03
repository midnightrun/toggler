package release

import (
	"context"

	"github.com/adamluzsi/frameless"
	"github.com/adamluzsi/frameless/iterators"

	"github.com/toggler-io/toggler/domains/deployment"
)

func NewRolloutManager(s Storage) *RolloutManager {
	return &RolloutManager{Storage: s}
}

// RolloutManager provides you with feature flag configurability.
// The manager use storage in a write heavy behavior.
//
// SRP: release manager
type RolloutManager struct{ Storage Storage }

const CtxPilotIpAddr = `pilot-ip-addr`

// GetAllReleaseFlagStatesOfThePilot check the flag states for every requested release flag.
// If a flag doesn't exist, then it will provide a turned off state to it.
// This helps with development where the flag name already agreed and used in the clients
// but the entity not yet been created in the system.
// Also this help if a flag is not cleaned up from the clients, the worst thing will be a disabled feature,
// instead of a breaking client.
// This also makes it harder to figure out `private` release flags
func (manager *RolloutManager) GetAllReleaseFlagStatesOfThePilot(ctx context.Context, pilotExternalID string, env deployment.Environment, flagNames ...string) (map[string]bool, error) {
	states := make(map[string]bool)

	for _, flagName := range flagNames {
		states[flagName] = false
	}

	flagIndexByID := make(map[string]*Flag)

	var pilotsIndex = make(map[string]*ManualPilotEnrollment)
	pilotsByExternalID := manager.Storage.FindReleasePilotsByExternalID(ctx, pilotExternalID)
	pilotsByExternalIDFilteredByEnv := iterators.Filter(pilotsByExternalID, func(p ManualPilotEnrollment) bool {
		return p.DeploymentEnvironmentID == env.ID
	})
	if err := iterators.ForEach(pilotsByExternalIDFilteredByEnv, func(p ManualPilotEnrollment) error {
		pilotsIndex[p.FlagID] = &p
		return nil
	}); err != nil {
		return nil, err
	}

	if err := iterators.ForEach(manager.Storage.FindReleaseFlagsByName(ctx, flagNames...), func(f Flag) error {
		flagIndexByID[f.ID] = &f

		if p, ok := pilotsIndex[f.ID]; ok {
			states[f.Name] = p.IsParticipating
			return nil
		}

		enrollment, err := manager.checkEnrollment(ctx, env, f, pilotExternalID, pilotsIndex)
		if err != nil {
			return err
		}

		states[f.Name] = enrollment

		return nil
	}); err != nil {
		return nil, err
	}

	return states, nil
}

func (manager *RolloutManager) checkEnrollment(ctx context.Context, env deployment.Environment, flag Flag, pilotExternalID string, manualPilotEnrollmentIndex map[string]*ManualPilotEnrollment) (bool, error) {
	if p, ok := manualPilotEnrollmentIndex[flag.ID]; ok {
		return p.IsParticipating, nil
	}

	var rollout Rollout
	found, err := manager.Storage.FindReleaseRolloutByReleaseFlagAndDeploymentEnvironment(ctx, flag, env, &rollout)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	return rollout.RolloutPlan.IsParticipating(ctx, pilotExternalID)
}

func (manager *RolloutManager) CreateFeatureFlag(ctx context.Context, flag *Flag) error {
	if flag == nil {
		return ErrMissingFlag
	}

	if err := flag.Validate(); err != nil {
		return err
	}

	if flag.ID != `` {
		return ErrInvalidAction
	}
	ff, err := manager.Storage.FindReleaseFlagByName(ctx, flag.Name)

	if err != nil {
		return err
	}

	if ff != nil {
		//TODO: this should be handled in transaction!
		// as mvp solution, it is acceptable for now,
		// but spec must be moved to the storage specs as `name is uniq across entries`
		// or transaction through context with serialization level must be used for this action
		return ErrFlagAlreadyExist
	}

	return manager.Storage.Create(ctx, flag)
}

func (manager *RolloutManager) UpdateFeatureFlag(ctx context.Context, flag *Flag) error {
	if flag == nil {
		return ErrMissingFlag
	}

	if err := flag.Validate(); err != nil {
		return err
	}

	return manager.Storage.Update(ctx, flag)
}

// TODO convert this into a stream
func (manager *RolloutManager) ListFeatureFlags(ctx context.Context) ([]Flag, error) {
	iter := manager.Storage.FindAll(ctx, Flag{})
	ffs := make([]Flag, 0) // empty slice required for null object pattern enforcement
	err := iterators.Collect(iter, &ffs)
	return ffs, err
}

func (manager *RolloutManager) UnsetPilotEnrollmentForFeature(ctx context.Context, flagID string, envID string, pilotExternalID string) error {

	var ff Flag

	found, err := manager.Storage.FindByID(ctx, &ff, flagID)

	if err != nil {
		return err
	}

	if !found {
		return frameless.ErrNotFound
	}

	pilot, err := manager.Storage.FindReleasePilotByReleaseFlagAndDeploymentEnvironmentAndExternalID(ctx, ff.ID, envID, pilotExternalID)

	if err != nil {
		return err
	}

	if pilot == nil {
		return nil
	}

	return manager.Storage.DeleteByID(ctx, pilot, pilot.ID)

}

func (manager *RolloutManager) SetPilotEnrollmentForFeature(ctx context.Context, flagID string, envID string, externalPilotID string, isEnrolled bool) error {

	var ff Flag

	found, err := manager.Storage.FindByID(ctx, &ff, flagID)

	if err != nil {
		return err
	}

	if !found {
		return frameless.ErrNotFound
	}

	pilot, err := manager.Storage.FindReleasePilotByReleaseFlagAndDeploymentEnvironmentAndExternalID(ctx, ff.ID, envID, externalPilotID)

	if err != nil {
		return err
	}

	if pilot != nil {
		pilot.IsParticipating = isEnrolled
		return manager.Storage.Update(ctx, pilot)
	}

	return manager.Storage.Create(ctx, &ManualPilotEnrollment{
		FlagID:                  ff.ID,
		DeploymentEnvironmentID: envID,
		ExternalID:              externalPilotID,
		IsParticipating:         isEnrolled,
	})

}

//TODO: make operation atomic between flags and pilots
// TODO delete ip addr allows as well
// TODO: rename
func (manager *RolloutManager) DeleteFeatureFlag(ctx context.Context, id string) error {
	if id == `` {
		return frameless.ErrIDRequired
	}

	var ff Flag
	found, err := manager.Storage.FindByID(ctx, &ff, id)

	if err != nil {
		return err
	}
	if !found {
		return frameless.ErrNotFound
	}

	if err := iterators.ForEach(manager.Storage.FindReleasePilotsByReleaseFlag(ctx, ff), func(pilot ManualPilotEnrollment) error {
		return manager.Storage.DeleteByID(ctx, pilot, pilot.ID)
	}); err != nil {
		return err
	}

	return manager.Storage.DeleteByID(ctx, Flag{}, id)
}
