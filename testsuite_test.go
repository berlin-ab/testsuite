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

			// whole inner suite runs
			"setup",
			"setup inner suite",

			"setup inner",
			"test inner",
			"teardown inner",

			"setup inner",
			"test another inner",
			"teardown inner",

			"teardown inner suite",
			"teardown",

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

		s.When("a subcontext happens", func(s *testsuite.S) {
			s.SetupSuite(func(t *testing.T) {
				contents = append(contents, "setup inner suite")
			})

			s.TeardownSuite(func(t *testing.T) {
				contents = append(contents, "teardown inner suite")
			})

			s.Setup(func(t *testing.T) {
				contents = append(contents, "setup inner")
			})

			s.Teardown(func(t *testing.T) {
				contents = append(contents, "teardown inner")
			})

			s.Run("inner test", func(t *testing.T) {
				contents = append(contents, "test inner")
			})

			s.Run("another inner test", func(t *testing.T) {
				contents = append(contents, "test another inner")
			})
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

	testsuite.New(t, "it allows for sub-suites", func(s *testsuite.S) {
		var topLevel string
		var midLevel string
		var leafLevel string

		t.Cleanup(func() {
			require.Equal(t, topLevel, "")
			require.Equal(t, midLevel, "")
			require.Equal(t, leafLevel, "")
		})

		s.Setup(func(t *testing.T) {
			topLevel = "top level"
		})

		s.TeardownSuite(func(t *testing.T) {
			topLevel = ""
		})

		s.When("a sub suite is defined", func(s *testsuite.S) {
			s.Setup(func(t *testing.T) {
				midLevel = "mid level"
			})

			s.TeardownSuite(func(t *testing.T) {
				midLevel = ""
			})

			s.When("a sub sub suite is defined", func(s *testsuite.S) {
				s.Setup(func(t *testing.T) {
					leafLevel = "leaf level"
				})

				s.TeardownSuite(func(t *testing.T) {
					leafLevel = ""
				})

				s.Run("a leaf test", func(t *testing.T) {
					require.Equal(t, topLevel, "top level")
					require.Equal(t, midLevel, "mid level")
					require.Equal(t, leafLevel, "leaf level")
				})
			})
		})
	})
}
