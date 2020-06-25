package recent

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/julienschmidt/httprouter"
)

// RecentRouter is an HTTP API for accessing recent synchronisations.
type RecentRouter struct {
	*httprouter.Router
	recent *RecentSynchronisations
}

// GePipelines fetches and returns the pipeline body.
func (a *RecentRouter) GetLatest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(makeSynchronisationResponse(a.recent.Latest()))

}

// NewRouter creates and returns a new RecentRouter.
func NewRouter(r *RecentSynchronisations) *RecentRouter {
	api := &RecentRouter{Router: httprouter.New(), recent: r}
	api.HandlerFunc(http.MethodGet, "/latest", api.GetLatest)
	return api
}

func makeSynchronisationResponse(s Synchronisation) responseSync {
	r := responseSync{
		Start:   s.Start.Format(time.RFC3339),
		End:     s.End.Format(time.RFC3339),
		SHA:     s.SHA,
		Error:   s.Error.Error(),
		Results: []responseSyncItem{},
	}
	for _, v := range s.Results {
		r.Results = append(r.Results, makeSyncItem(v))
	}

	return r
}

type responseSync struct {
	Start   string             `json:"startTime"`
	End     string             `json:"endTime"`
	SHA     string             `json:"sha"`
	Error   string             `json:"error"`
	Results []responseSyncItem `json:"results"`
}

type responseSyncItem struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Group     string            `json:"group"`
	Kind      string            `json:"kind"`
	Status    common.ResultCode `json:"status"`
	Message   string            `json:"message"`
}

func makeSyncItem(v common.ResourceSyncResult) responseSyncItem {
	return responseSyncItem{
		Message: v.Message, Status: v.Status,
		Name:      v.ResourceKey.Name,
		Namespace: v.ResourceKey.Namespace,
		Group:     v.ResourceKey.Group,
		Kind:      v.ResourceKey.Kind,
	}
}
