# testsuite

testsuite provides xUnit style test suite setup and teardown behavior
for golang's testing.T library.

## Usage: 

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
    }
