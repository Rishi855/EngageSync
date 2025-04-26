package Constant

var TANENT = "kanaka"

type Position int

const (
	Member   Position = iota
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
