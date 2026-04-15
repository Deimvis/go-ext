package xstringscase

// TODO: change Case to interface defining formatting
// - camel impl should have option for whether capitalize first token

// NOTE: Case will be changed to interface,
// do not leverage its type;
// for formatting use Case.String().
type Case string

func (c Case) String() string {
	return string(c)
}

const (
	Camel Case = "camel"
	Snake Case = "snake"
	Kebab Case = "kebab"
)
