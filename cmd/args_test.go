package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapAliasArgs_EmptyArgv(t *testing.T) {
	t.Parallel()
	assert.Nil(t, mapAliasArgs(nil))
	assert.Nil(t, mapAliasArgs([]string{}))
}

func TestMapAliasArgs_NoArgs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{}, mapAliasArgs([]string{"shopwell-cli"}))
}

func TestMapAliasArgs_RegularBinary(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"shopwell-cli", "project", "console", "debug:router"})

	assert.Equal(t, []string{"project", "console", "debug:router"}, args)
}

func TestMapAliasArgs_RenamedBinaryIsNotAliased(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/opt/bin/custom-cli", "cache:clear"})

	assert.Equal(t, []string{"cache:clear"}, args)
}

func TestMapAliasArgs_SwxAlias(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/usr/local/bin/swx", "debug:router", "--env=prod"})

	assert.Equal(t, []string{"project", "console", "debug:router", "--env=prod"}, args)
}

func TestMapAliasArgs_SwxAliasWithoutArgs(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/usr/local/bin/swx"})

	assert.Equal(t, []string{"project", "console", "list"}, args)
}

func TestMapAliasArgs_SwxExeAlias(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"C:\\tools\\swx.exe", "cache:clear"})

	assert.Equal(t, []string{"project", "console", "cache:clear"}, args)
}

func TestMapAliasArgs_SwxCaseInsensitive(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/usr/local/bin/SWX", "cache:clear"})

	assert.Equal(t, []string{"project", "console", "cache:clear"}, args)
}

func TestMapAliasArgs_SwxCompletion(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/usr/local/bin/swx", "completion", "bash"})

	assert.Equal(t, []string{"completion", "bash"}, args)
}

func TestMapAliasArgs_SwxInternalCompletion(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/usr/local/bin/swx", "__complete", "cache:clear"})

	assert.Equal(t, []string{"__complete", "project", "console", "cache:clear"}, args)
}

func TestMapAliasArgs_SwxInternalCompletionNoDesc(t *testing.T) {
	t.Parallel()
	args := mapAliasArgs([]string{"/usr/local/bin/swx", "__completeNoDesc", "cache:clear"})

	assert.Equal(t, []string{"__completeNoDesc", "project", "console", "cache:clear"}, args)
}

func TestMapAliasArgs_SwxHelp(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"project", "console", "--help"}, mapAliasArgs([]string{"/usr/local/bin/swx", "--help"}))
	assert.Equal(t, []string{"project", "console", "-h"}, mapAliasArgs([]string{"/usr/local/bin/swx", "-h"}))
}

func TestMapAliasArgs_SwxVersion(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"project", "console", "--version"}, mapAliasArgs([]string{"/usr/local/bin/swx", "--version"}))
	assert.Equal(t, []string{"project", "console", "-v"}, mapAliasArgs([]string{"/usr/local/bin/swx", "-v"}))
}

func TestMapAliasArgs_DoesNotMutateInput(t *testing.T) {
	t.Parallel()
	argv := []string{"/usr/local/bin/swx", "cache:clear"}

	mapAliasArgs(argv)

	assert.Equal(t, []string{"/usr/local/bin/swx", "cache:clear"}, argv)
}

func TestIsSwxAlias(t *testing.T) {
	t.Parallel()
	assert.True(t, isSwxAlias("swx"))
	assert.True(t, isSwxAlias("/usr/local/bin/swx"))
	assert.True(t, isSwxAlias(`C:\tools\Swx.exe`))
	assert.False(t, isSwxAlias("shopwell-cli"))
	assert.False(t, isSwxAlias("/usr/local/bin/swx-dev"))
	assert.False(t, isSwxAlias(""))
}

func TestCommandNameFromArgs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "shopwell-cli", commandNameFromArgs([]string{"/usr/local/bin/shopwell-cli"}))
	assert.Equal(t, "swx", commandNameFromArgs([]string{"C:\\tools\\swx.exe"}))
	assert.Equal(t, "custom-cli", commandNameFromArgs([]string{"custom-cli", "project"}))
	assert.Equal(t, "shopwell-cli", commandNameFromArgs(nil))
}

func TestCommandNameFromBinaryPath(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "shopwell-cli", commandNameFromBinaryPath("shopwell-cli"))
	assert.Equal(t, "shopwell-cli", commandNameFromBinaryPath("/usr/local/bin/shopwell-cli"))
	assert.Equal(t, "shopwell-cli", commandNameFromBinaryPath(`C:\Program Files\shopwell-cli.exe`))
	assert.Equal(t, "swx", commandNameFromBinaryPath("./swx"))
	assert.Equal(t, "shopwell-cli", commandNameFromBinaryPath(""), "empty path falls back to the root command name")
	assert.Equal(t, "shopwell-cli", commandNameFromBinaryPath(".exe"), "extension-only path falls back to the root command name")
}
