package modelevent

import "time"

type ZoneAddEvent struct {
	Zone      string    `json:"zone"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type ZoneChangeEvent struct {
	Zone      string    `json:"zone"`
	ChangedAt time.Time `json:"changed_at,omitempty"`
}

type ZoneDeleteEvent struct {
	Zone      string    `json:"zone"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}
