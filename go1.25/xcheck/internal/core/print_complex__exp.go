package core

// HIGHLY EXPERIMENTAL STRAIGHTFORWARD DUMB INTERFACE

type parameterFormat interface {
	literal() string
	restrictions() string
}

func newParameterFormat(literal, restrictions string) parameterFormat {
	return &parameterFormatImpl{literal, restrictions}

}

type parameterFormatImpl struct {
	literal_      string
	restrictions_ string
}

func (pf *parameterFormatImpl) literal() string {
	return pf.literal_
}

func (pf *parameterFormatImpl) restrictions() string {
	return pf.restrictions_
}
