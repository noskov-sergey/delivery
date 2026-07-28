package queries

type GetAllQueriersQuery struct {
	valid bool
}

func NewGetAllQueriersQuery() (GetAllQueriersQuery, error) {
	return GetAllQueriersQuery{
		valid: true,
	}, nil
}

func (g *GetAllQueriersQuery) IsValid() bool { return g.valid }
