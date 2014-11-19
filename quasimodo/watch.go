package quasimodo

import (
  "log"
  "os/exec"
  "time"
)

var watcher *Watcher

type Watcher struct {
  tasks       chan *[]Task
  interval    time.Duration
  prep        time.Duration
  startOffset time.Duration
}

func NewWatcher(conf *Config) *Watcher {
  if watcher != nil {
    panic("quasimodo: NewWatch must be called only once")
  }
  watcher = &Watcher{
    tasks:       make(chan *[]Task),
    interval:    toSec(conf.WatchIntervalSeconds),
    prep:        toSec(conf.PrepareSeconds),
    startOffset: toSec(conf.StartOffsetSeconds),
  }
  return watcher
}

func toSec(sec int64) time.Duration {
  return time.Duration(sec) * time.Second
}

func (w *Watcher) Watch() {
  w.watch()
  t := time.NewTicker(watcher.interval)
  for {
    select {
    case <-t.C:
      w.watch()
    }
  }
  t.Stop()
}

func (w *Watcher) watch() {
  log.Println("=====")
  log.Println("watch")
  log.Println("=====")
  tj, _ := taskStore.Marshal()
  log.Printf("taskStore: %s\n", string(tj))
  tasks := taskStore.Pop(time.Now().Add(w.prep))
  log.Printf("poped Tasks: %v\n", len(tasks))
  for _, t := range tasks {
    go w.execute(t)
  }
  pj, _ := procStore.Marshal()
  log.Printf("procStore: %s\n", string(pj))
  procs := procStore.Pop(time.Now())
  log.Printf("poped Procs: %v\n", len(procs))
  for _, p := range procs {
    go w.kill(p)
  }
}

func (w *Watcher) execute(t *Task) {
  d := t.From.Sub(time.Now()) - w.startOffset
  c := make(chan *Proc, 1)
  go func() {
  loop:
    for {
      select {
      case p := <-c:
        log.Println("receive proc")
        pid, err := procStore.Add(p)
        if err != nil {
          log.Printf("proc add error: %v\n", err)
        }
        log.Printf("process started: Pid(%v)\n", pid)
        break loop
      }
    }
  }()
  log.Println("wait to execute")
  time.AfterFunc(d, func() {
    cmd := exec.Command("sh", "-c", t.Command)
    err := cmd.Start()
    if err != nil {
      log.Printf("command error: %v\n", err)
      return
    }
    c <- &Proc{
      From: t.From.Time,
      To:   t.To.Time,
      Cmd:  cmd,
    }
    log.Printf("command executed: %s\n", t.Command)
  })
}

func (w *Watcher) kill(p *Proc) {
  pid, err := p.Kill()
  if err != nil {
    log.Printf("process kill error: %v\n", err)
  } else {
    log.Printf("process killed: Pid(%v)\n", pid)
  }
}
