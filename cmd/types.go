package cmd

import (
	"time"

	"github.com/joshdk/go-junit"
)

// TestDocument represents the results of a single test run.
type TestDocument struct {
	// SuiteName is a descriptor to the suite the current test belongs to
	SuiteName string `json:"suitename" yaml:"name"`

	// Name is a descriptor given to the test.
	Name string `json:"name" yaml:"name"`

	// Classname is an additional descriptor for the hierarchy of the test.
	Classname string `json:"classname" yaml:"classname"`

	// Published is a timestamp to determine when the test result has been stored in elasticsearch.
	Published time.Time `json:"published" yaml:"published"`

	// Duration is the total time taken to run the tests.
	Duration time.Duration `json:"duration" yaml:"duration"`

	// Status is the result of the test. Status values are passed, skipped,
	// failure, & error.
	Status junit.Status `json:"status" yaml:"status"`

	// Error is a record of the failure or error of a test, if applicable.
	//
	// The following relations should hold true.
	//   Error == nil && (Status == Passed || Status == Skipped)
	//   Error != nil && (Status == Failed || Status == Error)
	Error error `json:"error" yaml:"error"`

	// Additional properties from XML node attributes.
	// Some tools use them to store additional information about test location.
	Properties map[string]string `json:"properties" yaml:"properties"`

	// SystemOut is textual output for the test case. Usually output that is
	// written to stdout.
	SystemOut string `json:"stdout,omitempty" yaml:"stdout,omitempty"`

	// SystemErr is textual error output for the test case. Usually output that is
	// written to stderr.
	SystemErr string `json:"stderr,omitempty" yaml:"stderr,omitempty"`
}

// Error represents an erroneous test result.
type Error struct {
	// Message is a descriptor given to the error. Purpose and values differ by
	// environment.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Type is a descriptor given to the error. Purpose and values differ by
	// framework. Value is typically an exception class, such as an assertion.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Body is extended text for the error. Purpose and values differ by
	// framework. Value is typically a stacktrace.
	Body string `json:"body,omitempty" yaml:"body,omitempty"`
}

// Error returns a textual description of the test error.
func (err Error) Error() string {
	return err.Body
}
