package trip

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(trip *Trip) error {
	args := m.Called(trip)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id int64) (*Trip, error) {
	args := m.Called(id)
	return args.Get(0).(*Trip), args.Error(1)
}

func (m *MockRepository) Update(trip *Trip) error {
	args := m.Called(trip)
	return args.Error(0)
}

func (m *MockRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTrip(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockRepository)
	h := NewHandler(mockRepo)

	trip := &Trip{
		Name:        "Test Trip",
		Description: "A test trip",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("Create", mock.AnythingOfType("*trip.Trip")).Return(nil)

	jsonBody, _ := json.Marshal(trip)
	req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.CreateTrip(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var createdTrip Trip
		json.Unmarshal(rec.Body.Bytes(), &createdTrip)
		assert.Equal(t, trip.Name, createdTrip.Name)
	}

	mockRepo.AssertExpectations(t)
}
