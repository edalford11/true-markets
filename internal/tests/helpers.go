package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/gin-gonic/gin"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/require"
)

// TestResponse tests response body.
func TestResponse[T any](t *testing.T, w *httptest.ResponseRecorder, expected T, code int) {
	t.Helper()
	require.NotNil(t, w)
	require.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("expected a %d response but got %d: %s", code, w.Code, w.Body.String()),
	)

	actualMap := make(map[string]any)
	err := json.Unmarshal(w.Body.Bytes(), &actualMap)
	require.NoError(t, err, w.Body.String())

	expectedMap := make(map[string]any)
	expectedBytes, err := json.Marshal(expected)
	require.NoError(t, err, w.Body.String())
	err = json.Unmarshal(expectedBytes, &expectedMap)
	require.NoError(t, err, w.Body.String())

	require.Equal(t, len(expectedMap), len(actualMap), "len marshalled fields do not match")

	var actual T
	err = json.Unmarshal(w.Body.Bytes(), &actual)
	require.NoError(t, err, w.Body.String())

	changelog, err := diff.Diff(expected, actual)
	require.NoError(t, err, "failed to diff expected and actual response")

	if len(changelog) > 0 {
		buffer := bytes.NewBuffer(nil)
		w := tabwriter.NewWriter(buffer, 0, 4, 4, '\t', 0)
		fmt.Fprintln(w, "response body did not match expected")
		for _, change := range changelog {
			fmt.Fprintf(w, "%s\nEXPECTED:%#v\nACTUAL:%#v\n\n",
				strings.Join(change.Path, "."),
				change.From,
				change.To)
		}
		err = w.Flush()
		require.NoError(t, err)

		t.Fatal(buffer.String())
	}
}

type TestSuite interface {
	T() *testing.T
}

type ApiTestSuite interface {
	TestSuite
	Router() *gin.Engine
}

func serve(suite ApiTestSuite, method string, endpoint string) *httptest.ResponseRecorder {
	suite.T().Helper()
	request, err := http.NewRequest(method, endpoint, nil)
	require.NoError(suite.T(), err, "request setup should not get an error")
	request.Header.Set("Authorization", "Bearer abc123")
	recorder := httptest.NewRecorder()
	suite.Router().ServeHTTP(recorder, request)
	return recorder
}

func ServeGet(suite ApiTestSuite, endpoint string) *httptest.ResponseRecorder {
	suite.T().Helper()
	return serve(suite, http.MethodGet, endpoint)
}
