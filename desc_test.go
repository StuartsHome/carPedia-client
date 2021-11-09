package carpedia

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDesc(t *testing.T) {
	client, mux, _ := setup()
	body := "{null, null}"

	mux.HandleFunc("/desc", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))

		testMethod(t, r, http.MethodGet)
	})

	ctx := context.Background()
	opts := DescOptions{}
	descs, err := client.Desc.Desc(ctx, &opts)
	require.Nil(t, err)

	actual, err := io.ReadAll(descs.Body)
	require.Nil(t, err)

	assert.Equal(t, "{null, null}", string(actual))

}

// TODO
func TestAddDesc(t *testing.T) {
	// client, mux, _ := setup()
	// body := "{null, null}"

	// mux.HandleFunc("/desc", func(w http.ResponseWriter, r *http.Request) {

	// 	w.Header().Add("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(body))

	// 	testMethod(t, r, http.MethodPost)
	// })
	// ctx := context.Background()
	// opts := DescOptions{}

}
