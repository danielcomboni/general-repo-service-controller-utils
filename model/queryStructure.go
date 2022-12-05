package model

type QueryStructure struct {
	ParamName     string      `json:"paramName"`
	ParamValue    interface{} `json:"paramValue"`
	DbTableColumn string      `json:"dbTableColumn"`
}

