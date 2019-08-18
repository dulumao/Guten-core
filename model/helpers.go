package model

var m *Model

func M() *Model {
	if m == nil {
		m = &Model{}
	}

	return m
}
