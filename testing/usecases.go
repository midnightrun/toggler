package testing

import (
	"github.com/adamluzsi/testcase"

	"github.com/toggler-io/toggler/usecases"
)

const ExampleUseCasesLetVar = `ExampleUseCasesLetVar`

func init() {
	setups = append(setups, func(s *testcase.Spec) {
		s.Let(ExampleUseCasesLetVar, func(t *testcase.T) interface{} {
			return usecases.NewUseCases(ExampleStorage(t))
		})
	})
}

func GetUseCases(t *testcase.T, varName string) *usecases.UseCases {
	return t.I(varName).(*usecases.UseCases)
}

func ExampleUseCases(t *testcase.T) *usecases.UseCases {
	return GetUseCases(t, ExampleUseCasesLetVar)
}
