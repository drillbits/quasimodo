package quasimodo

import (
  "net/http"
  "strings"
  "time"
)

type HTTPHandler struct {
  allowedHosts []string
  basePath     string
}

func NewHTTPHandler(conf *Config) *HTTPHandler {
  h := &HTTPHandler{
    allowedHosts: conf.AllowedHosts,
    basePath:     conf.BasePath,
  }
  http.Handle(conf.BasePath, h)
  return h
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if !strings.HasPrefix(r.URL.Path, h.basePath) {
    http.Error(w, "HTTPHandler serving unexpected path: "+r.URL.Path, http.StatusForbidden)
    return
  }
  if !h.hostIsAllowed(strings.Split(r.RemoteAddr, ":")[0]) {
    http.Error(w, "forbidden host: "+r.Host, http.StatusForbidden)
    return
  }

  switch r.Method {
  case "GET":
    h.getCommandList(w, r)
  case "POST":
    h.putCommand(w, r)
  default:
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
  }
}

func (h *HTTPHandler) hostIsAllowed(host string) bool {
  for _, allowedHost := range h.allowedHosts {
    if host == allowedHost {
      return true
    }
  }
  return false
}

func (h *HTTPHandler) getCommandList(w http.ResponseWriter, _ *http.Request) {
  if taskStore == nil {
    http.Error(w, "taskStore must be initialized", http.StatusInternalServerError)
    return
  }

  j, err := taskStore.Marshal()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(j)
}

func (h *HTTPHandler) putCommand(w http.ResponseWriter, r *http.Request) {
  if taskStore == nil {
    http.Error(w, "taskStore must be initialized", http.StatusInternalServerError)
    return
  }

  err := r.ParseForm()
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  cmd := r.FormValue("command")
  f, err := parseTime(r.FormValue("from"))
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  t, err := parseTime(r.FormValue("to"))
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  taskStore.Add(cmd, f, t)

  h.getCommandList(w, r)
}

func parseTime(s string) (time.Time, error) {
  return time.ParseInLocation(timefmt, s, time.Local)
}
