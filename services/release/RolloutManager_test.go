package release_test

import (
	"context"
	"math/rand"
	"net/url"
	"testing"
	"time"

	"github.com/adamluzsi/testcase"
	"github.com/toggler-io/toggler/services/release"
	. "github.com/toggler-io/toggler/testing"

	"github.com/adamluzsi/frameless/iterators"
	"github.com/stretchr/testify/require"
)

func TestRolloutManager(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()
	SetupSpecCommonVariables(s)

	s.Let(`GeneratedRandomSeed`, func(t *testcase.T) interface{} {
		return time.Now().Unix()
	})

	s.Let(`RolloutManager`, func(t *testcase.T) interface{} {
		return &release.RolloutManager{
			Storage: GetStorage(t),

			RandSeedGenerator: func() int64 {
				return GetGeneratedRandomSeed(t)
			},
		}
	})

	s.Before(func(t *testcase.T) {
		require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Flag{}))
		require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Pilot{}))
	})

	SpecRolloutManagerCreateFeatureFlag(s)
	SpecRolloutManagerDeleteFeatureFlag(s)
	SpecRolloutManagerUpdateFeatureFlag(s)
	SpecRolloutManagerListFeatureFlags(s)
	SpecRolloutManagerSetPilotEnrollmentForFeature(s)
	SpecRolloutManagerUnsetPilotEnrollmentForFeature(s)
}

func SpecRolloutManagerDeleteFeatureFlag(s *testcase.Spec) {
	s.Describe(`DeleteFeatureFlag`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) error {
			return manager(t).DeleteFeatureFlag(context.TODO(), t.I(`flag ID`).(string))
		}

		s.Let(`FeatureFlagName`, func(t *testcase.T) interface{} { return ExampleName() })
		s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return nil })
		s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return rand.Intn(101) })
		s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(42) })
		s.Let(`ReleaseFlag`, func(t *testcase.T) interface{} {
			ff := &release.Flag{Name: t.I(`ReleaseFlagName`).(string)}
			ff.Rollout.RandSeed = t.I(`RolloutSeedSalt`).(int64)
			ff.Rollout.Strategy.Percentage = t.I(`RolloutPercentage`).(int)
			ff.Rollout.Strategy.DecisionLogicAPI = GetRolloutApiURL(t)
			return ff
		})

		s.When(`release flag id is empty`, func(s *testcase.Spec) {
			s.Let(`flag ID`, func(t *testcase.T) interface{} { return `` })

			s.Then(`it will return error about it`, func(t *testcase.T) {
				require.Error(t, subject(t))
			})
		})

		s.When(`it had been persisted previously`, func(s *testcase.Spec) {
			s.Let(`flag ID`, func(t *testcase.T) interface{} {
				return GetReleaseFlag(t).ID
			})

			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Save(context.TODO(), GetReleaseFlag(t)))
				require.NotEmpty(t, GetReleaseFlag(t).ID)
			})

			s.Then(`flag will be deleted`, func(t *testcase.T) {
				id := GetReleaseFlag(t).ID
				require.NotEmpty(t, id)
				require.Nil(t, subject(t))
			})

			s.And(`there are pilots manually set for the feature`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, manager(t).SetPilotEnrollmentForFeature(CTX(t), t.I(`flag ID`).(string), ExampleExternalPilotID(), true))
					require.Nil(t, manager(t).SetPilotEnrollmentForFeature(CTX(t), t.I(`flag ID`).(string), ExampleExternalPilotID(), false))
				})

				s.Then(`it will remove the pilots as well for the feature`, func(t *testcase.T) {
					require.Nil(t, subject(t))
					pilotCount, err := iterators.Count(GetStorage(t).FindAll(CTX(t), release.Pilot{}))
					require.Nil(t, err)
					require.Equal(t, 0, pilotCount)
				})

				s.And(`if other flags have pilots enrolled as well`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						othFlag := *GetReleaseFlag(t)
						othFlag.Name = `oth flag`
						othFlag.ID = ``

						require.Nil(t, GetStorage(t).Save(CTX(t), &othFlag))
						require.Nil(t, manager(t).SetPilotEnrollmentForFeature(CTX(t), othFlag.ID, ExampleExternalPilotID(), true))
						require.Nil(t, manager(t).SetPilotEnrollmentForFeature(CTX(t), othFlag.ID, ExampleExternalPilotID(), false))
					})

					s.Then(`they will be unaffected by the subject flag removal`, func(t *testcase.T) {
						require.Nil(t, subject(t))
						pilotCount, err := iterators.Count(GetStorage(t).FindAll(CTX(t), release.Pilot{}))

						require.Nil(t, err)
						require.Equal(t, 2, pilotCount)
					})
				})
			})
		})

	})
}

