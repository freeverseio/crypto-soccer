package storage

type Team struct {  
	ID    string  `db:"id"` 
	Name  string  `db:"name"`
}

func CreateTeam(id int, name string) error {
	return nil;
}