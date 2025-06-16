package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"
)

// stubServer simulates a minimal IMDS for testing.
func stubServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/latest/meta-data/", func(w http.ResponseWriter, r *http.Request) {
		// root: one leaf and one folder
		w.Write([]byte("instance-id\nfolder/\n"))
	})
	mux.HandleFunc("/latest/meta-data/instance-id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("i-abcdef123456"))
	})
	mux.HandleFunc("/latest/meta-data/folder/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("subkey\n"))
	})
	mux.HandleFunc("/latest/meta-data/folder/subkey", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("subvalue"))
	})
	return httptest.NewServer(mux)
}

func TestMetadataClient_GetAll(t *testing.T) {
	srv := stubServer()
	defer srv.Close()

	client := NewMetadataClient(srv.URL + "/latest/meta-data/")
	client.httpClient = srv.Client()

	all, err := client.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}

	// Verify "instance-id"
	if id, ok := all["instance-id"].(string); !ok || id != "i-abcdef123456" {
		t.Errorf("expected instance-id 'i-abcdef123456', got %v", all["instance-id"])
	}

	// Verify nested folder
	folder, ok := all["folder"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected folder map, got %T", all["folder"])
	}
	if sub, ok := folder["subkey"].(string); !ok || sub != "subvalue" {
		t.Errorf("expected folder.subkey 'subvalue', got %v", folder["subkey"])
	}
}

func TestMetadataClient_GetKey(t *testing.T) {
	srv := stubServer()
	defer srv.Close()

	client := NewMetadataClient(srv.URL + "/latest/meta-data/")
	client.httpClient = srv.Client()

	// Leaf key
	val, err := client.Get("instance-id")
	if err != nil {
		t.Fatalf("Get(instance-id) error: %v", err)
	}
	if val.(string) != "i-abcdef123456" {
		t.Errorf("expected 'i-abcdef123456', got %v", val)
	}

	// Nested key
	val2, err := client.Get("folder/subkey")
	if err != nil {
		t.Fatalf("Get(folder/subkey) error: %v", err)
	}
	if val2.(string) != "subvalue" {
		t.Errorf("expected 'subvalue', got %v", val2)
	}
}

func TestMetadataClient_GetMissing(t *testing.T) {
	srv := stubServer()
	defer srv.Close()

	client := NewMetadataClient(srv.URL + "/latest/meta-data/")
	client.httpClient = srv.Client()

	_, err := client.Get("nonexistent")
	if err == nil || !strings.Contains(err.Error(), `key "nonexistent" not found`) {
		t.Errorf("expected key not found error, got %v", err)
	}
}

func TestFetchEC2Metadata_NoAWSConfig(t *testing.T) {
	// Unset any AWS environment to force config load error
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")

	ctx := context.Background()
	_, err := fetchEC2Metadata(ctx, "i-unknown")
	if err == nil || !strings.Contains(err.Error(), "failed to load AWS config") {
		t.Errorf("expected AWS config load error, got %v", err)
	}
}
