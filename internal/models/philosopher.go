package models

import "time"

type Philosopher struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	DateBorn    string    `db:"date_born"`
	DateDied    string    `db:"date_died"`
	Birthplace  string    `db:"birthplace"`
	Interests   []string  `db:"interests"`
	PortraitURI string    `db:"portrait_uri"`
	Bio         string    `db:"bio"`
	CreatedAt   time.Time `db:"created_at"`
}
