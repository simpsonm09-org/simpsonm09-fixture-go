package api

// Problem is an RFC 7807 problem detail returned for item errors.
type Problem struct {
	Type   string `json:"type" example:"about:blank"`
	Title  string `json:"title" example:"Item not found"`
	Status int    `json:"status" example:"404"`
	Detail string `json:"detail" example:"Item 9 was not found"`
}
