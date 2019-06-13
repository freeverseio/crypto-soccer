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

func CountTeams() (int, error) {
	count := 0
	row := db.QueryRow("SELECT COUNT(*) FROM teams;")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}