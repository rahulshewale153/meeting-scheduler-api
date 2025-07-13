package model

type EventRequest struct {
	Event
	ProposedSlots []EventSlot `json:"proposed_slots" validate:"required,dive,required"`
}

type Event struct {
	ID              int64  `json:"id"`
	Title           string `json:"title" validate:"required"`
	OrganizerID     int64  `json:"organizer_id" validate:"required"`
	DurationMinutes int    `json:"duration_minutes" validate:"required"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

type EventSlot struct {
	ID        int64  `json:"id"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}
