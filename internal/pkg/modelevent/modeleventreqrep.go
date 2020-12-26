package modelevent

import "time"

type ZoneDataRequestEvent struct {
	Zone        string    `json:"zone"`
	OnlySOA     bool      `json:"only_soa"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
}

type ZoneDataReplyEvent struct {
	Zone             string    `json:"zone"`
	ZoneData         string    `json:"zone_data"`
	DnssecEnabled    bool      `json:"dnssec_enabled"`
	PresignedRecords bool      `json:"presigned_records"`
	RepliedAt        time.Time `json:"replied_at,omitempty"`
}
