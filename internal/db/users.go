package db

import "context"

type User struct {
	Name  string
	Photo *string
}

func (mgr *Mgr) EnsureUserExists(user *User) error {
	const query = "insert into users (name, photo) values ($1, $2) on conflict do nothing"

	_, err := mgr.conn.Exec(context.Background(), query, user.Name, user.Photo)

	return err
}
