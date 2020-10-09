package db

import "context"

type Session struct {
	UserName string
	Value    string
}

func (mgr *Mgr) GetSession(value string) (*Session, error) {
	const query = "select user_name from sessions where value = $1"

	session := Session{Value: value}
	err := mgr.conn.QueryRow(context.Background(), query, value).Scan(&session.UserName)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (mgr *Mgr) CreateSession(session *Session) error {
	const query = "insert into sessions (user_name, value) values ($1, $2)"

	_, err := mgr.conn.Exec(context.Background(), query, session.UserName, session.Value)

	return err
}