func SpecRolloutManagerCreateFeatureFlag(s *testcase.Spec) {
	s.Describe(`CreateFeatureFlag`, func(s *testcase.Spec) {
		subjectWithArgs := func(t *testcase.T, f *release.Flag) error {
			return manager(t).CreateFeatureFlag(context.TODO(), f)
		}

		subject := func(t *testcase.T) error {
			return subjectWithArgs(t, GetReleaseFlag(t))
		}

		s.Let(`FeatureFlagName`, func(t *testcase.T) interface{} { return ExampleName() })
		s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return nil })
		s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return rand.Intn(101) })
		s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(42) })
		s.Let(`ReleaseFlag`, func(t *testcase.T) interface{} {
			ff := &release.Flag{Name: t.I(`ReleaseFlagName`).(string)}
			ff.Rollout.RandSeed = t.I(`RolloutSeedSalt`).(int64)
			ff.Rollout.Strategy.Percentage = t.I(`RolloutPercentage`).(int)
			ff.Rollout.Strategy.DecisionLogicAPI = GetRolloutApiURL(t)
			return ff
		})

		s.Then(`on valid input the values persisted`, func(t *testcase.T) {
			require.Nil(t, subject(t))
			require.NotNil(t, FindStoredReleaseFlagByName(t))
			require.Equal(t, GetReleaseFlag(t), FindStoredReleaseFlagByName(t))
		})

		s.When(`name is empty`, func(s *testcase.Spec) {
			s.Let(`ReleaseFlagName`, func(t *testcase.T) interface{} { return "" })

			s.Then(`it will fail with invalid feature name`, func(t *testcase.T) {
				require.Equal(t, release.ErrNameIsEmpty, subject(t))
			})
		})

		s.When(`url`, func(s *testcase.Spec) {
			s.Context(`is not a valid request url`, func(s *testcase.Spec) {
				s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return `http//example.com` })

				s.Then(`it will fail with invalid url`, func(t *testcase.T) {
					require.Equal(t, release.ErrInvalidRequestURL, subject(t))
				})
			})

			s.Context(`is not defined or nil`, func(s *testcase.Spec) {
				s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return nil })

				s.Then(`it will be saved and will represent that no custom domain decision url used`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Nil(t, FindStoredReleaseFlagByName(t).Rollout.Strategy.DecisionLogicAPI)
				})
			})

			s.Context(`is a valid request URL`, func(s *testcase.Spec) {
				s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return `https://example.com` })

				s.Then(`it will persist for future usage`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, `https://example.com`, FindStoredReleaseFlagByName(t).Rollout.Strategy.DecisionLogicAPI.String())
				})
			})
		})

		s.When(`percentage`, func(s *testcase.Spec) {
			s.Context(`less than 0`, func(s *testcase.Spec) {
				s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return -1 + (rand.Intn(1024) * -1) })

				s.Then(`it will report error regarding the percentage`, func(t *testcase.T) {
					require.Equal(t, release.ErrInvalidPercentage, subject(t))
				})
			})

			s.Context(`greater than 100`, func(s *testcase.Spec) {
				s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return 101 + rand.Intn(1024) })

				s.Then(`it will report error regarding the percentage`, func(t *testcase.T) {
					require.Equal(t, release.ErrInvalidPercentage, subject(t))
				})
			})

			s.Context(`is a number between 0 and 100`, func(s *testcase.Spec) {
				s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return rand.Intn(101) })

				s.Then(`it will persist the received percentage`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, t.I(`RolloutPercentage`).(int), FindStoredReleaseFlagByName(t).Rollout.Strategy.Percentage)
				})
			})
		})

		s.When(`pseudo random seed salt`, func(s *testcase.Spec) {
			s.Context(`is 0`, func(s *testcase.Spec) {
				s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(0) })

				s.Then(`random seed generator used for setting seed value`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, GetGeneratedRandomSeed(t), FindStoredReleaseFlagByName(t).Rollout.RandSeed)
				})
			})

			s.Context(`something else`, func(s *testcase.Spec) {
				s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(12) })

				s.Then(`it will persist the value`, func(t *testcase.T) {
					require.Nil(t, subject(t))

					require.Equal(t, int64(12), FindStoredReleaseFlagByName(t).Rollout.RandSeed)
				})
			})
		})

		s.When(`release flag`, func(s *testcase.Spec) {
			s.Context(`is nil`, func(s *testcase.Spec) {
				s.Let(`ReleaseFlag`, func(t *testcase.T) interface{} { return nil })

				s.Then(`it will return error about it`, func(t *testcase.T) {
					require.Error(t, subject(t))
				})
			})

			s.Context(`was not stored until now`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Flag{}))
				})

				s.Then(`it will be persisted`, func(t *testcase.T) {
					require.Nil(t, subject(t))
					require.NotNil(t, FindStoredReleaseFlagByName(t))
					require.Equal(t, GetReleaseFlag(t), FindStoredReleaseFlagByName(t))
				})
			})

			s.Context(`had been persisted previously`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Save(context.TODO(), GetReleaseFlag(t)))
					require.NotEmpty(t, GetReleaseFlag(t).ID)
				})

				s.When(`the id is not referring to the existing one`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						GetReleaseFlag(t).ID = ``
					})

					s.Then(`it will report release flag already exists error`, func(t *testcase.T) {
						require.Equal(t, release.ErrFlagAlreadyExist, subject(t))
					})
				})

				s.When(`the ID is set and pointing to an existing flag`, func(s *testcase.Spec) {
					s.Before(func(t *testcase.T) {
						require.NotEmpty(t, GetReleaseFlag(t).ID)
						var ff release.Flag
						found, err := GetStorage(t).FindByID(context.Background(), &ff, GetReleaseFlag(t).ID)
						require.Nil(t, err)
						require.True(t, found)
						require.Equal(t, GetReleaseFlag(t), &ff)
					})

					s.Then(`it will report invalid action error`, func(t *testcase.T) {
						require.Equal(t, release.ErrInvalidAction, subject(t))
					})
				})
			})
		})
	})
}

