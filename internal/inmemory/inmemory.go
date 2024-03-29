package inmemory

import (
	"errors"

	"github.com/asankov/shortener/internal/links"
	"github.com/asankov/shortener/internal/random"
	"github.com/asankov/shortener/internal/users"
)

type DB struct {
	links  map[string]*links.Link
	users  map[string]*users.User
	random *random.Random
}

func NewDB() *DB {
	return &DB{
		links: make(map[string]*links.Link),
		users: map[string]*users.User{
			"admin@asankov.dev": {
				Email: "admin@asankov.dev",
				Roles: []users.Role{users.RoleAdmin},
			},
		},
		random: random.New(),
	}
}

func (d *DB) GetByID(id string) (*links.Link, error) {
	link, found := d.links[id]
	if !found {
		return nil, links.ErrLinkNotFound
	}
	return link, nil
}

func (d *DB) GetAll() ([]*links.Link, error) {
	all := make([]*links.Link, 0)
	for _, link := range d.links {
		all = append(all, link)
	}
	return all, nil
}

func (d *DB) Create(id string, url string) error {
	d.links[id] = &links.Link{ID: id, URL: url, Metrics: &links.Metrics{Clicks: 0}}
	return nil
}

func (d *DB) Delete(id string) error {
	delete(d.links, id)
	return nil
}

func (d *DB) IncrementClicks(id string) error {
	link, ok := d.links[id]
	if !ok {
		return links.ErrLinkNotFound
	}
	link.Metrics.Clicks++
	return nil
}

func (d *DB) GenerateID() (string, error) {
	var (
		conflictCount        int
		allowedConflictCount int = 4
		idLength             int = 3

		maxAllowedConflicts = 50
	)
	for {
		if conflictCount > allowedConflictCount {
			idLength++
			allowedConflictCount *= 2
		}

		if conflictCount > maxAllowedConflicts {
			return "", links.ErrIDNotGenerated
		}

		id := d.random.ID(idLength)
		_, err := d.GetByID(id)

		// An item with this ID is not found, so we can safely use it.
		if err != nil && errors.Is(err, links.ErrLinkNotFound) {
			return id, nil
		}
		conflictCount++
	}
}
