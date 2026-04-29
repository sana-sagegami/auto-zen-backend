package http

// OuraWebhookPayload は Oura から送られる Webhook イベントの構造体。
type OuraWebhookPayload struct {
	EventType  string `json:"event_type"`  // "create" | "update" | "delete"
	DataType   string `json:"data_type"`   // "daily_readiness" | "daily_sleep" など
	Day        string `json:"day"`         // "YYYY-MM-DD"
	DocumentID string `json:"document_id"` // べき等性確認用
}
