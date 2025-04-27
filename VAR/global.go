package Constant

import "fmt"

var TANENT = "kanaka"

type Position int

const (
	Member Position = iota
	Manager
	TechLead
)

type TechStack int

const (
	Golang TechStack = iota
	NodeJS
	Python
	Java
	React
	Angular
	Vue
	Postgres
	MongoDB
)

func (p Position) String() string {
	switch p {
	case Member:
		return "Member"
	case Manager:
		return "Manager"
	case TechLead:
		return "TechLead"
	default:
		return "Unknown"
	}
}

func (t TechStack) String() string {
	switch t {
	case Golang:
		return "Golang"
	case NodeJS:
		return "NodeJS"
	case Python:
		return "Python"
	case Java:
		return "Java"
	case React:
		return "React"
	case Angular:
		return "Angular"
	case Vue:
		return "Vue"
	case Postgres:
		return "Postgres"
	case MongoDB:
		return "MongoDB"
	default:
		return "Unknown"
	}
}
func ParseRole(role string) (int, error) {
	switch role {
	case "Member":
		return int(Member), nil
	case "Manager":
		return int(Manager), nil
	case "TechLead":
		return int(TechLead), nil
	default:
		return -1, fmt.Errorf("invalid role")
	}
}

func ParseTechnology(tech string) (int, error) {
	switch tech {
	case "Golang":
		return int(Golang), nil
	case "NodeJS":
		return int(NodeJS), nil
	case "Python":
		return int(Python), nil
	case "Java":
		return int(Java), nil
	case "React":
		return int(React), nil
	case "Angular":
		return int(Angular), nil
	case "Vue":
		return int(Vue), nil
	case "Postgres":
		return int(Postgres), nil
	case "MongoDB":
		return int(MongoDB), nil
	default:
		return -1, fmt.Errorf("invalid technology")
	}
}
