package testsuite_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/berlin-ab/testsuite"
)

func TestTestSuite(t *testing.T) {
	contents := []string{}

	//
	// setup assertion in cleanup block to ensure
	// teardown suite gets called before calling this
	// final cleanup
	//
	t.Cleanup(func() {
		expectedOrder := []string{
			// first suite in same toplevel test function
			"setup suite",

			// first top level Run executes
			"setup",
			"some test",
			"teardown",

			// second top level Run executes
			"setup",
			"some other test",
			"teardown",

			// first leve suite tears down
			"teardown suite",

			// second suite in same toplevel test function
			"setup other suite",
			"some other suite test",
			"teardown other suite",
		}

		actualOrder := contents
		require.EqualValues(t, expectedOrder, actualOrder)
	})

	testsuite.New(t, "Testing Test Suite", func(s *testsuite.S) {
		s.SetupSuite(func(t *testing.T) {
			contents = append(contents, "setup suite")
		})

		s.TeardownSuite(func(t *testing.T) {
			contents = append(contents, "teardown suite")
		})

		s.Setup(func(t *testing.T) {
			contents = append(contents, "setup")
		})

		s.Teardown(func(t *testing.T) {
			contents = append(contents, "teardown")
		})

		s.Run("some test", func(t *testing.T) {
			contents = append(contents, "some test")
		})

		s.Run("some other test", func(t *testing.T) {
			contents = append(contents, "some other test")
		})
	})

	testsuite.New(t, "Testing Other Test Suite", func(s *testsuite.S) {
		s.SetupSuite(func(t *testing.T) {
			contents = append(contents, "setup other suite")
		})

		s.TeardownSuite(func(t *testing.T) {
			contents = append(contents, "teardown other suite")
		})

		s.Run("some other suite test", func(t *testing.T) {
			contents = append(contents, "some other suite test")
		})
	})

	testsuite.New(t, "it defaults to no-op setup/teardown functions", func(s *testsuite.S) {
		s.Run("do nothing", func(t *testing.T) {
			require.True(t, true)
		})
	})

	testsuite.New(t, "it allows a testsuite to be skipped", func(s *testsuite.S) {
		s.SetupSuite(func(t *testing.T) {
			t.Skip()
		})

		s.Run("do nothing", func(t *testing.T) {
			require.True(t, false)
		})
	})

	testsuite.New(t, "it does nothing when no tests are specified", func(s *testsuite.S) {
		s.Setup(func(t *testing.T) {
			require.True(t, false)
		})

		s.SetupSuite(func(t *testing.T) {
			require.True(t, false)
		})

		s.Teardown(func(t *testing.T) {
			require.True(t, false)
		})

		s.TeardownSuite(func(t *testing.T) {
			require.True(t, false)
		})
	})

	testsuite.New(t, "calling hooks out of order panics", func(s *testsuite.S) {
		var paniced bool
		var originalPanic func(message string)

		s.SetupSuite(func(t *testing.T) {
			originalPanic = testsuite.PanicFunc()

			testsuite.SetPanicFunc(func(message string) {
				paniced = true
			})
		})

		s.Setup(func(t *testing.T) {
			paniced = false
		})

		s.TeardownSuite(func(t *testing.T) {
			testsuite.SetPanicFunc(originalPanic)
		})

		s.Run("it does not allow for setup to be defined after starting the tests", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.Run("foo", func(t *testing.T) {
				})

				s.Setup(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it does not allow for setup suite to be defined after starting the tests", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				paniced = false

				s.Run("foo", func(t *testing.T) {
				})

				s.SetupSuite(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it does not allow for teardown to be defined after starting the tests", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.Run("foo", func(t *testing.T) {
				})

				s.Teardown(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it does not allow for teardown suite to be defined after starting the tests", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.Run("foo", func(t *testing.T) {
				})

				s.TeardownSuite(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it considers calling hooks within run", func(t *testing.T) {
			testsuite.New(t, "it does not allow for teardown suite to be defined after starting the tests", func(s *testsuite.S) {
				s.Run("foo", func(t *testing.T) {
					s.Setup(func(t *testing.T) {
					})
				})
			})
			require.True(t, paniced)
		})
	})

	testsuite.New(t, "calling hooks twice panics", func(s *testsuite.S) {
		var paniced bool
		var originalPanic func(message string)

		s.SetupSuite(func(t *testing.T) {
			originalPanic = testsuite.PanicFunc()

			testsuite.SetPanicFunc(func(message string) {
				paniced = true
			})
		})

		s.Setup(func(t *testing.T) {
			paniced = false
		})

		s.TeardownSuite(func(t *testing.T) {
			testsuite.SetPanicFunc(originalPanic)
		})

		s.Run("it does not allow setup to be run twice", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.Setup(func(t *testing.T) {
				})

				s.Setup(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it does not allow teardown to be run twice", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.Teardown(func(t *testing.T) {
				})

				s.Teardown(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it does not allow setup suite to be run twice", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.SetupSuite(func(t *testing.T) {
				})

				s.SetupSuite(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})

		s.Run("it does not allow teardown suite to be run twice", func(t *testing.T) {
			testsuite.New(t, "suite", func(s *testsuite.S) {
				s.TeardownSuite(func(t *testing.T) {
				})

				s.TeardownSuite(func(t *testing.T) {
				})
			})

			require.True(t, paniced)
		})
	})

	t.Run("panic func panics", func(t *testing.T) {
		require.Panics(t, func() {
			testsuite.PanicFunc()("message")
		})
	})
}
