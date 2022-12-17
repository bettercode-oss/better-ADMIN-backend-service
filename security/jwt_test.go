package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJwtToken_GetRefreshTokenExpiresForCookie(t *testing.T) {
	// given
	now := time.Now()
	token := JwtToken{
		AccessToken:         "test-access-token",
		RefreshToken:        "test-refresh-token",
		RefreshTokenExpires: now,
	}

	// when
	actual := token.GetRefreshTokenExpiresForCookie()

	// then
	// KST로 테스트함.
	expected := now.Add(9 * time.Hour)
	assert.Equal(t, expected, actual)
}
