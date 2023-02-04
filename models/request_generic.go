package models

type GenericRequest[T any] struct {
	ParamNames  []QueryStructure       `json:"paramNames,omitempty"`
	Preloads    []string               `json:"preloads,omitempty"`
	QueryMap    map[string]interface{} `json:"queryMap,omitempty"`
	Property    []string               `json:"property,omitempty"`
	RequestType string                 `json:"requestType"`
	GeneralAuth Auth                   `json:"generalAuth,omitempty"`
}
