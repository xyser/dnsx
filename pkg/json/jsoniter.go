package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// Marshal is exported by gin/json package.
	Marshal = json.Marshal
	// MarshalToString marshal struct to string
	MarshalToString = jsoniter.MarshalToString

	// Unmarshal is exported by gin/json package.
	Unmarshal = json.Unmarshal

	// UnmarshalFromString is exported by module/json package.
	UnmarshalFromString = json.UnmarshalFromString

	// MarshalIndent is exported by gin/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by gin/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by gin/json package.
	NewEncoder = json.NewEncoder
)
