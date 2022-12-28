package plugin

type monitoringRequest struct {
	FolderID    string `json:"folderId"`
	Aggregation string `json:"aggregation"`
	Alias       string `json:"alias"`
	QueryText   string `json:"queryText"`
}

const apiKeyJsonInSettings = "apiKeyJson"

type monitoringConfig struct {
	APIEndpoint        string `json:"apiEndpoint"`
	MonitoringEndpoing string `json:"monitoringEndpoint"`
	FolderID           string `json:"folderId"`
}
