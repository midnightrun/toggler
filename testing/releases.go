package testing

import (
	"fmt"

	"github.com/adamluzsi/frameless/fixtures"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	"github.com/toggler-io/toggler/domains/release"
)

const (
	ExampleReleaseManualPilotEnrollmentLetVar = `ExampleReleaseManualPilotEnrollment`

	ExampleReleaseRolloutLetVar  = `ReleaseRollout`
	ExampleReleaseFlagLetVar     = `ReleaseFlag`
	ExampleReleaseFlagNameLetVar = `ReleaseFlagName`

	ExamplePilotExternalIDLetVar = `PilotExternalID`
	ExamplePilotLetVar           = `ManualPilotEnrollment`
	ExamplePilotEnrollmentLetVar = `PilotEnrollment`
)

func init() {
	setups = append(setups, func(s *testcase.Spec) {
		s.Let(ExampleReleaseManualPilotEnrollmentLetVar, func(t *testcase.T) interface{} {
			mpe := Create(release.ManualPilotEnrollment{}).(*release.ManualPilotEnrollment)
			mpe.FlagID = ExampleReleaseFlag(t).ID
			mpe.DeploymentEnvironmentID = ExampleDeploymentEnvironment(t).ID
			mpe.ExternalID = ExampleExternalPilotID(t)
			require.Nil(t, ExampleStorage(t).Create(GetContext(t), mpe))
			t.Defer(ExampleStorage(t).DeleteByID, GetContext(t), release.ManualPilotEnrollment{}, mpe.ID)
			return mpe
		})

		s.Let(ExamplePilotExternalIDLetVar, func(t *testcase.T) interface{} {
			return RandomExternalPilotID()
		})

		s.Let(ExamplePilotEnrollmentLetVar, func(t *testcase.T) interface{} {
			return fixtures.Random.Bool()
		})

		s.Let(ExamplePilotLetVar, func(t *testcase.T) interface{} {
			// domains/release/specs/FlagFinderSpec.go:53:1: DEPRECATED, clean it up
			return &release.ManualPilotEnrollment{
				FlagID:                  ExampleReleaseFlag(t).ID,
				DeploymentEnvironmentID: ExampleDeploymentEnvironment(t).ID,
				ExternalID:              t.I(ExamplePilotExternalIDLetVar).(string),
				IsParticipating:         t.I(ExamplePilotEnrollmentLetVar).(bool),
			}
		})

		GivenWeHaveReleaseFlag(s, ExampleReleaseFlagLetVar)

		GivenWeHaveReleaseRollout(s,
			ExampleReleaseRolloutLetVar,
			ExampleReleaseFlagLetVar,
			ExampleDeploymentEnvironmentLetVar,
		)
	})
}

func ExampleReleaseManualPilotEnrollment(t *testcase.T) *release.ManualPilotEnrollment {
	return t.I(ExampleReleaseManualPilotEnrollmentLetVar).(*release.ManualPilotEnrollment)
}

func ExampleExternalPilotID(t *testcase.T) string {
	return t.I(ExamplePilotExternalIDLetVar).(string)
}

func ExampleReleaseFlagName(t *testcase.T) string {
	return t.I(ExampleReleaseFlagNameLetVar).(string)
}

func GetPilot(t *testcase.T, vn string) *release.ManualPilotEnrollment {
	return t.I(vn).(*release.ManualPilotEnrollment)
}

func ExamplePilot(t *testcase.T) *release.ManualPilotEnrollment {
	return GetPilot(t, ExamplePilotLetVar)
}

func FindStoredExampleReleaseFlagByName(t *testcase.T) *release.Flag {
	return FindStoredReleaseFlagByName(t, ExampleReleaseFlagName(t))
}

func FindStoredReleaseFlagByName(t *testcase.T, name string) *release.Flag {
	f, err := ExampleStorage(t).FindReleaseFlagByName(GetContext(t), name)
	require.Nil(t, err)
	require.NotNil(t, f)
	return f
}

func ExampleReleaseRollout(t *testcase.T) *release.Rollout {
	return GetReleaseRollout(t, ExampleReleaseRolloutLetVar)
}

func getReleaseRolloutPlanLetVar(vn string) string {
	return fmt.Sprintf(`%s.plan`, vn)
}

func GetReleaseRolloutPlan(t *testcase.T, rolloutLVN string) release.RolloutDefinition {
	return t.I(getReleaseRolloutPlanLetVar(rolloutLVN)).(release.RolloutDefinition)
}

