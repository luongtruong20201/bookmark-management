package csv

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRecord struct {
	Name string `csv:"name"`
	Age  int    `csv:"age"`
}

func newMultipartFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, fh, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("failed to retrieve file header: %v", err)
	}

	return fh
}

func TestParseFromMultipartFile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		header      *multipart.FileHeader
		checkResult func(t *testing.T, records []*testRecord, err error)
	}{
		{
			name: "success - valid CSV",
			header: newMultipartFileHeader(t, "people.csv", []byte(
				"name,age\nAlice,30\nBob,25\n",
			)),
			checkResult: func(t *testing.T, records []*testRecord, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(records) != 2 {
					t.Fatalf("expected 2 records, got %d", len(records))
				}
				if records[0].Name != "Alice" || records[0].Age != 30 {
					t.Fatalf("unexpected first record: %+v", records[0])
				}
				if records[1].Name != "Bob" || records[1].Age != 25 {
					t.Fatalf("unexpected second record: %+v", records[1])
				}
			},
		},
		{
			name: "error - invalid CSV content (bad integer)",
			header: newMultipartFileHeader(t, "invalid.csv", []byte(
				"name,age\nAlice,not_an_int\n",
			)),
			checkResult: func(t *testing.T, _ []*testRecord, err error) {
				if err == nil {
					t.Fatalf("expected error for invalid CSV, got nil")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var records []*testRecord
			err := ParseFromMultipartFile(tc.header, &records)
			tc.checkResult(t, records, err)
		})
	}
}
