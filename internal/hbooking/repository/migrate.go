package repository

import (
	"context"
	"fmt"
)

func (r *Repository) Migrate() error {
	query := `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS workshop_bookings 
(
    booking_id      UUID DEFAULT uuid_generate_v4() NOT NULL,
    workshop_id     BIGINT                          NOT NULL,
    begin_at        TIMESTAMP                       NOT NULL,
    end_at          TIMESTAMP                       NOT NULL,
    client_id 	 	TEXT                            NOT NULL,
    client_timezone TEXT                            NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_workshop_bookings_booking_id
	ON workshop_bookings (booking_id);

CREATE INDEX IF NOT EXISTS idx_workshop_bookings
	ON workshop_bookings (workshop_id, begin_at, end_at);

CREATE TABLE IF NOT EXISTS workshop_schedules
(
	workshop_id       BIGINT NOT NULL,
	workshop_timezone TEXT   NOT NULL,
	begin_at          TIME   NOT NULL,
	end_at            TIME   NOT NULL
);

CREATE INDEX IF NOT EXISTS workshop_schedule_workshop_id_idx
	ON workshop_schedules (workshop_id);
`

	_, err := r.pool.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	return nil
}
