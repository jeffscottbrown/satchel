//go:build !production

package repository

import "testing"

// ConfigureRepositoryForTest sets the global Repository variable to the
// provided test repository. It also ensures that the original repository
// is restored after the test completes. This is useful for isolating
// tests and ensuring that they do not interfere with each other.
func ConfigureRepositoryForTest(t *testing.T, testRepo EmployeeRepository) {
	originalRepository := employeeRepository
	t.Cleanup(func() {
		employeeRepository = originalRepository
	})
	employeeRepository = testRepo
}
