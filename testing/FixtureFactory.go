package testing

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/adamluzsi/frameless/fixtures"
	"github.com/adamluzsi/frameless/reflects"
	"github.com/google/uuid"

	"github.com/toggler-io/toggler/domains/release"
	"github.com/toggler-io/toggler/domains/security"
)

func NewFixtureFactory() *FixtureFactory {
	return &FixtureFactory{}
}

type FixtureFactory struct {
	fixtures.GenericFixtureFactory
}

// this ensures that the randoms have better variety between test runs with -count n
var rnd = rand.New(rand.NewSource(time.Now().Unix()))

func (ff FixtureFactory) Create(EntityType interface{}) interface{} {
	switch reflects.BaseValueOf(EntityType).Interface().(type) {
	case release.Flag:
		flag := ff.GenericFixtureFactory.Create(EntityType).(*release.Flag)
		flag.Name = fmt.Sprintf(`%s - %s`, flag.Name, uuid.New().String())

		flag.Rollout.Strategy.DecisionLogicAPI = nil

		if rnd.Intn(2) == 0 {
			u, err := url.ParseRequestURI(fmt.Sprintf(`http://google.com/%s`, url.PathEscape(fixtures.Random.String())))

			if err != nil {
				panic(err)
			}

			flag.Rollout.Strategy.DecisionLogicAPI = u
		}

		flag.Rollout.Strategy.Percentage = fixtures.Random.IntBetween(0, 101)

		return flag

	case release.Pilot:
		pilot := ff.GenericFixtureFactory.Create(EntityType).(*release.Pilot)
		pilot.ExternalID = uuid.New().String()
		return pilot

	case security.Token:
		t := ff.GenericFixtureFactory.Create(EntityType).(*security.Token)
		t.SHA512 = uuid.New().String()
		return t

	default:
		return ff.GenericFixtureFactory.Create(EntityType)
	}
}
