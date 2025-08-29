package mysql

type DBChat struct {
	ID         string `db:"id"`
	ExternalID string `db:"external_id"`
}
