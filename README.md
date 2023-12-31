# testsuite

[![Go Reference](https://pkg.go.dev/badge/github.com/berlin-ab/testsuite.svg)](https://pkg.go.dev/github.com/berlin-ab/testsuite)
[![Go Report Card](https://goreportcard.com/badge/github.com/berlin-ab/testsuite)](https://goreportcard.com/report/github.com/berlin-ab/testsuite)
[![Coverage](coverage_badge.png)](coverage-report.html)

testsuite provides xUnit style test suite setup and teardown behavior
for golang's testing.T library missing from the Go ecosystem:

https://awesome-go.com/testing/

## Features

- SetupSuite/TeardownSuite
- Setup/Teardown per test
- Multiple suites per test function
- Simple abstraction wrapping standard testing.T
- Enforces simplicity by not allowing multiple definitions of hooks
- Enforces simplicity by not allowing sub suites which can make it difficult to follow changes to the state of the test

## Usage: 

```golang
package testsuite_test

import (
    "testing"

    "github.com/berlin-ab/testsuite"
)

func TestExample(t *testing.T) {
    testsuite.New(t, "defines a test suite with setup/teardown", func(s *testsuite.S) {
        s.SetupSuite(func(t *testing.T) {
            // executes before all tests
        })

        s.Setup(func(t *testing.T) {
            // executes before each test
        })

        s.Teardown(func(t *testing.T) {
            // executes after each test
        })

        s.TeardownSuite(func(t *testing.T) {
            // executes after all tests
        })

        s.Run("it does something", func(t *testing.T) {
            // defines a test using standard *testing.T
        })
    })

    testsuite.New(t, "another suite", func(s *testsuite.S) {
        s.Run("it does something else", func(t *testing.T) {
            // defines a test using standard *testing.T
        })
    })
}
```


## Example:

```golang

package testsuite_test

import (
    "testing"

    "github.com/berlin-ab/testsuite"
)

func TestSample(t *testing.T) {
    testsuite.New(t, "first example", func(s *testsuite.S) {
        s.SetupSuite(func(t *testing.T) {
            t.Log("setup first suite")
        })

        s.TeardownSuite(func(t *testing.T) {
            t.Log("teardown first suite")
        })

        s.Setup(func(t *testing.T) {
            t.Log("setup test")
        })

        s.Teardown(func(t *testing.T) {
            t.Log("teardown test")
        })

        s.Run("a test", func(t *testing.T) {
            t.Log("running test")
        })
    })

    testsuite.New(t, "second example", func(s *testsuite.S) {
        s.SetupSuite(func(t *testing.T) {
            t.Log("setup second suite")
        })

        s.TeardownSuite(func(t *testing.T) {
            t.Log("teardown second suite")
        })

        s.Setup(func(t *testing.T) {
            t.Log("setup test")
        })

        s.Teardown(func(t *testing.T) {
            t.Log("teardown test")
        })

        s.Run("a test", func(t *testing.T) {
            t.Log("running test")
        })
    })
}

$ go test -v ./sample_test.go | grep sample_test
sample_test.go:12: setup first suite
sample_test.go:20: setup test
sample_test.go:28: running test
sample_test.go:24: teardown test
sample_test.go:16: teardown first suite
sample_test.go:34: setup second suite
sample_test.go:42: setup test
sample_test.go:50: running test
sample_test.go:46: teardown test
sample_test.go:38: teardown second suite
```
