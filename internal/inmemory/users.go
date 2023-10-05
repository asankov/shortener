package inmemory

import "github.com/asankov/shortener/internal/users"

func (d *DB) GetUser(email, password string) (*users.User, error) {
	user, found := d.users[email]
	if !found {
		return nil, users.ErrUserNotFound
	}
	return user, nil
}

func (d *DB) CreateUser(email, password string, roles []users.Role) error {
	d.users[email] = &users.User{Email: email, Roles: roles}
	return nil
}

func (d *DB) ShouldCreateInitialUser() (bool, error) {
	return true, nil
}
