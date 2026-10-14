package appgui

import (
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	services "restclient/internal/services/fs"
	"restclient/internal/services/logger"
	"restclient/internal/types"

	"github.com/go-gui-org/go-gui/gui"
)

const saveDelay = 500 * time.Millisecond

var (
	saveMu     sync.Mutex
	saveTimers = map[string]*time.Timer{}
	writeMu    sync.Mutex
)

func editTab(w *gui.Window, id string, fn func(t *types.RequestTab)) {
	updateTab(w, id, fn)
	scheduleSave(w, id)
}

func scheduleSave(w *gui.Window, id string) {
	saveMu.Lock()
	defer saveMu.Unlock()

	if t, ok := saveTimers[id]; ok {
		t.Stop()
	}
	saveTimers[id] = time.AfterFunc(saveDelay, func() {
		w.QueueCommand(func(w *gui.Window) { persistTab(w, id) })
	})
}

func flushSave(w *gui.Window, id string) {
	saveMu.Lock()
	t, pending := saveTimers[id]
	if pending {
		t.Stop()
	}
	saveMu.Unlock()

	if pending {
		persistTab(w, id)
	}
}

func persistTab(w *gui.Window, id string) {
	saveMu.Lock()
	delete(saveTimers, id)
	saveMu.Unlock()

	a := gui.State[types.App](w)

	var snap types.RequestTab
	found := false
	for _, t := range a.OpenTabs {
		if t.ID == id {
			snap = t
			found = true
			break
		}
	}
	if !found || snap.Path == "" {
		return
	}

	snap.Params = slices.Clone(snap.Params)
	snap.Headers = slices.Clone(snap.Headers)

	if meta, ok := a.RequestLookup[id]; ok {
		meta.Method = snap.Method
		a.RequestLookup[id] = meta
	}

	go func() {
		writeMu.Lock()
		defer writeMu.Unlock()

		err := services.WriteRequestFile(snap.Path, func(r *types.Request) {
			applyTabToRequest(snap, r)
		})
		if err != nil {
			logger.Info("gagal menyimpan request: ", err)
		}
	}()
}

func loadTab(meta types.RequestTab) types.RequestTab {
	if meta.Path == "" {
		return meta
	}

	reqCh := make(chan types.Request, 1)
	errCh := make(chan error, 1)
	services.ReadRequestFile(meta.Path, reqCh, errCh)

	select {
	case req := <-reqCh:
		return tabFromRequest(meta, req)
	case <-errCh:
		return meta
	}
}

func tabFromRequest(meta types.RequestTab, req types.Request) types.RequestTab {
	t := meta
	if req.Request.Method != "" {
		t.Method = req.Request.Method
	}
	t.URL = req.Request.URL.Raw
	t.Params = paramsFromURL(t.URL)
	t.Headers = headersFromStrings(req.Request.Header)
	t.Body = req.Request.Body.Content
	t.AuthType = req.Request.Auth.Type
	if t.AuthType == "" {
		t.AuthType = "none"
	}
	return t
}

func applyTabToRequest(t types.RequestTab, req *types.Request) {
	host, path := hostAndPath(t.URL)

	req.Request.Method = t.Method
	req.Request.URL.Raw = t.URL
	req.Request.URL.Host = host
	req.Request.URL.Path = path
	req.Request.URL.Query = queryFromParams(t.Params)
	req.Request.Header = headersToStrings(t.Headers)
	req.Request.Body.Content = t.Body
	req.Request.Auth.Type = t.AuthType
}

func paramsFromURL(raw string) []types.KeyValue {
	i := strings.IndexByte(raw, '?')
	if i < 0 {
		return []types.KeyValue{}
	}

	out := []types.KeyValue{}
	for _, part := range strings.Split(raw[i+1:], "&") {
		if part == "" {
			continue
		}
		k, v, _ := strings.Cut(part, "=")
		out = append(out, types.KeyValue{Key: k, Value: v})
	}
	return out
}

func urlWithParams(raw string, params []types.KeyValue) string {
	base := raw
	if i := strings.IndexByte(raw, '?'); i >= 0 {
		base = raw[:i]
	}

	parts := make([]string, 0, len(params))
	for _, p := range params {
		if p.Key == "" {
			continue
		}
		if p.Value == "" {
			parts = append(parts, p.Key)
		} else {
			parts = append(parts, p.Key+"="+p.Value)
		}
	}

	if len(parts) == 0 {
		return base
	}
	return base + "?" + strings.Join(parts, "&")
}

func queryFromParams(params []types.KeyValue) []types.QueryRequest {
	out := make([]types.QueryRequest, 0, len(params))
	for _, p := range params {
		if p.Key == "" {
			continue
		}
		out = append(out, types.QueryRequest{Key: p.Key, Value: p.Value})
	}
	return out
}

func headersFromStrings(hs []string) []types.KeyValue {
	out := make([]types.KeyValue, 0, len(hs))
	for _, h := range hs {
		k, v, _ := strings.Cut(h, ":")
		out = append(out, types.KeyValue{
			Key:   strings.TrimSpace(k),
			Value: strings.TrimSpace(v),
		})
	}
	return out
}

func headersToStrings(hs []types.KeyValue) []string {
	out := make([]string, 0, len(hs))
	for _, h := range hs {
		if h.Key == "" {
			continue
		}
		out = append(out, h.Key+": "+h.Value)
	}
	return out
}

func hostAndPath(raw string) ([]string, []string) {
	base := raw
	if i := strings.IndexByte(raw, '?'); i >= 0 {
		base = raw[:i]
	}

	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return []string{}, []string{}
	}

	path := []string{}
	if p := strings.Trim(u.Path, "/"); p != "" {
		path = strings.Split(p, "/")
	}
	return strings.Split(u.Hostname(), "."), path
}
