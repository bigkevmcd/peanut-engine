package recent

import (
	"container/ring"
	"testing"
	"time"

	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-cmp/cmp"
)

// TODO: What sort of access should we give to the ring?

func TestLatest(t *testing.T) {
	syncs := NewRecentSynchronisations(ring.New(5))
	start, end := time.Now(), time.Now()
	sha := "7f193461f0b44fc5e397a63f2ddba8d9453e7a3f"

	syncs.Add(start, end, plumbing.NewHash(sha), []common.ResourceSyncResult{})

	want := Synchronisation{
		Start:   start,
		End:     end,
		SHA:     sha,
		Results: []common.ResourceSyncResult{},
	}
	if diff := cmp.Diff(want, syncs.Latest()); diff != "" {
		t.Fatalf("latest sync failed:\n%s", diff)
	}
}
