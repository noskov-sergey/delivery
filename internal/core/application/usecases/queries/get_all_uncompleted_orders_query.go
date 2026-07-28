package queries

type GetAllUncompletedOrdersQuery struct {
	valid bool
}

func NewGetAllUncompletedOrdersQuery() (GetAllUncompletedOrdersQuery, error) {
	return GetAllUncompletedOrdersQuery{
		valid: true,
	}, nil
}

func (g *GetAllUncompletedOrdersQuery) IsValid() bool { return g.valid }
