package utils

import (
	"encoding/json"
	"io"
)

func FromJSON(i interface{}, r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(i)
}

func ToJSON(i interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(i)
}
