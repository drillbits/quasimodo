package quasimodo

import (
  "errors"
  "log"
  "os/exec"
  "sync"
  "time"
  "encoding/json"
)

var (
  taskStore *TaskStore
  procStore *ProcStore
)

const timefmt = "2006/01/02 15:04"

type TaskStore struct {
  Tasks []*Task `json:"tasks"`
  mu    *sync.Mutex
}

type Task struct {
  From    Time   `json:"from"`
  To      Time   `json:"to"`
  Command string `json:"command"`
}

type ProcStore struct {
  ProcMap map[int]*Proc
  mu      *sync.Mutex
}

type Proc struct {
  From time.Time
  To   time.Time
  Cmd  *exec.Cmd
}

type Time struct {
  time.Time
  fmt string
}

func NewTime(t time.Time) Time {
  return Time{t, timefmt}
}

func (t Time) format() string {
  return t.Time.Format(t.fmt)
}

func (t Time) MarshalText() ([]byte, error) {
  return []byte(t.format()), nil
}

func (t Time) MarshalJSON() ([]byte, error) {
  return []byte(`"`+t.format()+`"`), nil
}

func NewStore(conf *Config) {
  if taskStore != nil || procStore != nil {
    panic("quasimodo: NewStore must be called only once")
  }
  var tasks []*Task
  if conf.MaxTasks > 0 {
    tasks = make([]*Task, 0, conf.MaxTasks)
  } else {
    tasks = make([]*Task, 0)
  }
  taskStore = &TaskStore{
    Tasks: tasks,
    mu:    new(sync.Mutex),
  }
  procStore = &ProcStore{
    ProcMap: make(map[int]*Proc),
    mu:      new(sync.Mutex),
  }
}

func (s *TaskStore) Add(cmd string, from time.Time, to time.Time) {
  s.mu.Lock()
  defer s.mu.Unlock()
  t := &Task{
    From:    NewTime(from),
    To:      NewTime(to),
    Command: cmd,
  }
  s.Tasks = append(s.Tasks, t)
}

func (s *TaskStore) Pop(prep time.Time) []*Task {
  s.mu.Lock()
  defer s.mu.Unlock()
  log.Printf("task pop time: %v\n", prep.Format(timefmt))
  popTasks := make([]*Task, 0, len(s.Tasks))
  remTasks := make([]*Task, 0, len(s.Tasks))
  for _, t := range s.Tasks {
    log.Printf("task.From: %v\n", t.From.Time)
    if prep.After(t.From.Time) {
      popTasks = append(popTasks, t)
    } else {
      remTasks = append(remTasks, t)
    }
  }
  s.Tasks = remTasks
  return popTasks
}

func (s *TaskStore) Marshal() ([]byte, error) {
  s.mu.Lock()
  defer s.mu.Unlock()
  return json.Marshal(s)
}

func (p *Proc) Pid() (int, error) {
  if p.Cmd.Process == nil {
    return 0, errors.New("process does not started")
  }
  return p.Cmd.Process.Pid, nil
}

func (p *Proc) Kill() (int, error) {
  pid, err := p.Pid()
  if err != nil {
    return 0, err
  }
  err = p.Cmd.Process.Kill()
  return pid, err
}

func (s *ProcStore) Add(p *Proc) (int, error) {
  s.mu.Lock()
  defer s.mu.Unlock()
  pid, err := p.Pid()
  if err != nil {
    return 0, err
  }
  s.ProcMap[pid] = p
  return pid, nil
}

func (s *ProcStore) Pop(now time.Time) []*Proc {
  s.mu.Lock()
  defer s.mu.Unlock()
  log.Printf("proc pop time: %v\n", now.Format(timefmt))
  popProcs := make([]*Proc, 0, len(s.ProcMap))
  for pid, p := range s.ProcMap {
    log.Printf("proc.To: %v\n", p.To)
    if now.After(p.To) {
      popProcs = append(popProcs, p)
      delete(s.ProcMap, pid)
    }
  }
  return popProcs
}

func (s *ProcStore) Marshal() ([]byte, error) {
  s.mu.Lock()
  defer s.mu.Unlock()
  return json.Marshal(s)
}
