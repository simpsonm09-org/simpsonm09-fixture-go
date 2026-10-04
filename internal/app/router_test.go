package app_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/simpsonm09-org/simpsonm09-fixture-go/internal/app"
)

type apiClient struct {
	t      *testing.T
	base   string
	client *http.Client
}

func newClient(t *testing.T) *apiClient {
	t.Helper()
	server := httptest.NewServer(app.NewRouter())
	t.Cleanup(server.Close)
	return &apiClient{t: t, base: server.URL, client: server.Client()}
}

// do sends a request and returns the status, headers, and body.
func (c *apiClient) do(method, path, body string) (int, http.Header, []byte) {
	c.t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	request, err := http.NewRequest(method, c.base+path, reader)
	if err != nil {
		c.t.Fatalf("build %s %s request: %v", method, path, err)
	}
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := c.client.Do(request)
	if err != nil {
		c.t.Fatalf("%s %s: %v", method, path, err)
	}
	defer func() { _ = response.Body.Close() }()
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		c.t.Fatalf("read %s %s body: %v", method, path, err)
	}
	return response.StatusCode, response.Header, payload
}

func (c *apiClient) decode(payload []byte) map[string]any {
	c.t.Helper()
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		c.t.Fatalf("decode %s: %v", payload, err)
	}
	return decoded
}

type item struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func TestListReturnsTheSeeds(t *testing.T) {
	client := newClient(t)

	status, _, body := client.do(http.MethodGet, "/items", "")

	if status != http.StatusOK {
		t.Fatalf("GET /items status = %d, want 200", status)
	}
	var items []item
	if err := json.Unmarshal(body, &items); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("GET /items returned %d items, want 3", len(items))
	}
	if items[0].Name != "Widget" || items[1].Name != "Gadget" || items[2].Name != "Gizmo" {
		t.Fatalf("seed names = %q, %q, %q", items[0].Name, items[1].Name, items[2].Name)
	}
}

func TestCreateReadUpdateDeleteLifecycle(t *testing.T) {
	client := newClient(t)

	status, _, body := client.do(http.MethodPost, "/items", `{"name":"Widget","description":"A small widget"}`)
	if status != http.StatusCreated {
		t.Fatalf("POST /items status = %d, want 201", status)
	}
	created := client.decode(body)
	id := int64(created["id"].(float64))
	if id <= 0 {
		t.Fatalf("created id = %d, want positive", id)
	}
	if created["name"] != "Widget" || created["description"] != "A small widget" {
		t.Fatalf("created = %v", created)
	}

	status, _, body = client.do(http.MethodGet, "/items/"+itoa(id), "")
	if status != http.StatusOK {
		t.Fatalf("GET /items/%d status = %d, want 200", id, status)
	}
	one := client.decode(body)
	if one["id"].(float64) != float64(id) || one["name"] != "Widget" {
		t.Fatalf("GET /items/%d = %v", id, one)
	}

	status, _, body = client.do(http.MethodPut, "/items/"+itoa(id), `{"name":"Renamed","description":"Still here"}`)
	if status != http.StatusOK {
		t.Fatalf("PUT /items/%d status = %d, want 200", id, status)
	}
	updated := client.decode(body)
	if updated["name"] != "Renamed" || updated["description"] != "Still here" {
		t.Fatalf("PUT /items/%d = %v", id, updated)
	}

	status, _, _ = client.do(http.MethodDelete, "/items/"+itoa(id), "")
	if status != http.StatusNoContent {
		t.Fatalf("DELETE /items/%d status = %d, want 204", id, status)
	}

	status, _, _ = client.do(http.MethodGet, "/items/"+itoa(id), "")
	if status != http.StatusNotFound {
		t.Fatalf("GET deleted item status = %d, want 404", status)
	}
}

func TestMissingDescriptionIsStoredAsNull(t *testing.T) {
	client := newClient(t)

	status, _, body := client.do(http.MethodPost, "/items", `{"name":"No description"}`)

	if status != http.StatusCreated {
		t.Fatalf("POST /items status = %d, want 201", status)
	}
	if value, present := client.decode(body)["description"]; !present || value != nil {
		t.Fatalf("description = %v, want null", value)
	}
}

func TestUnknownIDReturnsProblemDetail(t *testing.T) {
	client := newClient(t)

	cases := []struct{ method, body string }{
		{http.MethodGet, ""},
		{http.MethodPut, `{"name":"X"}`},
		{http.MethodDelete, ""},
	}
	for _, testCase := range cases {
		status, header, body := client.do(testCase.method, "/items/999", testCase.body)
		if status != http.StatusNotFound {
			t.Fatalf("%s /items/999 status = %d, want 404", testCase.method, status)
		}
		if contentType := header.Get("Content-Type"); !strings.Contains(contentType, "application/problem+json") {
			t.Fatalf("%s /items/999 Content-Type = %q, want application/problem+json", testCase.method, contentType)
		}
		problem := client.decode(body)
		if problem["type"] != "about:blank" || problem["title"] != "Item not found" {
			t.Fatalf("%s /items/999 problem = %v", testCase.method, problem)
		}
		if problem["status"].(float64) != 404 || problem["detail"] != "Item 999 was not found" {
			t.Fatalf("%s /items/999 problem = %v", testCase.method, problem)
		}
	}
}

func TestInvalidBodiesReturnBadRequest(t *testing.T) {
	client := newClient(t)

	cases := []struct {
		name string
		body string
	}{
		{"empty name", `{"name":""}`},
		{"blank name", `{"name":"   "}`},
		{"over-length name", `{"name":"` + strings.Repeat("a", 201) + `"}`},
		{"over-length description", `{"name":"Ok","description":"` + strings.Repeat("a", 2001) + `"}`},
		{"malformed json", `{"name": "Broken"`},
	}
	for _, testCase := range cases {
		status, header, body := client.do(http.MethodPost, "/items", testCase.body)
		if status != http.StatusBadRequest {
			t.Fatalf("POST /items with %s status = %d, want 400", testCase.name, status)
		}
		if contentType := header.Get("Content-Type"); !strings.Contains(contentType, "application/problem+json") {
			t.Fatalf("POST /items with %s Content-Type = %q, want application/problem+json", testCase.name, contentType)
		}
		problem := client.decode(body)
		if problem["type"] != "about:blank" || problem["title"] != "Validation failed" {
			t.Fatalf("POST /items with %s problem = %v", testCase.name, problem)
		}
		if problem["status"].(float64) != 400 {
			t.Fatalf("POST /items with %s problem status = %v, want 400", testCase.name, problem["status"])
		}
	}
}

func itoa(id int64) string {
	return strconv.FormatInt(id, 10)
}
