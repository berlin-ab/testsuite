package testsuite

import (
	"sync"
	"testing"
)

// New provides a new test suite that wraps the standard *testing.T behavior
// with xUnit setup/teardown behavior
func New(t *testing.T, suiteName string, userProvidedSuite func(*S)) {
	t.Run(suiteName, func(t *testing.T) {
		userProvidedSuite(&S{
			t: t,
			setupSuite: func(t *testing.T) {
				// no-op
			},
			teardownSuite: func(t *testing.T) {
				// no-op
			},
			setup: func(t *testing.T) {
				// no-op
			},
			teardown: func(t *testing.T) {
				// no-op
			},
		})
	})
}

// S holds the suite's configuration
type S struct {
	t             *testing.T
	setup         func(*testing.T)
	teardown      func(*testing.T)
	setupSuite    func(t *testing.T)
	teardownSuite func(t *testing.T)

	// lock
	performedSuiteSetup bool
	mu                  sync.Mutex
}

// Run performs a *testing.T Run behavior with the setup/teardown behavior of the
// testsuite
func (s2 *S) Run(name string, f func(*testing.T)) {
	// Ensure we only call setup suite and teardown suite once
	if s2.needsSetup() {
		s2.setupSuite(s2.t)
		s2.t.Cleanup(func() {
			s2.teardownSuite(s2.t)
		})
	}

	s2.t.Run(name, func(t *testing.T) {
		// configure this particular test with setup and teardown
		s2.setup(t)
		defer func() {
			s2.teardown(t)
		}()

		// run the specified test
		f(t)
	})
}

// Setup specifies behavior that should be run before each test in the suite
func (s2 *S) Setup(f func(t *testing.T)) {
	s2.setup = f
}

// Teardown specifies behavior that should be run after each test in the suite
func (s2 *S) Teardown(f func(t *testing.T)) {
	s2.teardown = f
}

// SetupSuite specifies behavior that should be run before running any tests in the suite
func (s2 *S) SetupSuite(f func(t *testing.T)) {
	s2.setupSuite = f
}

// TeardownSuite specifies behavior that should be run after running all tests in the suite
func (s2 *S) TeardownSuite(f func(t *testing.T)) {
	s2.teardownSuite = f
}

func (s2 *S) needsSetup() bool {
	// double check that the reader always gets the right value
	s2.mu.Lock()
	defer func() {
		s2.mu.Unlock()
	}()

	val := s2.performedSuiteSetup
	s2.performedSuiteSetup = true
	return !val
}
