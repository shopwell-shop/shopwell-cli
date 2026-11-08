package tracking

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func enableVersionLookup(t *testing.T) {
	t.Helper()
	t.Setenv("DO_NOT_TRACK", "")
	require.NoError(t, os.Unsetenv("DO_NOT_TRACK"))
}

func TestShopwellVersionLookupDoesNotWait(t *testing.T) {
	enableVersionLookup(t)
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		lookup := startShopwellVersion("project.console", func() string {
			<-release
			return "6.7.11.0"
		})
		synctest.Wait()
		assert.Empty(t, lookup.Ready())
		close(release)
		// The worker must finish even though its result has not been consumed.
		synctest.Wait()
		assert.Equal(t, "6.7.11.0", lookup.Ready())
	})
}

func TestShopwellVersionLookupSkipsIneligibleCommands(t *testing.T) {
	enableVersionLookup(t)
	for _, name := range []string{"", "project", "project.create", "extension.build"} {
		lookup := startShopwellVersion(name, func() string {
			t.Error("unexpected lookup")
			return ""
		})
		assert.Nil(t, lookup)
		assert.Empty(t, lookup.Ready())
	}
}

func TestShopwellVersionLookupRespectsDoNotTrack(t *testing.T) {
	for _, value := range []string{"", "1", "0"} {
		t.Setenv("DO_NOT_TRACK", value)
		lookup := startShopwellVersion("project.console", func() string {
			t.Error("unexpected lookup")
			return ""
		})
		assert.Nil(t, lookup)
		assert.Empty(t, lookup.Ready())
	}
}

func TestProjectShopwellVersion(t *testing.T) {
	enableVersionLookup(t)
	root := t.TempDir()
	t.Setenv("PROJECT_ROOT", root)
	for version, want := range map[string]string{
		"v6.7.11.0":       "6.7.11.0",
		"v6.7.0.0-rc1":    "6.7.0.0-rc1",
		"dev-trunk":       "dev-trunk",
		"v6.7.0.0-RC1":    "",
		"6.6.x-dev":       "",
		"6.7.12.x-dev":    "",
		"dev-tmp-0c8ec10": "",
		"dev-feature-xyz": "",
		"6.7.11.0-dev":    "",
	} {
		t.Run(version, func(t *testing.T) {
			lock := fmt.Sprintf(`{"packages":[{"name":"shopwell/core","version":%q}],"packages-dev":[]}`, version)
			require.NoError(t, os.WriteFile(filepath.Join(root, "composer.lock"), []byte(lock), 0o644))
			assert.Equal(t, want, projectShopwellVersion())
		})
	}
	t.Setenv("PROJECT_ROOT", t.TempDir())
	assert.Empty(t, projectShopwellVersion())
	t.Setenv("PROJECT_ROOT", "")
	t.Chdir(t.TempDir())
	assert.Empty(t, projectShopwellVersion())
}

func TestReportableShopwellVersion(t *testing.T) {
	assert.Equal(t, "6.6.0.0", ReportableShopwellVersion("6.6.0.0"))
	assert.Equal(t, "dev-trunk", ReportableShopwellVersion("dev-trunk"))
	assert.Empty(t, ReportableShopwellVersion("dev-feature-xyz"))
}

func TestShopwellVersionLookupUnavailable(t *testing.T) {
	enableVersionLookup(t)
	synctest.Test(t, func(t *testing.T) {
		lookup := startShopwellVersion("project.console", func() string { return "" })
		synctest.Wait()
		assert.Empty(t, lookup.Ready())
	})
}
