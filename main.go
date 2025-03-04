package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func SleepAndComplete(t *Task, d time.Duration) {
    time.Sleep(d)
    fmt.Println("task finished with id " + t.id.String())

    t.result = rand.IntN(1000)
    t.finished = true
}

func (db *Database) getStatusHandler(w http.ResponseWriter, r *http.Request) {
    taskId := path.Base(r.URL.Path)
    fmt.Println("got get status request on id: " + taskId)

    uuid, err := uuid.Parse(taskId)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    } else if _, ok := db.tasks[uuid]; !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    var result string

    if db.tasks[uuid].finished {
        result = "{\"status\": \"ready\"}"
    } else {
        result = "{\"status\": \"in_progress\"}"
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, result)
}

func (db *Database) getResultHandler(w http.ResponseWriter, r *http.Request) {
    taskId := path.Base(r.URL.Path)
    fmt.Println("got get result request on id: " + taskId)

    uuid, err := uuid.Parse(taskId)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    } else if _, ok := db.tasks[uuid]; !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    
    if db.tasks[uuid].finished {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "{\"result\": \"" + strconv.Itoa(db.tasks[uuid].result) + "\"}")
    } else {
        w.WriteHeader(http.StatusProcessing)
    }
}

func (db *Database) postTaskHandler(w http.ResponseWriter, r *http.Request) {
    taskInfo, _ := io.ReadAll(r.Body)
    stringTime := string(taskInfo)
    fmt.Println("got post request with task time: " + string(taskInfo))

    var taskDur time.Duration

    if len(stringTime) == 0 {
        taskDur = time.Second
    } else {
        var err error
        taskDur, err = time.ParseDuration(stringTime)

        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, err.Error())
            return
        }
    }

    newUuid := uuid.New()
    db.tasks[newUuid] = &Task{
        id: newUuid,
    }

    fmt.Println("created task with id " + newUuid.String())

    go SleepAndComplete(db.tasks[newUuid], taskDur)

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "{\"task_id\": \"" + newUuid.String() + "\"}")
}

func createServer(db *Database, addr string) error {
    r := chi.NewRouter()
    r.Route("/", func(r chi.Router) {
        r.Get("/status/{task_id}", db.getStatusHandler)
        r.Get("/result/{task_id}", db.getResultHandler)
        r.Post("/task", db.postTaskHandler)
    })

    s := &http.Server{
        Addr: addr,
        Handler: r,
    }

    return s.ListenAndServe()
}

func main() {
    db := newDatabase()
    err := createServer(db, ":8000")

    if err != nil {
        _ = fmt.Errorf("%s", "failed to start: " + err.Error())
    }
}
