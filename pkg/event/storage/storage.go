package storage

import (
	"context"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/event"
)

const dbName = "event_log"

type (
	Storage struct {
		db db.Database
	}

	Options struct {
		DB db.Database
	}
)

// New builds Storage from Options
func New(opt Options) *Storage {
	return &Storage{
		db: opt.DB,
	}
}

func (s *Storage) Persist(ctx context.Context, ev event.Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	id := ev.Metadata.ID
	user := ev.Tenant.ID
	name := ev.Metadata.Name
	q, _, err := goqu.
		Insert(dbName).
		Cols("eventid", "eventname", "userid", "info").
		Vals(goqu.Vals{id, name, user, data}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, q)
	if err != nil {
		return err
	}
	return nil
}
