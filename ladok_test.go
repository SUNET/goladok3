package goladok3

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SUNET/goladok3/ladokmocks"
	"github.com/SUNET/goladok3/ladoktypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockGenericEndpointServer(t *testing.T, mux *http.ServeMux, contentType, method, url string, reply []byte, statusCode int) {
	mux.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(statusCode)
			testMethod(t, r, method)
			testURL(t, r, url)
			w.Write(reply)
		},
	)
}

func mockNewClient(t *testing.T, env, url string) *Client {
	certPEM, cert, privateKeyPEM, _ := ladokmocks.MockCertificateAndKey(t, env, 0, 100)
	cfg := X509Config{
		URL:            url,
		Certificate:    cert,
		CertificatePEM: certPEM,
		PrivateKeyPEM:  privateKeyPEM,
	}
	client, err := NewX509(cfg)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return client
}

func mockSetup(t *testing.T, env string) (*http.ServeMux, *httptest.Server, *Client) {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	client := mockNewClient(t, env, server.URL)

	return mux, server, client
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func testURL(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.RequestURI)
}

func testBody(t *testing.T, r *http.Request, want string) {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(r.Body)
	assert.NoError(t, err)

	got := buffer.String()
	require.JSONEq(t, want, got)
}

func TestNewX509Validation(t *testing.T) {
	tts := []struct {
		name    string
		config  X509Config
		wantErr bool
	}{
		{
			name:    "empty config",
			config:  X509Config{},
			wantErr: true,
		},
		{
			name: "missing certificate",
			config: X509Config{
				URL:           "https://example.com",
				PrivateKeyPEM: []byte("key"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewX509(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewOIDC(t *testing.T) {
	client, err := NewOIDC(OidcConfig{})
	assert.NoError(t, err)
	assert.Nil(t, client)
}

func TestCheckResponse(t *testing.T) {
	tts := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{"200 OK", 200, nil},
		{"201 Created", 201, nil},
		{"202 Accepted", 202, nil},
		{"204 No Content", 204, nil},
		{"304 Not Modified", 304, nil},
		{"401 Unauthorized", 401, ErrNotAllowedRequest},
		{"500 Internal Server Error", 500, ErrInvalidRequest},
		{"403 Forbidden", 403, ErrInvalidRequest},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			err := checkResponse(resp)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewRequestWithBody(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	body := map[string]string{"key": "value"}

	req, err := client.newRequest(context.Background(), "application/json", "POST", "/test", body)
	assert.NoError(t, err)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, "application/json", req.Header.Get("Accept"))

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(req.Body)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"data"`)
}

func TestNewRequestInvalidURL(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "://bad-url")
	_, err := client.newRequest(context.Background(), "application/json", "GET", "/test", nil)
	assert.Error(t, err)
}

func TestDoInvalidContentType(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, server.URL)
	req, err := client.newRequest(context.Background(), "text/plain", "GET", "/test", nil)
	assert.NoError(t, err)

	var result map[string]interface{}
	_, err = client.do(context.Background(), req, &result)
	assert.ErrorIs(t, err, ladoktypes.ErrNoValidContentType)
}

func TestEnvironment(t *testing.T) {
	tts := []struct {
		name    string
		env     string
		want    string
		wantErr bool
	}{
		{"IntTestAPI", ladoktypes.EnvIntTestAPI, ladoktypes.EnvIntTestAPI, false},
		{"TestAPI", ladoktypes.EnvTestAPI, ladoktypes.EnvTestAPI, false},
		{"ProdAPI", ladoktypes.EnvProdAPI, ladoktypes.EnvProdAPI, false},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			client := mockNewClient(t, tt.env, "https://example.com")
			got, err := client.environment(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsStudent(t *testing.T) {
	mux, server, client := mockSetup(t, ladoktypes.EnvIntTestAPI)
	defer server.Close()

	mockGenericEndpointServer(t, mux, ContentTypeStudentinformationJSON, "GET",
		fmt.Sprintf("/studentinformation/student/%s", ladokmocks.Students[0].StudentUID),
		ladokmocks.StudentJSON(ladokmocks.Students[0]), 200)

	_, err := client.IsStudent(context.Background(), &IsStudentReq{UID: ladokmocks.Students[0].StudentUID})
	assert.NoError(t, err)
}

func TestGetMyStudentDegrees(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	degrees, err := client.GetMyStudentDegrees(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, degrees)
}

func TestMarshalPDF(t *testing.T) {
	degrees := &MyStudentDegrees{}
	degrees.MarshalPDF()
}

func TestCheck(t *testing.T) {
	type sample struct {
		Name string `validate:"required"`
	}
	assert.Error(t, Check(&sample{}))
	assert.NoError(t, Check(&sample{Name: "test"}))
}

func TestGetBehorighetsprofil_ValidationError(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	_, _, err := client.Kataloginformation.GetBehorighetsprofil(context.Background(), &GetBehorighetsprofilerReq{})
	assert.Error(t, err)
}

func TestGetStudent_ValidationError(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	_, _, err := client.Studentinformation.GetStudent(context.Background(), &GetStudentReq{})
	assert.Error(t, err)
}

func TestGetAktivPaLarosate_ValidationError(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	_, _, err := client.Studentinformation.GetAktivPaLarosate(context.Background(), &GetAktivPaLarosateReq{})
	assert.Error(t, err)
}

func TestHistorical_ValidationError(t *testing.T) {
	client := mockNewClient(t, ladoktypes.EnvIntTestAPI, "https://example.com")
	_, _, err := client.Feed.Historical(context.Background(), &HistoricalReq{})
	assert.Error(t, err)
}
