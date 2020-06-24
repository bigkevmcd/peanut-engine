package recent

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
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
		Results: []responseSyncItem{},
	}
	for _, v := range s.Results {
		r.Results = append(r.Results, responseSyncItem{GVK: v.ResourceKey, Message: v.Message, Status: v.Status})
	}

	return r
}

type responseSync struct {
	Start   string             `json:"startTime"`
	End     string             `json:"endTime"`
	SHA     string             `json:"sha"`
	Results []responseSyncItem `json:"results"`
}

type responseSyncItem struct {
	// TODO: break this out in some way otherwise the JSON encoding is weird).
	GVK     kube.ResourceKey  `json:"gvk"`
	Status  common.ResultCode `json:"status"`
	Message string            `json:"message"`
}
