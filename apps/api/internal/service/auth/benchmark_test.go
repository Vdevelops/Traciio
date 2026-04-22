package auth

import (
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// BenchmarkValidateToken benchmarks token validation (critical path)
func BenchmarkValidateToken(b *testing.B) {
	jwtManager := jwt.NewJWTManager(
		"test-secret-key-for-benchmarking-only",
		24*time.Hour,
		7*24*time.Hour,
	)

	// Generate a valid token
	token, _ := jwtManager.GenerateAccessToken("user-1", "test@example.com", "admin")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = jwtManager.ValidateToken(token)
	}
}

// BenchmarkGenerateToken benchmarks token generation
func BenchmarkGenerateToken(b *testing.B) {
	jwtManager := jwt.NewJWTManager(
		"test-secret-key-for-benchmarking-only",
		24*time.Hour,
		7*24*time.Hour,
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = jwtManager.GenerateAccessToken("user-1", "test@example.com", "admin")
	}
}

// BenchmarkPasswordHashing benchmarks bcrypt password hashing (expensive operation)
func BenchmarkPasswordHashing(b *testing.B) {
	password := "password123"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	}
}

// BenchmarkPasswordVerification benchmarks password verification
func BenchmarkPasswordVerification(b *testing.B) {
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	}
}

// BenchmarkPasswordVerificationWithDifferentCost benchmarks with different bcrypt costs
func BenchmarkPasswordVerificationWithCost10(b *testing.B) {
	benchmarkPasswordWithCost(b, 10)
}

func BenchmarkPasswordVerificationWithCost12(b *testing.B) {
	benchmarkPasswordWithCost(b, 12)
}

func benchmarkPasswordWithCost(b *testing.B, cost int) {
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), cost)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	}
}
