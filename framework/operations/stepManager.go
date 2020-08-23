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
	manager := &StepManager{
		regexTestMethods:   make(map[string]func(*model.Step) error),
		literalTestMethods: make(map[string]func(*model.Step) error),
	}
	for step, function := range getDefaultSteps() {
		manager.AddStepToManager(step, function)
	}

	return manager
}

// AddStepToManager adds a Step Description and its associated method to the StepManager so it knows what it needs to do
func (stepManager *StepManager) AddStepToManager(description string, method func(*model.Step) error) error {
	logger.Trace().
		Str("step", description).
		Msg("Adding step to step manager")
	if stepManager.isRegexDescription(description) {
		return stepManager.addRegexTestStep(description, method)
	}
	logger.Trace().
		Str("step", description).
		Bool("isRegex", false).
		Msg("Step is a literal description")
	return stepManager.addLiteralTestStep(description, method)
}

func (stepManager *StepManager) isRegexDescription(description string) bool {
	return strings.Contains(description, "'${string}'")
}

func (stepManager *StepManager) addRegexTestStep(description string, method func(*model.Step) error) error {
	logger.Trace().
		Str("step", description).
		Bool("isRegex", true).
		Msg("Adding step to regex test steps")
	parsedString := strings.ReplaceAll(description, "'${string}'", "('[A-Za-z]+')")
	_, err := regexp.Compile(parsedString)
	if err != nil {
		return err
	}

	if stepManager.regexTestMethods[parsedString] != nil {
		err := fmt.Errorf("Error: Step description '%s' already exists", description)
		logger.Error().Msg(err.Error())
		return err
	}
	stepManager.regexTestMethods[parsedString] = method
	logger.Trace().
		Str("parsedStep", parsedString).
		Msg("Added step to regex test steps")
	return nil
}

func (stepManager *StepManager) addLiteralTestStep(description string, method func(*model.Step) error) error {
	logger.Trace().
		Str("step", description).
		Bool("isRegex", false).
		Msg("Adding step to literal test steps")
	if stepManager.literalTestMethods[description] != nil {
		err := fmt.Errorf("Error: Step description '%s' already exists", description)
		logger.Error().Msg(err.Error())
		return err
	}
	stepManager.literalTestMethods[description] = method
	logger.Trace().
		Msg("Added step to literal test steps")
	return nil
}

// GetTestMethod will return the associated function based on the description string
func (stepManager *StepManager) GetTestMethod(description string) (func(*model.Step) error, error) {
	logger.Trace().
		Str("step", description).
		Msg("Retrieving step from step manager")
	if stepManager.isRegexTestDescription(description) {
		return stepManager.getRegexMethod(description)
	}
	return stepManager.getLiteralMethod(description)
}

func (stepManager *StepManager) getRegexMethod(description string) (func(*model.Step) error, error) {
	logger.Trace().
		Str("step", description).
		Bool("isRegex", true).
		Msg("Retrieving regex step from step manager")
	for regex, function := range stepManager.regexTestMethods {
		if matched, _ := regexp.MatchString(regex, description); matched {
			return function, nil
		}
	}
	err := fmt.Errorf("Could not find test that matches description: %s", description)
	logger.Error().
		Err(err).
		Str("step", description).
		Bool("isRegex", true).
		Msg(err.Error())
	return nil, err
}

func (stepManager *StepManager) getLiteralMethod(description string) (func(*model.Step) error, error) {
	logger.Trace().
		Str("step", description).
		Bool("isRegex", false).
		Msg("Retrieving literal step from step manager")
	function, ok := stepManager.literalTestMethods[description]
	if !ok {
		err := fmt.Errorf("Step '%s' is not registered in step list", description)
		logger.Error().
			Err(err).
			Str("step", description).
			Bool("isRegex", false).
			Msg("Could not find literal step in step manager")
		return nil, err
	}
	return function, nil
}

func (stepManager *StepManager) isRegexTestDescription(description string) bool {
	return strings.Contains(description, "'")
}
