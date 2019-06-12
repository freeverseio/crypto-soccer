package storage

type Team struct {  
	ID    string  `db:"id"` 
	Name  string  `db:"name"`
}

func CreateTeam(id int, name string) error {
	//  TODO: check for db is initialized
	_, err := db.Query("INSERT INTO teams (id, name) VALUES ($1, $2);", id, name)
	if err != nil {
	  	return err
	}

	return nil;
}