func SpecRolloutManagerUpdateFeatureFlag(s *testcase.Spec) {
	s.Describe(`UpdateFeatureFlag`, func(s *testcase.Spec) {
		subjectWithArgs := func(t *testcase.T, f *release.Flag) error {
			return manager(t).UpdateFeatureFlag(context.TODO(), f)
		}

		subject := func(t *testcase.T) error {
			return subjectWithArgs(t, GetReleaseFlag(t))
		}

		s.Let(`FeatureFlagName`, func(t *testcase.T) interface{} { return ExampleName() })
		s.Let(`RolloutApiURL`, func(t *testcase.T) interface{} { return nil })
		s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return rand.Intn(101) })
		s.Let(`RolloutSeedSalt`, func(t *testcase.T) interface{} { return int64(42) })
		s.Let(`ReleaseFlag`, func(t *testcase.T) interface{} {
			ff := &release.Flag{Name: t.I(`ReleaseFlagName`).(string)}
			ff.Rollout.RandSeed = t.I(`RolloutSeedSalt`).(int64)
			ff.Rollout.Strategy.Percentage = t.I(`RolloutPercentage`).(int)
			ff.Rollout.Strategy.DecisionLogicAPI = GetRolloutApiURL(t)
			return ff
		})

		s.When(`input is invalid for the release flag Verify low level domain requirement`, func(s *testcase.Spec) {
			s.Let(`RolloutPercentage`, func(t *testcase.T) interface{} { return 128 })

			s.Then(`it will report error`, func(t *testcase.T) {
				require.Error(t, subject(t))
			})
		})

		s.When(`release flag`, func(s *testcase.Spec) {
			s.Context(`is nil`, func(s *testcase.Spec) {
				s.Let(`ReleaseFlag`, func(t *testcase.T) interface{} { return nil })

				s.Then(`it will return error about it`, func(t *testcase.T) {
					require.Error(t, subject(t))
				})
			})

			s.Context(`was not stored until now`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Flag{}))
				})

				s.Then(`it will report error about the missing ID`, func(t *testcase.T) {
					require.Error(t, subject(t))
				})
			})

			s.Context(`had been persisted previously`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Save(context.TODO(), GetReleaseFlag(t)))
					require.NotEmpty(t, GetReleaseFlag(t).ID)
				})

				s.Then(`latest values are persisted in the storage`, func(t *testcase.T) {
					flag := *GetReleaseFlag(t) // pass by value copy
					newName := flag.Name + ` v2`
					flag.Name = newName
					flag.Rollout.Strategy.Percentage = 42
					u, err := url.Parse(`https://example.com`)
					require.Nil(t, err)
					flag.Rollout.Strategy.DecisionLogicAPI = u
					require.Nil(t, subjectWithArgs(t, &flag))

					var storedFlag release.Flag
					found, err := GetStorage(t).FindByID(context.Background(), &storedFlag, GetReleaseFlag(t).ID)
					require.Nil(t, err)
					require.True(t, found)
					require.Equal(t, u, storedFlag.Rollout.Strategy.DecisionLogicAPI)
					require.Equal(t, 42, storedFlag.Rollout.Strategy.Percentage)
					require.Equal(t, newName, storedFlag.Name)
				})
			})
		})

	})
}

