package recent

import (
	"container/ring"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-cmp/cmp"
)

func TestGetLatest(t *testing.T) {
	ts, s := makeServer(t)
	start, end := time.Date(2020, time.June, 24, 22, 0, 0, 0, time.UTC), time.Date(2020, time.June, 24, 22, 1, 0, 0, time.UTC)
	sha := "7f193461f0b44fc5e397a63f2ddba8d9453e7a3f"
	testErr := errors.New("this is an error")
	s.Add(start, end, plumbing.NewHash(sha), testErr, []common.ResourceSyncResult{
		{
			Status:  common.ResultCodeSyncFailed,
			Message: "service/taxi failed",
			ResourceKey: kube.ResourceKey{
				Group:     "v1",
				Kind:      "ConfigMap",
				Namespace: "test",
				Name:      "test-cfg",
			},
		},
	})

	req := makeClientRequest(t, fmt.Sprintf("%s/latest", ts.URL))
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assertJSONResponse(t, res, map[string]interface{}{
		"startTime": "2020-06-24T22:00:00Z",
		"endTime":   "2020-06-24T22:01:00Z",
		"sha":       sha,
		"error":     testErr.Error(),
		"results": []interface{}{
			map[string]interface{}{
				"group":     "v1",
				"kind":      "ConfigMap",
				"name":      "test-cfg",
				"namespace": "test",
				"message":   "service/taxi failed",
				"status":    "SyncFailed",
			},
		},
	})
}

func makeClientRequest(t *testing.T, path string) *http.Request {
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func makeServer(t *testing.T) (*httptest.Server, *RecentSynchronisations) {
	syncs := NewRecentSynchronisations(ring.New(5))
	router := NewRouter(syncs)
	ts := httptest.NewTLSServer(router)
	t.Cleanup(ts.Close)
	return ts, syncs
}

// TODO: assert the content-type.
func assertJSONResponse(t *testing.T, res *http.Response, want map[string]interface{}) {
	t.Helper()
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		errMsg, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("didn't get a successful response: %v (%s)", res.StatusCode, strings.TrimSpace(string(errMsg)))
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]interface{}{}
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("failed to parse %s: %s", b, err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("JSON response failed:\n%s", diff)
	}
}
