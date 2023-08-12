package inmemory

import "github.com/asankov/shortener/internal/users"

func (d *DB) Get(email, password string) (*users.User, error) {
	user, found := d.users[email]
	if !found {
		return nil, users.ErrUserNotFound
	}
	return user, nil
}
