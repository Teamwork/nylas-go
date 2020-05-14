package nylas

import (
	"context"
	"encoding/base64"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const b64TestPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAACklEQVR4nGMAAQAABQABDQottAAAAABJRU5ErkJggg=="

func testPNG() io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64TestPNG))
}

func TestFile(t *testing.T) {
	accessToken := "accessToken"
	id := "br57kcekhf1hsjq04y8aonkit"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/files/"+id)

		_, _ = w.Write(getFileJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.File(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := File{
		ID:          "8cid1lhd0m7x9k5wjrkpufs1a",
		Object:      "file",
		AccountID:   "43jf3n4e***",
		ContentType: "image/png",
		Filename:    "test.png",
		Size:        24429,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("File: (-got +want):\n%s", diff)
	}
}

func TestUploadFile(t *testing.T) {
	accessToken := "accessToken"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodPost, "/files")

		mp, err := r.MultipartReader()
		if err != nil {
			t.Fatalf("expected multipart: %v", err)
		}

		p, err := mp.NextPart()
		if err != nil && err != io.EOF {
			t.Fatalf("next part: %v", err)
		}
		if p == nil || p.FormName() != "file" {
			t.Fatal("expected single field named \"file\"")
		}

		wantContentType := "image/png"
		if ct := p.Header.Get("Content-Type"); ct != wantContentType {
			t.Errorf("expected Content-Type to be detected, got: %q; want %q",
				ct, wantContentType)
		}

		if _, err := png.Decode(p); err != nil {
			t.Fatalf("part not valid png: %v", err)
		}

		_, _ = w.Write(uploadFileJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.UploadFile(context.Background(), "test.png", testPNG())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := File{
		ID:          "8cid1lhd0m7x9k5wjrkpufs1a",
		Object:      "file",
		AccountID:   "43jf3n4e***",
		ContentType: "image/png",
		Filename:    "test.png",
		Size:        24429,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("File: (-got +want):\n%s", diff)
	}
}

func TestDownloadFile(t *testing.T) {
	accessToken := "accessToken"
	id := "br57kcekhf1hsjq04y8aonkit"
	want := []byte(`body`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/files/"+id+"/download")

		_, _ = w.Write(want)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	file, err := client.DownloadFile(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if diff := cmp.Diff(data, want); diff != "" {
		t.Errorf("File: (-got +want):\n%s", diff)
	}
}

func TestDeleteFile(t *testing.T) {
	accessToken := "accessToken"
	id := "br57kcekhf1hsjq04y8aonkit"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodDelete, "/files/"+id)

	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	err := client.DeleteFile(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

var getFileJSON = []byte(`{
    "account_id": "43jf3n4e***",
    "content_type": "image/png",
    "filename": "test.png",
    "id": "8cid1lhd0m7x9k5wjrkpufs1a",
    "object": "file",
    "size": 24429
}`)

var uploadFileJSON = []byte(`[
    {
        "account_id": "43jf3n4e***",
        "content_type": "image/png",
        "filename": "test.png",
        "id": "8cid1lhd0m7x9k5wjrkpufs1a",
        "object": "file",
        "size": 24429
    }
]`)
