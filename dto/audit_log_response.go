package dto

import (
	"Blog/models"
	"time"
)

type AuditLogResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Object    string    `json:"object"`
	ObjectID  uint      `json:"object_id"`
	Timestamp time.Time `json:"timestamp"`

	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	MetaData  string `json:"metadata"`
}

func ToAuditLogResponse(log models.AuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		Object:    log.Object,
		ObjectID:  log.ObjectID,
		Timestamp: log.Timestamp,
		IP:        log.IP,
		UserAgent: log.UserAgent,
		MetaData:  log.Metadata,
	}
}

func ToAuditLogList(logs []models.AuditLog) []AuditLogResponse {
	result := make([]AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, ToAuditLogResponse(log))
	}
	return result
}
