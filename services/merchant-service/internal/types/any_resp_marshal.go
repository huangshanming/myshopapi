package types

import "encoding/json"

// MarshalJSON unwraps Data so OkJson matches legacy DTO-only bodies
// (not {"data": ...}). Used while many routes still declare returns (AnyResp).
func (a AnyResp) MarshalJSON() ([]byte, error) {
	if a.Data == nil {
		return []byte("null"), nil
	}
	return json.Marshal(a.Data)
}
