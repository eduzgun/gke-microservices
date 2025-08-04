package models

type Philosopher struct {
	Name        string   `db:"name"`
	DateBorn    string   `db:"date_born"`
	DateDied    string   `db:"date_died"`
	Birthplace  string   `db:"birthplace"`
	Interests   []string `db:"interests"`
	PortraitURI string   `db:"portrait_uri"`
}
