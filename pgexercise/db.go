package pgexercise

import (
	"context"
	"database/sql"
)

type Queries struct {
	GetMemberWithRecommender *sql.Stmt
}

func (q *Queries) Prepare(ctx context.Context, db *sql.DB) error {
	var err error
	q.GetMemberWithRecommender, err = db.PrepareContext(ctx, `
		SELECT 
			m.firstname || ' ' ||  m.surname AS member, 
			r.firstname || ' ' || r.surname AS recommnededby 
		FROM cd.members m 
		LEFT JOIN cd.members r ON m.recommendedby = r.memid;
	`)
	if err != nil {
		return err
	}
	return nil
}

type MemberRecommenderPair struct {
	MemberName      string
	RecommenderName string
}

func (q *Queries) GetAllMembersWithRecommender(ctx context.Context) ([]MemberRecommenderPair, error) {
	rows, err := q.GetMemberWithRecommender.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v []MemberRecommenderPair
	for rows.Next() {
		var mem, rec sql.NullString
		err := rows.Scan(&mem, &rec)
		if err != nil {
			return nil, err
		}

		v = append(v, MemberRecommenderPair{mem.String, rec.String})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		GetMemberWithRecommender: tx.Stmt(q.GetMemberWithRecommender),
	}
}