func GetReleaseRollout(t *testcase.T, vn string) *release.Rollout {
	return t.I(vn).(*release.Rollout)
}

func GivenWeHaveReleaseRollout(s *testcase.Spec, vn, flagLVN, envLVN string) {
	s.Let(fmt.Sprintf(`%s.plan`, vn), func(t *testcase.T) interface{} {
		def := release.NewRolloutDecisionByPercentage()
		return &def
	})
	s.Let(vn, func(t *testcase.T) interface{} {
		rf := GetReleaseFlag(t, flagLVN)
		de := GetDeploymentEnvironment(t, envLVN)

		rollout := FixtureFactory{}.Create(release.Rollout{}).(*release.Rollout)
		rollout.FlagID = rf.ID
		rollout.DeploymentEnvironmentID = de.ID
		rollout.RolloutPlan = GetReleaseRolloutPlan(t, vn)
		require.Nil(t, rollout.RolloutPlan.Validate())

		// TODO: replace when rollout manager has function for this
		require.Nil(t, ExampleRolloutManager(t).Storage.Create(GetContext(t), rollout))
		t.Defer(ExampleRolloutManager(t).Storage.DeleteByID, GetContext(t), *rollout, rollout.ID)
		return rollout
	})
}

func GivenWeHaveReleaseFlag(s *testcase.Spec, vn string) {
	s.Let(vn, func(t *testcase.T) interface{} {
		rf := FixtureFactory{}.Create(release.Flag{}).(*release.Flag)
		require.Nil(t, ExampleRolloutManager(t).Storage.Create(GetContext(t), rf))
		t.Defer(ExampleRolloutManager(t).DeleteFeatureFlag, GetContext(t), rf.ID)
		return rf
	})
}

// FIXME replace caller release flag LVNs with rollout LVNs
func AndReleaseFlagPercentageIs(s *testcase.Spec, rolloutLVN string, percentage int) {
	s.Before(func(t *testcase.T) {
		byPercentage, ok := GetReleaseRolloutPlan(t, rolloutLVN).(*release.RolloutDecisionByPercentage)
		require.True(t, ok, `unexpected release rollout plan definition for AndReleaseFlagPercentageIs helper`)
		byPercentage.Percentage = percentage
	})
}

func GetReleaseFlag(t *testcase.T, lvn string) *release.Flag {
	ff := t.I(lvn)
	if ff == nil {
		return nil
	}
	return ff.(*release.Flag)
}

func ExampleReleaseFlag(t *testcase.T) *release.Flag {
	return GetReleaseFlag(t, ExampleReleaseFlagLetVar)
}

func ExampleRolloutManager(t *testcase.T) *release.RolloutManager {
	return release.NewRolloutManager(ExampleStorage(t))
}

func SpecPilotEnrolmentIs(t *testcase.T, enrollment bool) {
	if ExampleReleaseFlag(t).ID == `` {
		require.Nil(t, ExampleStorage(t).Create(GetContext(t), ExampleReleaseFlag(t)))
	}

	rm := release.NewRolloutManager(ExampleStorage(t))
	require.Nil(t, rm.SetPilotEnrollmentForFeature(GetContext(t),
		ExampleReleaseFlag(t).ID,
		ExampleDeploymentEnvironment(t).ID,
		ExampleExternalPilotID(t),
		enrollment))
}

func NoReleaseFlagPresentInTheStorage(s *testcase.Spec) {
	s.Before(func(t *testcase.T) {
		// TODO: replace with flag manager list+delete
		require.Nil(t, ExampleStorage(t).DeleteAll(GetContext(t), release.Flag{}))
	})
}

func AndExamplePilotManualParticipatingIsSetTo(s *testcase.Spec, isParticipating bool) {
	s.Before(func(t *testcase.T) {
		require.Nil(t, ExampleRolloutManager(t).SetPilotEnrollmentForFeature(
			GetContext(t),
			ExampleReleaseFlag(t).ID,
			ExampleDeploymentEnvironment(t).ID,
			ExampleExternalPilotID(t),
			isParticipating,
		))

		t.Defer(ExampleRolloutManager(t).UnsetPilotEnrollmentForFeature,
			GetContext(t),
			ExampleReleaseFlag(t).ID,
			ExampleDeploymentEnvironment(t).ID,
			ExampleExternalPilotID(t),
		)
	})
}
