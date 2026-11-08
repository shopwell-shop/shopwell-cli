package directory

// integrations is the hardwired directory of Shopwell AI integrations. There is
// no remote source and nothing is fetched at runtime, so the directory lives
// directly in Go. Repository is the GitHub "owner/repo" the skill is installed
// from; the shopwell-cli skills ship in this repository.
var integrations = Directory{
	Integrations: []Integration{
		{
			Name:          "shopwell-cli",
			DisplayName:   "Shopwell CLI",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Use Shopwell CLI effectively for project, extension, and account workflows.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/products/cli/",
			Delivery:      Delivery{Kind: DeliveryGit, Repository: "shopwell-shop/shopwell-cli"},
		},
		{
			Name:          "shopwell-cli-docker",
			DisplayName:   "Shopwell CLI (Docker)",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Run commands in Docker-backed Shopwell projects through Shopwell CLI.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/products/cli/",
			Delivery:      Delivery{Kind: DeliveryGit, Repository: "shopwell-shop/shopwell-cli"},
		},
		{
			Name:          "deployment-helper",
			DisplayName:   "Shopwell Deployment Helper",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Use Shopwell CLI and Deployment Helper together for build and deploy workflows.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/guides/hosting/installation-updates/deployments/deployment-helper/index.html",
			Delivery:      Delivery{Kind: DeliveryGit, Repository: "shopwell-shop/deployment-helper"},
		},
		{
			Name:          "acceptance-test-suite",
			DisplayName:   "Shopwell Acceptance Test Suite",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Write and run Playwright acceptance and accessibility tests for a Shopwell project with the Acceptance Test Suite.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/guides/development/testing/e2e-playwright/",
			Delivery:      Delivery{Kind: DeliveryGit, Repository: "shopwell-shop/acceptance-test-suite"},
		},
		{
			Name:          "shopwell-cli-extension-anonymize",
			DisplayName:   "Shopwell Extension Anonymize",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Bootstrap an extension's anonymize config so project dump --anonymize clears its personal data and secrets.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/products/cli/",
			Delivery:      Delivery{Kind: DeliveryGit, Repository: "shopwell-shop/shopwell-cli"},
		},
		{
			Name:          "meteor",
			DisplayName:   "Meteor MCP Server",
			Type:          TypeMCP,
			Provider:      "shopwell",
			Description:   "Connect your AI agent to the Meteor MCP server for Shopwell development assistance.",
			Status:        StatusActive,
			Documentation: "https://meteor.shopwell.cn/documentation/getting-started/agents",
			Delivery: Delivery{
				Kind:      DeliveryMCP,
				URL:       "https://meteor.shopwell.cn/mcp",
				Transport: "http",
			},
		},
	},
}
