package operations

import (
	"fmt"
	"regexp"
	"strings"

	model "github.com/julianGoh17/simple-e2e/framework/models"
	"gopkg.in/yaml.v2"
)

/*
	Things to Accomplish:

	- Create mapping function which maps string to function
	- Create stage manager which determines which stages will be run (if any stages requested)
	- Create way a single manager which will manage the test run and determine when things test steps fail
	-
*/

// Controller is able to understand which stages and steps to run based on the test file. It is responsible for understanding if a test step has
// failed and will stop the test run prematurely if so.
type Controller struct {
	regexTestMethods   map[string]func(*model.Step)
	literalTestMethods map[string]func(*model.Step)
	procedure          *model.Procedure
}

// NewController is a constructor function which returns a pointer to the variable to work with
func NewController() *Controller {
	return &Controller{
		regexTestMethods:   make(map[string]func(*model.Step)),
		literalTestMethods: make(map[string]func(*model.Step)),
	}
}

// AddTestStep adds a Step Description and its associated method to the Controller so it knows what needs to do
func (controller *Controller) AddTestStep(description string, method func(*model.Step)) error {
	if isRegexDescription(description) {
		return controller.addRegexTestStep(description, method)
	}
	return controller.addLiteralTestStep(description, method)
}

func isRegexDescription(description string) bool {
	return strings.Contains(description, "'${string}'")
}

func (controller *Controller) addRegexTestStep(description string, method func(*model.Step)) error {
	parsedString := strings.ReplaceAll(description, "'${string}'", "('[A-Za-z]+')")
	_, err := regexp.Compile(parsedString)
	if err != nil {
		return err
	}

	if controller.regexTestMethods[parsedString] != nil {
		return fmt.Errorf("Error: Step description '%s' already exists", description)
	}
	controller.regexTestMethods[parsedString] = method

	return nil
}

func (controller *Controller) addLiteralTestStep(description string, method func(*model.Step)) error {
	if controller.literalTestMethods[description] != nil {
		return fmt.Errorf("Error: Step description '%s' already exists", description)
	}
	controller.literalTestMethods[description] = method
	return nil
}

// SetProcedure takes the read byte data from the test file and converts it to the Procedure object
func (controller *Controller) SetProcedure(procedureData []byte) error {
	procedure := &model.Procedure{}

	if err := yaml.UnmarshalStrict(procedureData, procedure); err != nil {
		return err
	}

	// TODO: Check for ill formatted testfile
	controller.procedure = procedure

	return nil
}
