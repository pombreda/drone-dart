package server

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {

	sdk, err := version.NewVersion("2.5.3")
	if err != nil {
		t.Error(err)
	}

	parts := strings.Split(">=0.8.10+6 <2.0.0", " ")
	for _, part := range parts {
		constraints, err := version.NewConstraint(part)
		if err != nil {
			t.Error(err)
		}

		result := constraints.Check(sdk)
		fmt.Println(result)
	}
}
