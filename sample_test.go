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
