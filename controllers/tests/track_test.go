package mock_controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/cschen13/spotitube/controllers"
	"github.com/gorilla/mux"
)

func TestGetTracks(t *testing.T) {
	tests := []struct {
		requestVars map[string]string
	}{
		{nil},
	}

	// TODO: Create SessionManager interface so it can be mocked
	ctrl := controllers.NewTrackController()
	router := mux.NewRouter()
	ctrl.Register(router)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest()
	router.ServeHTTP(rr, req)

	// make assertions
}
