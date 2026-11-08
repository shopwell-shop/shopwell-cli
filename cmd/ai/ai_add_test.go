package ai

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shopwell-shop/shopwell-cli/internal/ai/install"
)

func TestWriteAddResultInstalled(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, writeAddResult(&buf, install.AddResult{
		Name: "shopwell-cli", Agent: "claude-code", Scope: install.ScopeGlobal,
		ResolvedRevision: "v0.18.4",
	}))
	assert.Contains(t, buf.String(), "Installed shopwell-cli for claude-code (global) @v0.18.4")
}

func TestWriteAddResultDryRun(t *testing.T) {
	var buf bytes.Buffer
	require.NoError(t, writeAddResult(&buf, install.AddResult{
		Name: "shopwell-cli", Agent: "claude-code", Scope: install.ScopeGlobal,
		DryRun: true, Command: []string{"npx", "skills", "add"},
	}))
	assert.Contains(t, buf.String(), "[dry-run]")
	assert.Contains(t, buf.String(), "npx skills add")
}
