package models

import "time"

type Philosopher struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	DateBorn    string    `db:"date_born" json:"date_born"`
	DateDied    string    `db:"date_died" json:"date_died"`
	Birthplace  string    `db:"birthplace" json:"birthplace"`
	Interests   []string  `db:"interests" json:"interests"`
	PortraitURI *string   `db:"portrait_uri" json:"portrait_uri"`
	Bio         string    `db:"bio" json:"bio"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
