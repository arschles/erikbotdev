package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type CounterData struct {
	Name        string
	Value       uint64    `json:"value"`
	TimeUpdated time.Time `json:"time_updated"`
}

var db *bbolt.DB
var USER_BUCKET = []byte("Users")
var FOLLOWER_BUCKET = []byte("Followers")
var COUNTER_BUCKET = []byte("Counters")

func IncrementCounter(counterName string) (current uint64) {
	db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(COUNTER_BUCKET)
		v := b.Get([]byte(counterName))

		oldSchemaErr := json.Unmarshal(v, &current)
		var val CounterData
		// if the data was in the old "schema" (just the counter value),
		// then upgrade it to the new schema
		if oldSchemaErr == nil {
			val.Value = current
			val.TimeUpdated = time.Now()
		} else {
			// otherwise it's in the new schema, so decode to that
			if err := json.Unmarshal(v, &val); err != nil {
				return errors.Wrap(
					err,
					fmt.Sprintf(
						"Unmarshaling counter %s to the new schema",
						counterName,
					),
				)
			}
		}

		val.Name = counterName
		val.Value++
		val.TimeUpdated = time.Now()

		toWrite, err := json.Marshal(&val)
		if err != nil {
			return err
		}

		b.Put([]byte(counterName), toWrite)

		return nil
	})

	return
}

func ListCounters() []CounterData {
	counters := make([]CounterData, 0)

	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(COUNTER_BUCKET)
		b.ForEach(func(k, v []byte) error {
			var cData CounterData
			if err := json.Unmarshal(v, &cData); err != nil {
				return err
			}
			counters = append(counters, cData)
			return nil
		})
		return nil
	})

	return counters
}

func CountersListTimer(
	counterNames []string,
	counterDataChan chan<- []CounterData,
) {
	// this thing needs to go into a goroutine

	ticker := time.NewTicker(5 * time.Second)
	counterNamesSet := map[string]struct{}{}
	for {
		<-ticker.C
		counters := ListCounters()
		retCounters := []CounterData{}
		for _, counter := range counters {
			if _, ok := counterNamesSet[counter.Name]; ok {
				retCounters = append(retCounters, counter)
			}
		}
		counterDataChan <- retCounters
	}
	// start a timer that ticks every 30 mins and iterate the counterNames.
	// print out the duration since each one was written
}

func InitDatabase(file string, mode os.FileMode) error {
	var err error
	db, err = bbolt.Open(file, mode, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Initialize (create any needed buckets, ensure they don't exists first)
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists(USER_BUCKET)
	if err != nil {
		return err
	}

	_, err = tx.CreateBucketIfNotExists(COUNTER_BUCKET)
	if err != nil {
		return err
	}

	_, err = tx.CreateBucketIfNotExists(FOLLOWER_BUCKET)
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		UpdateFollowers()
		t := time.NewTicker(5 * time.Minute)
		for range t.C {
			UpdateFollowers()
		}
	}()

	return nil
}
