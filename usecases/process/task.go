package process

import (
	"fmt"
	"http_server/repository/models"
	"math/rand/v2"
	"time"
)

func SleepAndComplete(t *models.Task, d time.Duration) {
    fmt.Println("task started (id: " + t.Id.String() + ")")
    time.Sleep(d)
    fmt.Println("task finished (id: " + t.Id.String() + ")")

    t.Result = rand.IntN(1000)
    t.Finished = true
}
