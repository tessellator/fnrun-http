package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tessellator/fnrun"
)

func Source(ctx context.Context, invoker fnrun.Invoker) error {
	http.HandleFunc("/", makeHandler(invoker))
	return http.ListenAndServe(":8080", nil)
}

func makeHandler(invoker fnrun.Invoker) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		r.Body.Close()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		input := fnrun.Input{Data: data}

		result, err := invoker.Invoke(context.Background(), &input)
		if err != nil {
			if err == fnrun.ErrAvailabilityTimeout {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			os.Stdout.WriteString(err.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(result.Status)
		w.Write(result.Data)
	}
}
