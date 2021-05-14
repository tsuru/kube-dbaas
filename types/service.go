package types

type CreateArgs struct {
	Name        string                 `form:"name"`
	Team        string                 `form:"team"`
	Plan        string                 `form:"plan"`
	Description string                 `form:"description"`
	Tags        []string               `form:"tags"`
	Parameters  map[string]interface{} `form:"parameters"`
}

type BindAppArgs struct {
	AppName          string   `form:"app-name"`
	AppHosts         []string `form:"app-hosts"`
	AppInternalHosts []string `form:"app-internal-hosts"`
	AppClusterName   string   `form:"app-cluster-name"`
	User             string   `form:"user"`
	EventID          string   `form:"eventid"`
}
