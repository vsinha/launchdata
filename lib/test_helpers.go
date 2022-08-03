package lib

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	golden "github.com/jimeh/go-golden"
	"github.com/stretchr/testify/require"
)

func timeParse(timeString string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeString)
	return t
}

func verify(t *testing.T, got interface{}) {
	transformJSON := cmp.FilterValues(func(x, y []byte) bool {
		return json.Valid(x) && json.Valid(y)
	}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
		if err := json.Unmarshal(in, &out); err != nil {
			panic(err) // should never occur given previous filter to ensure valid JSON
		}
		return out
	}))

	gotJson, err := json.MarshalIndent(&got, "", "  ")
	require.NoError(t, err)

	if golden.Update() {
		golden.Set(t, gotJson)
	}

	var want []byte
	want = golden.Get(t)

	if diff := cmp.Diff(gotJson, want, transformJSON); diff != "" {
		t.Errorf("diff (-got,+want:\n%s", diff)
	}
}