func SpecRolloutManagerListFeatureFlags(s *testcase.Spec) {
	s.Describe(`ListFeatureFlags`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]*release.Flag, error) {
			return manager(t).ListFeatureFlags(context.TODO())
		}

		onSuccess := func(t *testcase.T) []*release.Flag {
			ffs, err := subject(t)
			require.Nil(t, err)
			return ffs
		}

		s.When(`features are in the system`, func(s *testcase.Spec) {

			s.Before(func(t *testcase.T) {
				EnsureFlag(t, `a`, 0)
				EnsureFlag(t, `b`, 0)
				EnsureFlag(t, `c`, 0)
			})

			s.Then(`release flags are returned`, func(t *testcase.T) {
				flags := onSuccess(t)

				expectedFlagNames := []string{`a`, `b`, `c`}

				for _, ff := range flags {
					require.Contains(t, expectedFlagNames, ff.Name)
				}
			})

		})

		s.When(`no feature present in the system`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Flag{}))
			})

			s.Then(`release flags are returned`, func(t *testcase.T) {
				flags := onSuccess(t)

				require.Equal(t, []*release.Flag{}, flags)
			})
		})

	})
}

func SpecRolloutManagerSetPilotEnrollmentForFeature(s *testcase.Spec) {
	s.Describe(`SetPilotEnrollmentForFeature`, func(s *testcase.Spec) {

		GetNewEnrollment := func(t *testcase.T) bool {
			return t.I(`NewEnrollment`).(bool)
		}

		subject := func(t *testcase.T) error {
			return manager(t).SetPilotEnrollmentForFeature(context.TODO(), t.I(`FlagID`).(string), GetExternalPilotID(t), GetNewEnrollment(t))
		}

		s.Let(`FlagID`, func(t *testcase.T) interface{} {
			return GetReleaseFlag(t).ID
		})

		s.Let(`NewEnrollment`, func(t *testcase.T) interface{} {
			return rand.Intn(2) == 0
		})

		findFlag := func(t *testcase.T) *release.Flag {
			iter := GetStorage(t).FindAll(context.Background(), &release.Flag{})
			require.NotNil(t, iter)
			require.True(t, iter.Next())
			var ff release.Flag
			require.Nil(t, iter.Decode(&ff))
			require.False(t, iter.Next())
			require.Nil(t, iter.Err())
			return &ff
		}

		s.When(`no release flag is seen ever before`, func(s *testcase.Spec) {
			s.Let(`FlagID`, func(t *testcase.T) interface{} { return `` })
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Flag{}))
			})

			s.Then(`error returned`, func(t *testcase.T) {
				require.Error(t, subject(t))
			})
		})

		s.When(`release flag already configured`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Save(context.TODO(), GetReleaseFlag(t)))
			})

			s.Then(`flag will not be recreated`, func(t *testcase.T) {
				require.Nil(t, subject(t))

				count, err := iterators.Count(GetStorage(t).FindAll(context.Background(), release.Flag{}))
				require.Nil(t, err)
				require.Equal(t, 1, count)

				flag := findFlag(t)
				require.Equal(t, GetReleaseFlag(t), flag)
			})

			s.And(`pilot already exists`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Save(context.TODO(), GetPilot(t)))
				})

				s.And(`and pilot is has the opposite enrollment status`, func(s *testcase.Spec) {
					s.Let(`PilotEnrollment`, func(t *testcase.T) interface{} {
						return !GetNewEnrollment(t)
					})

					s.Then(`the original pilot is updated to the new enrollment status`, func(t *testcase.T) {
						require.Nil(t, subject(t))
						flag := findFlag(t)

						pilot, err := GetStorage(t).FindReleaseFlagPilotByPilotExternalID(context.Background(), flag.ID, GetExternalPilotID(t))
						require.Nil(t, err)

						require.NotNil(t, pilot)
						require.Equal(t, GetNewEnrollment(t), pilot.Enrolled)
						require.Equal(t, GetExternalPilotID(t), pilot.ExternalID)
						require.Equal(t, GetPilot(t), pilot)

						count, err := iterators.Count(GetStorage(t).FindAll(context.Background(), release.Pilot{}))
						require.Nil(t, err)
						require.Equal(t, 1, count)
					})
				})

				s.And(`pilot already has the same enrollment status`, func(s *testcase.Spec) {
					s.Let(`PilotEnrollment`, func(t *testcase.T) interface{} {
						return GetNewEnrollment(t)
					})

					s.Then(`pilot remain the same`, func(t *testcase.T) {

						require.Nil(t, subject(t))
						ff := findFlag(t)

						pilot, err := GetStorage(t).FindReleaseFlagPilotByPilotExternalID(context.Background(), ff.ID, GetExternalPilotID(t))
						require.Nil(t, err)

						require.NotNil(t, pilot)
						require.Equal(t, GetNewEnrollment(t), pilot.Enrolled)
						require.Equal(t, GetExternalPilotID(t), pilot.ExternalID)

						count, err := iterators.Count(GetStorage(t).FindAll(context.Background(), release.Pilot{}))
						require.Nil(t, err)
						require.Equal(t, 1, count)

					})
				})
			})
		})
	})
}

