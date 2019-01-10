package maps

type DistanceMatrix struct {
	DestinationAddress []string       `json:"destination_addresses"`
	OriginAddress      []string       `json:"origin_addresses"`
	Rows               []DistanceRows `json:"rows"`
	ErrorMessage       string         `json:"error_message"`
}

type DistanceRows struct {
	Elements []DistanceElements `json:"elements"`
}

type DistanceElements struct {
	Status   string   `json:"status"`
	Distance MetaData `json:"distance"`
	Duration MetaData `json:"duration"`
}

type MetaData struct {
	Text  string `json:"text"`
	Value int64  `json:"value"`
}
