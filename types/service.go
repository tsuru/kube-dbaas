package types

type CreateArgs struct {
	Name        string                 `form:"name"`
	Team        string                 `form:"team"`
	Plan        string                 `form:"plan"`
	Description string                 `form:"description"`
	Tags        []string               `form:"tags"`
	Parameters  map[string]interface{} `form:"parameters"`
}