func SpecRolloutManagerUnsetPilotEnrollmentForFeature(s *testcase.Spec) {
	s.Describe(`UnsetPilotEnrollmentForFeature`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) error {
			return manager(t).UnsetPilotEnrollmentForFeature(context.TODO(), t.I(`FlagID`).(string), GetExternalPilotID(t))
		}

		s.Let(`FlagID`, func(t *testcase.T) interface{} {
			return GetReleaseFlag(t).ID
		})

		findFlag := func(t *testcase.T) *release.Flag {
			iter := GetStorage(t).FindAll(context.Background(), &release.Flag{})
			require.NotNil(t, iter)
			require.True(t, iter.Next())
			var ff release.Flag
			require.Nil(t, iter.Decode(&ff))
			require.False(t, iter.Next())
			require.Nil(t, iter.Err())
			return &ff
		}

		s.When(`no release flag is seen ever before`, func(s *testcase.Spec) {
			s.Let(`FlagID`, func(t *testcase.T) interface{} { return `` })
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Truncate(context.Background(), release.Flag{}))
			})

			s.Then(`error returned`, func(t *testcase.T) {
				require.Error(t, subject(t))
			})
		})

		s.When(`release flag already configured`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).Save(context.TODO(), GetReleaseFlag(t)))
			})

			s.Then(`flag will not be recreated`, func(t *testcase.T) {
				require.Nil(t, subject(t))

				count, err := iterators.Count(GetStorage(t).FindAll(CTX(t), release.Flag{}))
				require.Nil(t, err)
				require.Equal(t, 1, count)

				flag := findFlag(t)
				require.Equal(t, GetReleaseFlag(t), flag)
			})

			s.And(`pilot not exist for the flag`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Truncate(CTX(t), release.Pilot{}))
				})

				s.Then(`it will return without any error`, func(t *testcase.T) {
					require.Nil(t, subject(t))
				})
			})

			s.And(`pilot already exists`, func(s *testcase.Spec) {
				s.Before(func(t *testcase.T) {
					require.Nil(t, GetStorage(t).Save(CTX(t), GetPilot(t)))
				})

				s.Then(`pilot manual enrollment will be removed`, func(t *testcase.T) {
					require.Nil(t, subject(t))
					flag := findFlag(t)

					pilot, err := GetStorage(t).FindReleaseFlagPilotByPilotExternalID(context.Background(), flag.ID, GetExternalPilotID(t))
					require.Nil(t, err)
					require.Nil(t, pilot)
				})
			})
		})
	})
}

func GetGeneratedRandomSeed(t *testcase.T) int64 {
	return t.I(`GeneratedRandomSeed`).(int64)
}

func manager(t *testcase.T) *release.RolloutManager {
	return t.I(`RolloutManager`).(*release.RolloutManager)
}
