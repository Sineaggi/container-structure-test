package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleContainerTools/container-structure-test/pkg/types/unversioned"
)

func TestFinalResults(t *testing.T) {
	t.Parallel()

	result := unversioned.SummaryObject{
		Pass:     1,
		Fail:     1,
		Total:    2,
		Duration: time.Duration(2),
		Results: []*unversioned.TestResult{
			{
				Name:     "my first test",
				Pass:     true,
				Stdout:   "it works!",
				Stderr:   "",
				Duration: time.Duration(1),
			},
			{
				Name:     "my fail",
				Pass:     false,
				Stdout:   "",
				Stderr:   "this failed",
				Errors:   []string{"this failed because of that"},
				Duration: time.Duration(1),
			},
		},
	}

	var finalResultsTests = []struct {
		actual   *bytes.Buffer
		format   unversioned.OutputValue
		expected string
	}{
		{
			actual:   bytes.NewBuffer([]byte{}),
			format:   unversioned.Junit,
			expected: `<?xml version="1.0" encoding="UTF-8"?><testsuites failures="1" tests="2" time="2e-09"><testsuite name="container-structure-test.test"><testcase name="my first test" time="1e-09"><system-out>it works!</system-out><system-err></system-err></testcase><testcase name="my fail" time="1e-09"><failure>this failed because of that</failure><system-out></system-out><system-err>this failed</system-err></testcase></testsuite></testsuites>`,
		},
		{
			actual:   bytes.NewBuffer([]byte{}),
			format:   unversioned.Json,
			expected: `{"Pass":1,"Fail":1,"Total":2,"Duration":2,"Results":[{"Name":"my first test","Pass":true,"Stdout":"it works!","Duration":1},{"Name":"my fail","Pass":false,"Stderr":"this failed","Errors":["this failed because of that"],"Duration":1}]}`,
		},
	}

	for _, test := range finalResultsTests {
		test := test

		t.Run(test.format.String(), func(t *testing.T) {
			t.Parallel()

			FinalResults(test.actual, test.format, result)

			if strings.TrimSpace(test.actual.String()) != test.expected {
				t.Errorf("expected %s but got %s", test.expected, test.actual)
			}
		})
	}
}
