package main

import (
	"context"
	"database/sql"
)

type UserContact struct {
	ID    string
	Email string
	Phone string
}

type ContactStore struct {
	Pool *sql.DB
}

func (s *ContactStore) ListUserContactsNPlusOne(ctx context.Context, userIDs []string) ([]UserContact, error) {
	var contacts []UserContact
	for _, userID := range userIDs {
		var c UserContact
		err := s.Pool.QueryRowContext(ctx, "SELECT id, email, phone FROM users WHERE id = $1", userID).Scan(&c.ID, &c.Email, &c.Phone)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (s *ContactStore) ListUserContacts(ctx context.Context, userIDs []string) ([]UserContact, error) {
	rows, err := s.Pool.QueryContext(ctx,
		`SELECT id, email, phone FROM users WHERE id = ANY($1)`, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]UserContact, 0, len(userIDs))
	for rows.Next() {
		var c UserContact
		if err := rows.Scan(&c.ID, &c.Email, &c.Phone); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (s *ContactStore) BatchGetByIDs(ctx context.Context, ids []int) ([]UserContact, error) {
	rows, err := s.Pool.QueryContext(ctx,
		`SELECT id, email, phone FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []UserContact
	for rows.Next() {
		var c UserContact
		rows.Scan(&c.ID, &c.Email, &c.Phone)
		result = append(result, c)
	}
	return result, nil
}
