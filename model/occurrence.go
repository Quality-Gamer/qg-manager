package model

type Occurrence struct {
	Id               string               `json:"occurrence_id"`
	Description      string            `json:"description"`
	SolveOccurrences []SolveOccurrence `json:"solve"`
}

type ManagerOccurrence struct {
	Id         string     `json:"id"`
	Occurrence Occurrence `json:"occurrence"`
	Status     string     `json:"status"`
}

type UserOccurrence struct {
	Id         string     `json:"id"`
	Occurrence Occurrence `json:"occurrence"`
	Status     string     `json:"status"`
	Level      int        `json:"level"`
}

func LoadOccurrenceList() []Occurrence {
	//Futuramente virá de um microsservice que acessa o mysql
	var o1 Occurrence
	var o2 Occurrence
	var o3 Occurrence

	o1.Id = "1"
	o1.Description = "Um dos seus programadores backend está doente"

	o2.Id = "2"
	o2.Description = "O time atrasou uma entrega"

	o3.Id = "3"
	o3.Description = "O cliente pediu para adicionar novas funcionalidades sem alterar o prazo"

	var list []Occurrence
	list = append(list, o1, o2, o3)
	return list
}

func LoadUserOccurrenceList() []Occurrence {
	//Futuramente virá de um microsservice que acessa o mysql
	var o1 Occurrence
	var o2 Occurrence

	o1.Id = "1"
	o1.Description = "x dos seus programadores não tem IDE para trabalhar"

	o2.Id = "2"
	o2.Description = "x dos seus designers não tem Software para trabalhar"

	var list []Occurrence
	list = append(list, o1, o2)
	return list
}
