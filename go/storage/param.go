package storage

import "database/sql"

type Param struct {
	Name  string
	Value string
}

func Params(tx *sql.Tx) ([]Param, error) {
	rows, err := tx.Query("SELECT name, value FROM params;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var params []Param
	for rows.Next() {
		param := Param{}
		if err = rows.Scan(
			&param.Name,
			&param.Value,
		); err != nil {
			return nil, err
		}
		params = append(params, param)
	}

	return params, nil
}

func ParamByName(tx *sql.Tx, name string) (*Param, error) {
	rows, err := tx.Query("SELECT value FROM params WHERE name=$1;", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	param := Param{}
	param.Name = name
	if err = rows.Scan(
		&param.Value,
	); err != nil {
		return nil, err
	}

	return &param, nil
}

func (b Param) InsertOrUpdate(tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO params (name, value) VALUES ($1, $2)
						ON CONFLICT (name) DO UPDATE  
						SET name = $1 , value = $2;`, b.Name, b.Value)
	return err
}
