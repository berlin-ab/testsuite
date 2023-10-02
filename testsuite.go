package testsuite

import (
	"fmt"
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
	teardownSuite func(t *testing.T)

	testCaseDefined     bool
	setupCalled         bool
	teardownCalled      bool
	setupSuiteCalled    bool
	teardownSuiteCalled bool
}

// Run performs a *testing.T Run behavior with the setup/teardown behavior of the
// testsuite
func (s *S) Run(name string, f func(*testing.T)) {
	if s.needsSetup() {
		s.observeSetup()
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

// Setup specifies behavior that should be run before each test in the suite
//
// note: must be specified before Run()
func (s *S) Setup(f func(t *testing.T)) {
	const hookName = "Setup()"

	s.preventSetupRedefinition(&s.setupCalled, hookName)
	s.preventHookMisuseAfterStart(hookName)

	s.setup = f
}

// Teardown specifies behavior that should be run after each test in the suite
//
// note: must be specified before Run()
func (s *S) Teardown(f func(t *testing.T)) {
	const hookName = "Teardown()"

	s.preventSetupRedefinition(&s.teardownCalled, hookName)
	s.preventHookMisuseAfterStart(hookName)

	s.teardown = f
}

// SetupSuite specifies behavior that should be run before running any tests in the suite
//
// note: must be specified before Run()
func (s *S) SetupSuite(f func(t *testing.T)) {
	const hookName = "SetupSuite()"

	s.preventSetupRedefinition(&s.setupSuiteCalled, hookName)
	s.preventHookMisuseAfterStart(hookName)

	s.setupSuite = f
}

// TeardownSuite specifies behavior that should be run after running all tests in the suite
//
// note: must be specified before Run()
func (s *S) TeardownSuite(f func(t *testing.T)) {
	const hookName = "TeardownSuite()"

	s.preventSetupRedefinition(&s.teardownSuiteCalled, hookName)
	s.preventHookMisuseAfterStart(hookName)

	s.teardownSuite = f
}

func (s *S) preventSetupRedefinition(val *bool, hookName string) {
	if *val {
		s.panicHookRedefinition(hookName)
	} else {
		*val = true
	}
}

func (s *S) preventHookMisuseAfterStart(hook string) {
	if s.testCaseDefined {
		panicFunc(fmt.Sprintf("%v called after Run() in testsuite", hook))
	}
}

func (s *S) needsSetup() bool {
	// only report needsSetup=true the first time a test case is
	// defined in a suite
	return !s.testCaseDefined
}

func (s *S) observeSetup() {
	s.testCaseDefined = true
}

func (s *S) panicHookRedefinition(hook string) {
	panicFunc(fmt.Sprintf("%v called twice. %v can only be called once", hook, hook))
}

func runSuite(t *testing.T, userProvidedSuite func(*S)) {
	suite := newSuite(t)

	defer func() {
		if suite.testCaseDefined {
			suite.teardownSuite(t)
		}
	}()

	userProvidedSuite(suite)
}

func newSuite(t *testing.T) *S {
	return &S{
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
	}
}

// panicFunc allows the tests to observe panics without actually
// panicking.
var panicFunc = func(message string) {
	panic(message)
}
