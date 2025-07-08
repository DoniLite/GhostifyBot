package utils_test

import (
	"fmt"
	"testing"

	"github.com/DoniLite/GhostifyBot/utils"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	mockError := fmt.Errorf("something went wrong with the value of %d", 20)
	reporter := utils.CreateNewReport()
	t.Run("error creation and priority", func(t *testing.T) {
		reporter.Priority = utils.HIGH
		priority_fail_meta := fmt.Sprintf("The expected Priority %s is not provided on the object \n %+v", *utils.HIGH, reporter)
		assert.Equal(t, utils.HIGH, reporter.Priority, priority_fail_meta)
	})
	t.Run("error persistence", func(t *testing.T) {
		reporter.Err = mockError.Error()
		err := reporter.PersistReport()
		assert.Equal(t, nil, err, "Failed to persist the error reporter")
	})
}
