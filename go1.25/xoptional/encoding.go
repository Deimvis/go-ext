package xoptional

import (
	"bytes"
	"encoding/json"

	yaml_v3 "gopkg.in/yaml.v3"
)

func (o T[U]) MarshalJSON() ([]byte, error) {
	if !o.set {
		return []byte("null"), nil
	}
	return json.Marshal(o.v)
}

func (o *T[U]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		o.set = false
		return nil
	}
	err := json.Unmarshal(b, &o.v)
	if err != nil {
		return err
	}
	o.set = true
	return nil
}

func (o T[U]) MarshalYAML() (any, error) {
	if !o.set {
		return nil, nil
	}
	return o.v, nil
}

func (o *T[U]) UnmarshalYAML(node *yaml_v3.Node) error {
	if node.IsZero() {
		o.set = false
		return nil
	}

	err := node.Decode(&o.v)
	if err != nil {
		return err
	}
	o.set = true
	return nil
}
