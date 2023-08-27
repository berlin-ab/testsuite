package testsuite

import (
	"sync"
	"testing"
)

// New provides a new test suite that wraps the standard *testing.T behavior
// with xUnit setup/teardown behavior
func New(t *testing.T, suiteName string, userProvidedSuite func(*S)) {
	t.Run(suiteName, func(t *testing.T) {
		runSuite(t, userProvidedSuite)
	})
}

// S holds the suite's configuration
type S struct {
	t             *testing.T
	setup         func(*testing.T)
	teardown      func(*testing.T)
	setupSuite    func(t *testing.T)
	teardownSuite *func(t *testing.T)

	// lock
	testCaseDefined bool
	mu              sync.Mutex
}

// Run performs a *testing.T Run behavior with the setup/teardown behavior of the
// testsuite
func (s *S) Run(name string, f func(*testing.T)) {
	if s.needsSetup() {
		s.setupSuite(s.t)
	}

	s.t.Run(name, func(t *testing.T) {
		// configure this particular test with setup and teardown
		s.setup(t)

		defer s.teardown(t)

		// run the specified test
		f(t)
	})
}

// When specifies a nested suite
func (s *S) When(context string, userProvidedContext func(s *S)) {
	s.Run(context, func(t *testing.T) {
		runSuite(t, userProvidedContext)
	})
}

// Setup specifies behavior that should be run before each test in the suite
func (s *S) Setup(f func(t *testing.T)) {
	s.setup = f
}

// Teardown specifies behavior that should be run after each test in the suite
func (s *S) Teardown(f func(t *testing.T)) {
	s.teardown = f
}

// SetupSuite specifies behavior that should be run before running any tests in the suite
func (s *S) SetupSuite(f func(t *testing.T)) {
	s.setupSuite = f
}

// TeardownSuite specifies behavior that should be run after running all tests in the suite
func (s *S) TeardownSuite(f func(t *testing.T)) {
	s.teardownSuite = &f
}

func (s *S) needsSetup() bool {
	// double check that the reader always gets the right value
	s.mu.Lock()
	defer s.mu.Unlock()

	// only report needsSetup=true the first time a test case is
	// defined in a suite
	val := s.testCaseDefined
	s.testCaseDefined = true
	return !val
}

func runSuite(t *testing.T, userProvidedSuite func(*S)) {
	suite := newSuite(t)

	defer func() {
		if suite.teardownSuite == nil {
			return // no teardown specified
		}

		if !suite.testCaseDefined {
			return // no tests specified, no need to teardown
		}

		teardownSuite := *suite.teardownSuite
		teardownSuite(t)
	}()

	userProvidedSuite(suite)
}

func newSuite(t *testing.T) *S {
	return &S{
		t: t,
		setupSuite: func(t *testing.T) {
			// no-op
		},
		teardownSuite: nil, // no-op
		setup: func(t *testing.T) {
			// no-op
		},
		teardown: func(t *testing.T) {
			// no-op
		},
	}
}
