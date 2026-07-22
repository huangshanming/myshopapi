package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexUint64 accepts JSON number or numeric string via encoding/json.
type FlexUint64 uint64

func (v FlexUint64) Uint64() uint64 { return uint64(v) }

func (v *FlexUint64) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*v = 0
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*v = 0
			return nil
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid id %q", s)
		}
		*v = FlexUint64(n)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		var u uint64
		if err2 := json.Unmarshal(b, &u); err2 != nil {
			return err
		}
		*v = FlexUint64(u)
		return nil
	}
	u, err := strconv.ParseUint(string(n), 10, 64)
	if err != nil {
		return err
	}
	*v = FlexUint64(u)
	return nil
}

func (v FlexUint64) MarshalJSON() ([]byte, error) {
	return strconv.AppendUint(nil, uint64(v), 10), nil
}
