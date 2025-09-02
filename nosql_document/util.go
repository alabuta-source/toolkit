package nosql_document

type WriteResponse struct {
	ID      any  `json:"id"`
	Success bool `json:"success"`
}

func newWriteResponse(id any, success bool) *WriteResponse {
	return &WriteResponse{
		ID:      id,
		Success: success,
	}
}
