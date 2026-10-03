package docker

// buildSearch builds the OpenSearch node the shop indexes into when
// shopwell/elasticsearch is installed.
func buildSearch(*Environment, service) composeService {
	return composeService{
		Image: "opensearchproject/opensearch:2",
		Environment: yamlMap[string]{}.
			set("OPENSEARCH_INITIAL_ADMIN_PASSWORD", "Shopwell123!").
			set("discovery.type", "single-node").
			set("plugins.security.disabled", "true"),
		Volumes: []string{"opensearch-data:/usr/share/opensearch/data"},
	}
}
