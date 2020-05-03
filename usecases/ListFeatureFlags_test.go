package usecases_test

import (
	"context"
	"testing"

	"github.com/toggler-io/toggler/domains/release"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"

	. "github.com/toggler-io/toggler/testing"
)

func TestUseCases_ListFeatureFlags(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()
	SetUp(s)

	subject := func(t *testcase.T) ([]release.Flag, error) {
		return GetProtectedUsecases(t).ListFeatureFlags(context.TODO())
	}

	onSuccess := func(t *testcase.T) []release.Flag {
		ffs, err := subject(t)
		require.Nil(t, err)
		return ffs
	}

	s.When(`there are at least one flag in the system`, func(s *testcase.Spec) {
		s.Before(func(t *testcase.T) { EnsureFlag(t, `42`, 0) })

		s.Then(`we receive back the release flag list`, func(t *testcase.T) {
			ffs := onSuccess(t)
			require.Equal(t, 1, len(ffs))
			require.Equal(t, `42`, ffs[0].Name)
		})
	})

	s.When(`there is no flag in the system`, func(s *testcase.Spec) {
		s.Before(func(t *testcase.T) {
			require.Nil(t, ExampleStorage(t).DeleteAll(context.Background(), release.Flag{}))
		})

		s.Then(`we receive back empty list`, func(t *testcase.T) {
			ffs := onSuccess(t)
			require.Equal(t, 0, len(ffs))
		})
	})

}
