package mysql

type PageQuery struct {
	Limit  int
	Offset int
}

func (p PageQuery) Normalize(defaultLimit, maxLimit int) PageQuery {
	limit := p.Limit
	offset := p.Offset

	if defaultLimit <= 0 {
		defaultLimit = 20
	}
	if maxLimit <= 0 {
		maxLimit = 200
	}

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return PageQuery{
		Limit:  limit,
		Offset: offset,
	}
}
