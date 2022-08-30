package doljara

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	RedisKeySuffixRaw string = "_raw"
	RedisKeySuffixRev string = "_rev"
)

func DoljaraRooms(response http.ResponseWriter, req *http.Request) {

	response.Header().Set("access-control-allow-origin", "*")
	response.Header().Set("content-type", "application/json; charset=utf-8")

	rid := req.URL.Query().Get("rid")
	if rid == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	opt, _ := redis.ParseURL("rediss://:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:xxxxx")
	opt.MaxRetries = -1
	client := redis.NewClient(opt)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		response.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	switch req.Method {
	case http.MethodGet:

		val, err := client.Get(ctx, rid+RedisKeySuffixRaw).Result()
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}

		// response.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(response, val)

	case http.MethodPut:

		origRev, err := strconv.Atoi(req.URL.Query().Get("origrev"))
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		val := Response{}
		val.Room.Code = rid
		if err := json.NewDecoder(req.Body).Decode(&val.Room.Raw); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = client.Get(ctx, rid+RedisKeySuffixRev).Err(); err != nil {
			client.Set(ctx, rid+RedisKeySuffixRev, origRev, 0)
		}

		err = client.Watch(ctx, func(tx *redis.Tx) error {

			currRev, _ := tx.Get(ctx, rid+RedisKeySuffixRev).Int()
			if currRev != origRev {
				return errors.New("")
			}

			currRev++
			val.Room.Raw.Rev = currRev
			json, _ := json.Marshal(val)

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, rid+RedisKeySuffixRev, currRev, 0)
				pipe.Set(ctx, rid+RedisKeySuffixRaw, json, 0)
				return nil
			})
			return err

		}, rid+RedisKeySuffixRev)
		if err != nil {
			response.WriteHeader(http.StatusConflict)
			return
		}

		response.WriteHeader(http.StatusOK)

	case http.MethodOptions:
		response.Header().Set("access-control-allow-headers", "content-type")
		response.Header().Set("access-control-allow-methods", "GET, PUT, OPTIONS")
		response.WriteHeader(http.StatusNoContent)

	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}

}
