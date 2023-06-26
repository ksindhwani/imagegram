package filesystem

import (
	"testing"

	"github.com/ksindhwani/imagegram/pkg/config"
	"github.com/ksindhwani/imagegram/pkg/filesystem/local"
	"github.com/stretchr/testify/assert"
)

func TestNewFileSystem(t *testing.T) {
	tests := []struct {
		Name          string
		ExpectedError error
		Expected      FileSystem
	}{
		{
			Name:          "Test New Local FileSystem",
			ExpectedError: nil,
			Expected: &local.LocalFileSystem{
				HostDirectory:  "test host directory",
				LocalDirectory: "test local directory",
			},
		},
		{
			Name:          "Test New Default FileSystem",
			ExpectedError: nil,
			Expected: &local.LocalFileSystem{
				HostDirectory:  "test host directory",
				LocalDirectory: "test local directory",
			},
		},
	}

	config := config.Config{
		HostImageDirectory:  "test host directory",
		LocalImageDirectory: "test local directory",
	}

	for _, test := range tests {
		result, err := New(LOCAL, &config)
		assert.Equal(t, test.Expected, result, test.Name)
		assert.Equal(t, test.ExpectedError, err, test.Name)
	}
}
