package dtos

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestMemberAccessLog_GetHumanizeBrowserUserAgent(t *testing.T) {
	// given
	accessLog := MemberAccessLog{
		BrowserUserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	}

	// when
	actual := accessLog.GetHumanizeBrowserUserAgent()

	// then
	expected := "{\"browser\":{\"name\":\"Chrome\",\"version\":\"99.0.4844.51\"},\"engine\":{\"name\":\"AppleWebKit\",\"version\":\"537.36\"},\"mobile\":false,\"os\":\"Intel Mac OS X 10_15_7\",\"platform\":\"Macintosh\"}"
	assert.Equal(t, expected, actual)
}
