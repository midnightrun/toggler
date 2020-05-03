package specs

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/adamluzsi/frameless/iterators"
	"github.com/adamluzsi/frameless/reflects"
	"github.com/adamluzsi/frameless/resources"
	"github.com/adamluzsi/testcase"
	"github.com/google/uuid"

	. "github.com/toggler-io/toggler/testing"

	"github.com/toggler-io/toggler/domains/release"

	"github.com/adamluzsi/frameless"
	"github.com/adamluzsi/frameless/resources/specs"
	"github.com/stretchr/testify/require"
)

type pilotFinderSpec struct {
	Subject interface {
		release.PilotFinder
		resources.Creator
		resources.Updater
		resources.Finder
		resources.Deleter
	}
	FixtureFactory specs.FixtureFactory
}

func (spec pilotFinderSpec) Test(t *testing.T) {
	s := testcase.NewSpec(t)
	SetUp(s)

	s.LetValue(ExampleStorageLetVar, spec.Subject)

	s.Test(`ManualPilotEnrollment`, func(t *testcase.T) {
		specs.CommonSpec{
			EntityType:     release.ManualPilotEnrollment{},
			FixtureFactory: spec.ff(t),
			Subject:        spec.Subject,
		}.Test(t.T)
	})

	s.Describe(`pilotFinderSpec`, func(s *testcase.Spec) {
		s.Before(func(t *testcase.T) {
			require.Nil(t, spec.Subject.DeleteAll(spec.ctx(), release.ManualPilotEnrollment{}))
		})

		s.After(func(t *testcase.T) {
			require.Nil(t, spec.Subject.DeleteAll(spec.ctx(), release.ManualPilotEnrollment{}))
		})

		s.Describe(`FindReleasePilotsByReleaseFlag`, func(s *testcase.Spec) {
			subject := func(t *testcase.T) frameless.Iterator {
				pilotEntriesIter := spec.Subject.FindReleasePilotsByReleaseFlag(spec.ctx(), *ExampleReleaseFlag(t))
				t.Defer(pilotEntriesIter.Close)
				return pilotEntriesIter
			}

			thenNoPilotsFound := func(s *testcase.Spec) {
				s.Then(`no pilots found`, func(t *testcase.T) {
					iter := subject(t)
					require.NotNil(t, iter)
					require.False(t, iter.Next())
					require.Nil(t, iter.Err())
					require.Nil(t, iter.Close())
				})
			}

			s.When(`feature object is nil`, func(s *testcase.Spec) {
				s.LetValue(ExampleReleaseFlagLetVar, nil)

				thenNoPilotsFound(s)
			})

			s.When(`flag was never persisted before`, func(s *testcase.Spec) {
				s.Let(ExampleReleaseFlagLetVar, func(t *testcase.T) interface{} {
					return Create(release.Flag{})
				})

				thenNoPilotsFound(s)
			})

			s.When(`flag is persisted`, func(s *testcase.Spec) {
				thenNoPilotsFound(s)

				s.And(`there are manual pilot configs for the release flag`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						expectedPilots := t.I(`expectedPilots`).([]*release.ManualPilotEnrollment)

						for _, pilot := range expectedPilots {
							require.Nil(t, spec.Subject.Create(spec.ctx(), pilot))
						}
					})

					s.Let(`expectedPilots`, func(t *testcase.T) interface{} {
						var expectedPilots []*release.ManualPilotEnrollment
						for i := 0; i < 5; i++ {
							pilot := &release.ManualPilotEnrollment{
								FlagID:                  ExampleReleaseFlag(t).ID,
								DeploymentEnvironmentID: ExampleDeploymentEnvironment(t).ID,
								ExternalID:              strconv.Itoa(i),
							}

							expectedPilots = append(expectedPilots, pilot)
						}
						return expectedPilots
					})

					s.Then(`it will return all of them`, func(t *testcase.T) {
						iter := subject(t)
						defer iter.Close()
						require.NotNil(t, iter)
						var actualPilots []*release.ManualPilotEnrollment
						for iter.Next() {
							var actually release.ManualPilotEnrollment
							require.Nil(t, iter.Decode(&actually))
							actualPilots = append(actualPilots, &actually)
						}
						require.Nil(t, iter.Err())

						expectedPilots := t.I(`expectedPilots`).([]*release.ManualPilotEnrollment)
						require.True(t, len(expectedPilots) == len(actualPilots))
						require.ElementsMatch(t, expectedPilots, actualPilots)
					})
				})
			})
		})

		s.Describe(`FindReleasePilotByReleaseFlagAndDeploymentEnvironmentAndExternalID`, func(s *testcase.Spec) {
			subject := func(t *testcase.T) (*release.ManualPilotEnrollment, error) {
				return spec.Subject.FindReleasePilotByReleaseFlagAndDeploymentEnvironmentAndExternalID(
					spec.ctx(),
					ExampleReleaseFlag(t).ID,
					ExampleDeploymentEnvironment(t).ID,
					ExampleID(t),
				)
			}

			s.Before(func(t *testcase.T) {
				require.Nil(t, spec.Subject.DeleteAll(spec.ctx(), release.ManualPilotEnrollment{}))
			})

			ThenNoPilotsFound := func(s *testcase.Spec) {
				s.Then(`no pilots found`, func(t *testcase.T) {
					pilot, err := subject(t)
					require.Nil(t, err)
					require.Nil(t, pilot)
				})
			}

			s.When(`flag is not persisted`, func(s *testcase.Spec) {
				s.Let(ExampleReleaseFlagLetVar, func(t *testcase.T) interface{} {
					return Create(release.Flag{})
				})

				ThenNoPilotsFound(s)
			})

			s.When(`flag persisted already exists`, func(s *testcase.Spec) {
				s.Let(`featureFlagID`, func(t *testcase.T) interface{} {
					return ExampleReleaseFlag(t).ID
				})

				ThenNoPilotsFound(s)

				s.And(`the given there is a registered pilot for the feature`, func(s *testcase.Spec) {
					s.Around(func(t *testcase.T) func() {
						pilot := &release.ManualPilotEnrollment{
							FlagID:                  ExampleReleaseFlag(t).ID,
							DeploymentEnvironmentID: ExampleDeploymentEnvironment(t).ID,
							ExternalID:              ExampleID(t),
						}
						require.Nil(t, spec.Subject.Create(spec.ctx(), pilot))
						return func() { require.Nil(t, spec.Subject.DeleteByID(spec.ctx(), *pilot, pilot.ID)) }
					})

					s.Then(`then pilots will be retrieved`, func(t *testcase.T) {
						pilot, err := subject(t)
						require.Nil(t, err)
						require.NotNil(t, pilot)

						require.Equal(t, ExampleID(t), pilot.ExternalID)
						require.Equal(t, ExampleReleaseFlag(t).ID, pilot.FlagID)
						require.Equal(t, ExampleDeploymentEnvironment(t).ID, pilot.DeploymentEnvironmentID)
					})
				})
			})
		})

		s.Describe(`FindReleasePilotsByExternalID`, func(s *testcase.Spec) {
			subject := func(t *testcase.T) frameless.Iterator {
				pilotEntriesIter := spec.Subject.FindReleasePilotsByExternalID(spec.ctx(), ExampleExternalPilotID(t))
				t.Defer(pilotEntriesIter.Close)
				return pilotEntriesIter
			}

			s.Let(`PilotExternalID`, func(t *testcase.T) interface{} {
				return RandomExternalPilotID()
			})

			s.When(`there is no pilot records`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) { require.Nil(t, spec.Subject.DeleteAll(spec.ctx(), release.ManualPilotEnrollment{})) })

				s.Then(`it will return an empty result set`, func(t *testcase.T) {
					count, err := iterators.Count(subject(t))
					require.Nil(t, err)
					require.Equal(t, 0, count)
				})
			})

			s.When(`the given pilot id has no records`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					ctx := spec.ctx()
					extID := RandomExternalPilotID()

					var newUUID = func() string {
						uuidV4, err := uuid.NewRandom()
						require.Nil(t, err)
						return uuidV4.String()
					}

					require.Nil(t, spec.Subject.Create(ctx, &release.ManualPilotEnrollment{FlagID: newUUID(), ExternalID: extID, IsParticipating: true}))
					require.Nil(t, spec.Subject.Create(ctx, &release.ManualPilotEnrollment{FlagID: newUUID(), ExternalID: extID, IsParticipating: true}))
					require.Nil(t, spec.Subject.Create(ctx, &release.ManualPilotEnrollment{FlagID: newUUID(), ExternalID: extID, IsParticipating: false}))
				})

				s.Then(`it will return an empty result set`, func(t *testcase.T) {
					count, err := iterators.Count(subject(t))
					require.Nil(t, err)
					require.Equal(t, 0, count)
				})
			})

			s.When(`pilot ext id has multiple records`, func(s *testcase.Spec) {
				s.Let(`expected pilots`, func(t *testcase.T) interface{} {
					var pilots []release.ManualPilotEnrollment

					for i := 0; i < rand.Intn(5)+5; i++ {
						uuidV4, err := uuid.NewRandom()
						require.Nil(t, err)

						pilot := release.ManualPilotEnrollment{
							FlagID:          uuidV4.String(),
							ExternalID:      ExampleExternalPilotID(t),
							IsParticipating: rand.Intn(1) == 0,
						}

						require.Nil(t, spec.Subject.Create(spec.ctx(), &pilot))
						pilots = append(pilots, pilot)
					}

					return pilots
				})

				s.Before(func(t *testcase.T) { t.I(`expected pilots`) }) // eager load let value

				s.Then(`it will return all of them`, func(t *testcase.T) {
					var pilots []release.ManualPilotEnrollment
					require.Nil(t, iterators.Collect(subject(t), &pilots))
					require.ElementsMatch(t, t.I(`expected pilots`).([]release.ManualPilotEnrollment), pilots)
				})
			})
		})
	})
}

