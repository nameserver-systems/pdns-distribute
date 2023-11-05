package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_convertStringToInt(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			num, err := ConvertStringToInt("5")

			require.NoError(t, err)
			assert.Equal(t, 5, num)
		})

		t.Run("negative", func(t *testing.T) {
			num, err := ConvertStringToInt("-5")

			require.NoError(t, err)
			assert.Equal(t, -5, num)
		})

		t.Run("zero", func(t *testing.T) {
			num, err := ConvertStringToInt("0")

			require.NoError(t, err)
			assert.Zero(t, num)
		})
	})

	t.Run("fail", func(t *testing.T) {
		t.Run("no_number", func(t *testing.T) {
			num, err := ConvertStringToInt("abcd5")

			require.Error(t, err)
			assert.Empty(t, num)
		})

		t.Run("empty", func(t *testing.T) {
			num, err := ConvertStringToInt("")

			require.Error(t, err)
			assert.Empty(t, num)
		})
	})
}

func Test_trimAndLowerString(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			txt := TrimAndLowerString("Test Service")

			assert.Equal(t, "testservice", txt)
		})

		t.Run("equal-input-output", func(t *testing.T) {
			txt := TrimAndLowerString("testservice")

			assert.Equal(t, "testservice", txt)
		})

		t.Run("empty", func(t *testing.T) {
			txt := TrimAndLowerString("")

			assert.Empty(t, txt)
		})

		t.Run("with-hyphen", func(t *testing.T) {
			txt := TrimAndLowerString("Test-Service")

			assert.Equal(t, "test-service", txt)
		})

		t.Run("camelcase", func(t *testing.T) {
			txt := TrimAndLowerString("testService")

			assert.Equal(t, "testservice", txt)
		})
	})
}

func Test_trimSpace(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		t.Run("one_space", func(t *testing.T) {
			txt := trimSpace("test service")

			assert.Equal(t, "testservice", txt)
		})
		t.Run("with_hyphen", func(t *testing.T) {
			txt := trimSpace("test-service")

			assert.Equal(t, "test-service", txt)
		})
		t.Run("with_multiple_spaces", func(t *testing.T) {
			txt := trimSpace("    test         Service    ")

			assert.Equal(t, "testService", txt)
		})
		t.Run("empty", func(t *testing.T) {
			txt := trimSpace("")

			assert.Empty(t, txt)
		})
	})
}

func Test_GenerateUUID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		id, err := GenerateUUID()

		require.NoError(t, err)
		assert.NotNil(t, id)
		assert.Len(t, id, 36)
	})
}

func Test_getHashedTime(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		time := getHashedTime()
		assert.NotEmpty(t, time)
	})
}

func Test_EnsurePathExist(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tmpDir, dirErr := os.MkdirTemp("", "")
		require.NoError(t, dirErr)
		require.DirExists(t, tmpDir)

		defer os.RemoveAll(tmpDir)

		t.Run("create_dir", func(t *testing.T) {
			newTmpDir := filepath.Join(tmpDir, "new-dir")

			err := EnsurePathExist(newTmpDir)

			require.NoError(t, err)
			assert.DirExists(t, newTmpDir)
		})

		t.Run("dir_exists", func(t *testing.T) {
			err := EnsurePathExist(tmpDir)

			require.NoError(t, err)
			assert.DirExists(t, tmpDir)
		})
	})

	t.Run("fail", func(t *testing.T) {
		t.Run("no_path_given", func(t *testing.T) {
			err := EnsurePathExist("")

			require.EqualError(t, err, "mkdir : no such file or directory")
		})
	})
}
