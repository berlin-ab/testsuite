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
		require.EqualValues(t, contents,
			[]string{
				// first suite in same toplevel test function
				"setup suite",
				"setup",
				"some test",
				"teardown",
				"setup",
				"some other test",
				"teardown",
				"teardown suite",

				// second suite in same toplevel test function
				"setup other suite",
				"some other suite test",
				"teardown other suite",
			})
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
}
