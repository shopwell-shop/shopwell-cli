package ai

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shopwell-shop/shopwell-cli/internal/ai/install"
)

func TestWriteRemoveResult(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, writeRemoveResult(&buf, install.RemoveResult{
		Name: "shopwell-cli", Agent: "claude-code", Scope: install.ScopeProject, Removed: true,
		Command: []string{"npx", "skills", "remove", "shopwell-cli"},
	}))
	assert.Contains(t, buf.String(), "Removed shopwell-cli for claude-code (project)")
}

func TestWriteRemoveResultNothingToRemove(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, writeRemoveResult(&buf, install.RemoveResult{
		Name: "shopwell-cli", Agent: "cursor", Scope: install.ScopeGlobal, Removed: false,
	}))
	assert.Contains(t, buf.String(), "nothing to remove")
}

func TestWriteRemoveResultDryRun(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, writeRemoveResult(&buf, install.RemoveResult{
		Name: "shopwell-cli", Agent: "claude-code", Scope: install.ScopeGlobal,
		Removed: true, DryRun: true, Command: []string{"npx", "skills", "remove", "shopwell-cli"},
	}))
	assert.Contains(t, buf.String(), "[dry-run] would remove")
}
