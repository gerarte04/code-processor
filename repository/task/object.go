package task

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
)

type Task struct {
    Id uuid.UUID
    Finished bool
    Result int
}

func SleepAndComplete(t *Task, d time.Duration) {
    time.Sleep(d)
    fmt.Println("task finished with id " + t.Id.String())

    t.Result = rand.IntN(1000)
    t.Finished = true
}
