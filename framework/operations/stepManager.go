package operations

import (
	"fmt"
	"regexp"
	"strings"

	model "github.com/julianGoh17/simple-e2e/framework/models"
)

// StepManager is the object that
type StepManager struct {
	regexTestMethods   map[string]func(*model.Step) error
	literalTestMethods map[string]func(*model.Step) error
}

// NewStepManager is the empty constructor which returns a functional StepManager to use
func NewStepManager() *StepManager {
	return &StepManager{
		regexTestMethods:   make(map[string]func(*model.Step) error),
		literalTestMethods: make(map[string]func(*model.Step) error),
	}
}

// AddStepToManager adds a Step Description and its associated method to the StepManager so it knows what it needs to do
func (stepManager *StepManager) AddStepToManager(description string, method func(*model.Step) error) error {
	if stepManager.isRegexDescription(description) {
		return stepManager.addRegexTestStep(description, method)
	}
	return stepManager.addLiteralTestStep(description, method)
}

func (stepManager *StepManager) isRegexDescription(description string) bool {
	return strings.Contains(description, "'${string}'")
}

func (stepManager *StepManager) addRegexTestStep(description string, method func(*model.Step) error) error {
	parsedString := strings.ReplaceAll(description, "'${string}'", "('[A-Za-z]+')")
	_, err := regexp.Compile(parsedString)
	if err != nil {
		return err
	}

	if stepManager.regexTestMethods[parsedString] != nil {
		return fmt.Errorf("Error: Step description '%s' already exists", description)
	}
	stepManager.regexTestMethods[parsedString] = method

	return nil
}

func (stepManager *StepManager) addLiteralTestStep(description string, method func(*model.Step) error) error {
	if stepManager.literalTestMethods[description] != nil {
		return fmt.Errorf("Error: Step description '%s' already exists", description)
	}
	stepManager.literalTestMethods[description] = method
	return nil
}

// GetTestMethod will return the associated function based on the description string
func (stepManager *StepManager) GetTestMethod(description string) (func(*model.Step) error, error) {
	if stepManager.isRegexTestDescription(description) {
		return stepManager.getRegexMethod(description)
	}
	return stepManager.getLiteralMethod(description)
}

func (stepManager *StepManager) getRegexMethod(description string) (func(*model.Step) error, error) {
	for regex, function := range stepManager.regexTestMethods {
		if matched, _ := regexp.MatchString(regex, description); matched {
			return function, nil
		}
	}
	return nil, fmt.Errorf("Could not find test that matches description: %s", description)
}

func (stepManager *StepManager) getLiteralMethod(description string) (func(*model.Step) error, error) {
	function := stepManager.literalTestMethods[description]
	if function == nil {
		return nil, fmt.Errorf("Step '%s' is not registered in step list", description)
	}
	return function, nil
}

func (stepManager *StepManager) isRegexTestDescription(description string) bool {
	return strings.Contains(description, "'")
}
