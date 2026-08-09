package tags

import (
	"dt2026/httpx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GetTagResponse struct {
	Ok     bool     `json:"ok"`
	Status int      `json:"status"`
	Tags   []string `json:"tags"`
}

func (r GetTagResponse) GetStatus() int {
	return r.Status
}

func Get(recordID int64, pool *pgxpool.Pool) httpx.Response {
	// TODO
	return nil
}
