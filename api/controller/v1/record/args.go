package record

// CreateArgs create record args
type CreateArgs struct {
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Value    string `json:"value" validate:"required"`
	TTL      uint32 `json:"ttl"`
	Priority int    `json:"priority"`
}
