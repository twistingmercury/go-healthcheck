package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/twistingmercury/go-healthcheck"
)

func testServer(status int, delay bool) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if delay {
				time.Sleep(4 * time.Second)
			}
			w.WriteHeader(status)
			fmt.Fprintln(w, "Hello, client")
		}))
}

func TestDependencyHandlerFunc(t *testing.T) {
	exp := healthcheck.HealthStatusResult{
		Status:          healthcheck.HealthStatusOK,
		Name:            "Test Func",
		RequestDuration: 42,
		Resource:        "None",
	}
	dep := healthcheck.DependencyDescriptor{
		HandlerFunc: func() (hsr healthcheck.HealthStatusResult) {
			return exp
		},
	}
	act := dep.HandlerFunc()

	assert.Equal(t, exp, act)
}

func TestCheckDepsInvokesHandlerFunc(t *testing.T) {
	exp := healthcheck.HealthStatusResult{
		Status:          healthcheck.HealthStatusOK,
		Name:            "Test Func",
		RequestDuration: 42,
		Resource:        "None",
	}

	ts := testServer(200, false)
	defer ts.Close()
	deps := []healthcheck.DependencyDescriptor{
		{Connection: ts.URL, Name: "Test URL 1", Type: "HTTP"},
		{Connection: "", Name: "Test custom", Type: "Custom", HandlerFunc: func() (hsr healthcheck.HealthStatusResult) { return exp }},
	}

	s, r := healthcheck.CheckDeps(deps)
	assert.Equal(t, healthcheck.HealthStatusOK, s)
	assert.Equal(t, 2, len(r))
}

func TestCheckUrlReturnsOK(t *testing.T) {
	ts := testServer(200, false)
	defer ts.Close()

	act := healthcheck.CheckURL(ts.URL)
	assert.Equal(t, healthcheck.HealthStatusOK, act.Status)
}

func TestCheckURLReturnsError(t *testing.T) {
	act := healthcheck.CheckURL("hqpn://wtf.is.this.url???")
	assert.Equal(t, healthcheck.HealthStatusCritical, act.Status)
}

func TestCheckUrlReturnWarning(t *testing.T) {
	ts := testServer(200, true)
	defer ts.Close()

	act := healthcheck.CheckURL(ts.URL)
	assert.Equal(t, healthcheck.HealthStatusWarning, act.Status)
}

func TestCheckUrlReturnCritical(t *testing.T) {
	ts := testServer(500, false)
	defer ts.Close()

	act := healthcheck.CheckURL(ts.URL)
	assert.Equal(t, healthcheck.HealthStatusCritical, act.Status)
}

func TestHandlerReturnCritical(t *testing.T) {
	tOK := testServer(200, false)
	tW1 := testServer(300, false)
	tW2 := testServer(200, true)
	tCrit := testServer(500, false)

	defer func() {
		tOK.Close()
		tW1.Close()
		tW2.Close()
		tCrit.Close()
	}()

	deps := []healthcheck.DependencyDescriptor{
		{Connection: tOK.URL, Name: "Test Good 1", Type: "HTTP"},
		{Connection: tCrit.URL, Name: "Test Critical 2", Type: "HTTP"},
		{Connection: tOK.URL, Name: "Test Good 3", Type: "HTTP"},
		{Connection: tW1.URL, Name: "Test Warn 300: SLOW", Type: "HTTP"},
		{Connection: tW2.URL, Name: "Test Warn 300", Type: "HTTP"},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.GET("/test", healthcheck.Handler(deps...))
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(resp, c.Request)

	data := resp.Body.Bytes()
	assert.Greater(t, len(data), 0)

	var hcr healthcheck.HealthResponse

	err := json.Unmarshal(data, &hcr)
	assert.NoError(t, err)
	assert.Equal(t, healthcheck.HealthStatusCritical, hcr.Status)

	str := string(data)
	exp := hcr.String()

	assert.Equal(t, exp, str)
}

func TestHealthStatusString(t *testing.T) {
	assert.Equal(t, "OK", healthcheck.HealthStatusOK.String())
	assert.Equal(t, "Warning", healthcheck.HealthStatusWarning.String())
	assert.Equal(t, "Critical", healthcheck.HealthStatusCritical.String())
	assert.Equal(t, "HealthStatus(5)", healthcheck.HealthStatus(5).String())
}

func TestHealthStatusParse(t *testing.T) {
	ok, err := healthcheck.ParseHealthStatus("OK")
	assert.NoError(t, err)
	assert.Equal(t, healthcheck.HealthStatusOK, ok)

	wn, err := healthcheck.ParseHealthStatus("Warning")
	assert.NoError(t, err)
	assert.Equal(t, healthcheck.HealthStatusWarning, wn)

	cr, err := healthcheck.ParseHealthStatus("Critical")
	assert.NoError(t, err)
	assert.Equal(t, healthcheck.HealthStatusCritical, cr)

	x, err := healthcheck.ParseHealthStatus("Fatal")
	assert.Error(t, err)
	assert.Equal(t, healthcheck.HealthStatusNotSet, x)
}

func TestHealthStatusUnmarshalText(t *testing.T) {
	var err error
	var hs healthcheck.HealthStatus
	err = hs.UnmarshalText([]byte("OK"))
	assert.NoError(t, err)

	err = hs.UnmarshalText([]byte("Warning"))
	assert.NoError(t, err)

	err = hs.UnmarshalText([]byte("Critical"))
	assert.NoError(t, err)

	err = hs.UnmarshalText([]byte("Fatal"))
	assert.Error(t, err)
}

func TestDependencyDescriptorString(t *testing.T) {
	desc := healthcheck.DependencyDescriptor{
		Connection:  "test",
		HandlerFunc: nil,
		Name:        "test",
		Type:        "test",
	}

	js := desc.String()

	assert.Greater(t, len(js), 0)
}

func TestHealthStatusResultString(t *testing.T) {
	r := healthcheck.HealthStatusResult{
		Status:          healthcheck.HealthStatusOK,
		Name:            "test",
		RequestDuration: 42,
		Resource:        "test",
	}

	js := r.String()

	assert.Greater(t, len(js), 0)
}

func TestDependencyTypeString(t *testing.T) {
	assert.Equal(t, "OK", healthcheck.HealthStatusOK.String())
	assert.Equal(t, "Warning", healthcheck.HealthStatusWarning.String())
	assert.Equal(t, "Critical", healthcheck.HealthStatusCritical.String())
}
