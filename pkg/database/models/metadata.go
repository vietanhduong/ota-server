package models

type Metadata struct {
	base
	Type      string
	ProfileId uint  
	Key       string
	Value     string
}