func (spec pilotFinderSpec) Benchmark(b *testing.B) {
	b.Run(`pilotFinderSpec`, func(b *testing.B) {
		b.Skip(`TODO`)
	})
}

func (spec pilotFinderSpec) ctx() context.Context {
	return spec.FixtureFactory.Context()
}

func (spec pilotFinderSpec) ff(t *testcase.T) specs.FixtureFactory {
	return &FixtureFactoryForPilots{
		FixtureFactory:             spec.FixtureFactory,
		GetFlagID:                  func() string { return ExampleReleaseFlag(t).ID },
		GetDeploymentEnvironmentID: func() string { return ExampleDeploymentEnvironment(t).ID },
	}
}

type FixtureFactoryForPilots struct {
	specs.FixtureFactory
	GetFlagID                  func() string
	GetDeploymentEnvironmentID func() string
}

func (ff *FixtureFactoryForPilots) Create(EntityType interface{}) interface{} {
	switch reflects.BaseValueOf(EntityType).Interface().(type) {
	case release.ManualPilotEnrollment:
		pilot := ff.FixtureFactory.Create(EntityType).(*release.ManualPilotEnrollment)
		pilot.FlagID = ff.GetFlagID()
		pilot.DeploymentEnvironmentID = ff.GetDeploymentEnvironmentID()
		return pilot

	default:
		return ff.FixtureFactory.Create(EntityType)
	}
}
