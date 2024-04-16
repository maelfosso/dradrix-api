package storage

type Queries struct {
	*DBCollections
}

func (q *Queries) WithTx() *Queries {
	return &Queries{}
}

func NewQueries(d Database) *Queries {
	return &Queries{
		DBCollections: d.GetAllCollections(),
	}
}
