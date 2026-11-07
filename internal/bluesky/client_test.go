package bluesky

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type call struct {
	method string
	path   string
	auth   string
	body   []byte
}

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *[]call) {
	t.Helper()
	var calls []call
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))
		calls = append(calls, call{r.Method, r.URL.Path, r.Header.Get("Authorization"), body})
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return New(Config{
		Handle:      "me.bsky.social",
		AppPassword: "app-pass",
		PDS:         srv.URL,
	}), &calls
}

func TestLogin(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/xrpc/com.atproto.server.createSession" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["identifier"] != "me.bsky.social" || body["password"] != "app-pass" {
			t.Errorf("unexpected login body: %v", body)
		}
		_, _ = w.Write([]byte(`{"accessJwt":"jwt-123","did":"did:plc:abc"}`))
	})
	if err := client.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	if client.DID() != "did:plc:abc" {
		t.Fatalf("DID = %q", client.DID())
	}
}

func TestUpdateDescriptionPreservesFields(t *testing.T) {
	client, calls := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			_, _ = w.Write([]byte(`{"accessJwt":"jwt-123","did":"did:plc:abc"}`))
		case "/xrpc/com.atproto.repo.getRecord":
			_, _ = w.Write([]byte(`{"cid":"cid-old","value":{"$type":"app.bsky.actor.profile","displayName":"Me","description":"old","avatar":{"$type":"blob","ref":{"$link":"bafy"}}}}`))
		case "/xrpc/com.atproto.repo.putRecord":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["swapRecord"] != "cid-old" {
				t.Errorf("swapRecord = %v", body["swapRecord"])
			}
			record := body["record"].(map[string]any)
			if record["displayName"] != "Me" {
				t.Errorf("displayName lost: %v", record)
			}
			if record["description"] != "new bio" {
				t.Errorf("description = %v", record["description"])
			}
			if _, ok := record["avatar"]; !ok {
				t.Errorf("avatar lost: %v", record)
			}
			_, _ = w.Write([]byte(`{"uri":"at://x","cid":"cid-new"}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})
	ctx := context.Background()
	if err := client.Login(ctx); err != nil {
		t.Fatal(err)
	}
	changed, err := client.UpdateDescription(ctx, "new bio")
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("changed = false, want true")
	}
	if got := (*calls)[len(*calls)-1].auth; got != "Bearer jwt-123" {
		t.Fatalf("Authorization = %q", got)
	}
}

func TestUpdateDescriptionSkipsUnchanged(t *testing.T) {
	client, calls := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			_, _ = w.Write([]byte(`{"accessJwt":"jwt-123","did":"did:plc:abc"}`))
		case "/xrpc/com.atproto.repo.getRecord":
			_, _ = w.Write([]byte(`{"cid":"cid-old","value":{"$type":"app.bsky.actor.profile","description":"same"}}`))
		case "/xrpc/com.atproto.repo.putRecord":
			t.Error("putRecord must not be called")
		}
	})
	ctx := context.Background()
	if err := client.Login(ctx); err != nil {
		t.Fatal(err)
	}
	changed, err := client.UpdateDescription(ctx, "same")
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("changed = true, want false")
	}
	for _, c := range *calls {
		if strings.HasSuffix(c.path, "putRecord") {
			t.Fatal("putRecord was called")
		}
	}
}

func TestUpdateDescriptionCreatesMissingRecord(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			_, _ = w.Write([]byte(`{"accessJwt":"jwt-123","did":"did:plc:abc"}`))
		case "/xrpc/com.atproto.repo.getRecord":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"RecordNotFound","message":"Could not locate record"}`))
		case "/xrpc/com.atproto.repo.putRecord":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			record := body["record"].(map[string]any)
			if record["$type"] != "app.bsky.actor.profile" {
				t.Errorf("$type = %v", record["$type"])
			}
			if record["description"] != "first bio" {
				t.Errorf("description = %v", record["description"])
			}
			if swap, ok := body["swapRecord"]; ok && swap != "" {
				t.Errorf("swapRecord = %v, want omitted", swap)
			}
			_, _ = w.Write([]byte(`{"uri":"at://x","cid":"cid-new"}`))
		}
	})
	ctx := context.Background()
	if err := client.Login(ctx); err != nil {
		t.Fatal(err)
	}
	changed, err := client.UpdateDescription(ctx, "first bio")
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("changed = false, want true")
	}
}

func TestErrorsDoNotLeakSecrets(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"AuthenticationRequired","message":"bad token"}`))
	})
	err := client.Login(context.Background())
	if err == nil {
		t.Fatal("expected an error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "401") || !strings.Contains(msg, "AuthenticationRequired") {
		t.Fatalf("error lacks status or code: %q", msg)
	}
	if strings.Contains(msg, "app-pass") || strings.Contains(msg, "jwt-123") {
		t.Fatalf("error leaks a secret: %q", msg)
	}
}

func TestPDSNormalization(t *testing.T) {
	c := New(Config{Handle: "h", AppPassword: "p", PDS: "bsky.social/"})
	if c.pds != "https://bsky.social" {
		t.Fatalf("pds = %q", c.pds)
	}
}
