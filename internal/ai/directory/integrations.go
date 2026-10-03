package directory

// integrations is the hardwired directory of Shopwell AI integrations. There is
// no remote source and nothing is fetched at runtime, so the directory lives
// directly in Go.
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
			Delivery:      Delivery{Kind: DeliveryBundled},
		},
		{
			Name:          "shopwell-cli-docker",
			DisplayName:   "Shopwell CLI (Docker)",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Run commands in Docker-backed Shopwell projects through Shopwell CLI.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/products/cli/",
			Delivery:      Delivery{Kind: DeliveryBundled},
		},
		{
			Name:          "deployment-helper",
			DisplayName:   "Shopwell Deployment Helper",
			Type:          TypeSkill,
			Provider:      "shopwell",
			Description:   "Use Shopwell CLI and Deployment Helper together for build and deploy workflows.",
			Status:        StatusActive,
			Documentation: "https://developer.shopwell.cn/docs/guides/hosting/installation-updates/deployments/deployment-helper/index.html",
			Delivery: Delivery{
				Kind:       DeliveryGit,
				Repository: "https://github.com/shopwell-shop/deployment-helper",
			},
			Compatibility: &Compatibility{Source: "owner"},
		},
	},
}
