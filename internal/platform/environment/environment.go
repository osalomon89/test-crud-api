package environment

type Environment int

const (
	Development Environment = iota
	Beta
	Production
)

func (d Environment) String() string {
	return [...]string{"development", "beta", "production"}[d]
}

func GetFromString(s string) Environment {
	switch s {
	case "production":
		return Production
	case "beta":
		return Beta
	default:
		return Development
	}
}
