package kkcity

import "github.com/jackc/pgx"

// Use the pool to do further operations.
// langs must follow ISO-639-1 (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
func Use(langs []string, gKey string, pool *pgx.ConnPool) {
	setupLanguage(langs)
	googleKey = gKey
	use(pool, getAll())
}
