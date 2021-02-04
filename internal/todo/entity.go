package todo

import (
	"database/sql"
	"time"
)

type Todo struct {
	id          int
	title       string
	desc        string
	done        bool
	createdAt   time.Time
	updatedAt   time.Time
	completedAt sql.NullTime
}

type TodoUpdate struct {
	title string
	desc  string
	done  bool
}
