package db

import "context"

type Session struct {
	UserName string
	Value    string
}

func (mgr *Mgr) CreateSession(session *Session) error {
	const query = "insert into sessions (user_name, value) values ($1, $2)"

	_, err := mgr.conn.Exec(context.Background(), query, session.UserName, session.Value)

	return err
}
