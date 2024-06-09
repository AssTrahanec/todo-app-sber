package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	todoListSber "todo-list-sber"
	"todo-list-sber/pkg/service"
	servicemocks "todo-list-sber/pkg/service/mocks"
)

func TestCreateTodoItemHandler(t *testing.T) {
	type mockBehavior func(r *servicemocks.MockTodoItem, item todoListSber.TodoItem)

	tests := []struct {
		name                 string
		inputBody            string
		inputItem            todoListSber.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title": "Test Task", "description": "Test Description", "date": "2024-06-05T20:00:00Z", "is_done": true}`,
			inputItem: todoListSber.TodoItem{
				Title:       "Test Task",
				Description: "Test Description",
				Date:        time.Date(2024, time.June, 5, 20, 0, 0, 0, time.UTC),
				IsDone:      true,
			},
			mockBehavior: func(r *servicemocks.MockTodoItem, item todoListSber.TodoItem) {
				r.EXPECT().Create(item).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{}`,
			inputItem:            todoListSber.TodoItem{},
			mockBehavior:         func(r *servicemocks.MockTodoItem, item todoListSber.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"title": "Test Task", "description": "Test Description", "date": "2024-06-05T20:00:00Z", "is_done": true}`,
			inputItem: todoListSber.TodoItem{
				Title:       "Test Task",
				Description: "Test Description",
				Date:        time.Date(2024, time.June, 5, 20, 0, 0, 0, time.UTC),
				IsDone:      true,
			},
			mockBehavior: func(r *servicemocks.MockTodoItem, item todoListSber.TodoItem) {
				r.EXPECT().Create(item).Return(0, errors.New("Something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"Something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			test.mockBehavior(mockTodoItem, test.inputItem)

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			router := gin.New()
			router.POST("/api/todo", handler.createTodoItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/todo",
				bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
func TestGetAllTodoItemsHandler(t *testing.T) {
	tests := []struct {
		name                 string
		mockBehavior         func(r *servicemocks.MockTodoItem)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success",
			mockBehavior: func(r *servicemocks.MockTodoItem) {
				expectedTodoItems := []todoListSber.TodoItem{
					{Id: 1, Title: "Task 1", Description: "Description 1", Date: time.Date(2024, time.June, 5, 20, 0, 0, 0, time.UTC), IsDone: false},
					{Id: 2, Title: "Task 2", Description: "Description 2", Date: time.Date(2024, time.June, 5, 20, 0, 0, 0, time.UTC), IsDone: false},
				}
				r.EXPECT().GetAll().Return(expectedTodoItems, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":[{"id":1,"title":"Task 1","description":"Description 1","date":"2024-06-05T20:00:00Z","is_done":false},{"id":2,"title":"Task 2","description":"Description 2","date":"2024-06-05T20:00:00Z","is_done":false}]}`,
		},
		{
			name: "Service Error",
			mockBehavior: func(r *servicemocks.MockTodoItem) {
				r.EXPECT().GetAll().Return(nil, errors.New("Service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			test.mockBehavior(mockTodoItem)

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			r := gin.New()
			r.GET("api/todo", handler.getAllTodoItems)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/todo", nil)

			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d; got %d", test.expectedStatusCode, w.Code)
			}

			if w.Body.String() != test.expectedResponseBody {
				t.Errorf("expected response body %q; got %q", test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
func TestGetTodoItemByIdHandler(t *testing.T) {
	tests := []struct {
		name                 string
		idParam              string
		mockBehavior         func(r *servicemocks.MockTodoItem, id int)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Success",
			idParam: "1",
			mockBehavior: func(r *servicemocks.MockTodoItem, id int) {
				expectedTodoItem := todoListSber.TodoItem{
					Id:          1,
					Title:       "Task 1",
					Description: "Description 1",
					Date:        time.Date(2024, time.June, 5, 20, 0, 0, 0, time.UTC),
					IsDone:      false,
				}
				r.EXPECT().GetById(id).Return(expectedTodoItem, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"data":{"id":1,"title":"Task 1","description":"Description 1","date":"2024-06-05T20:00:00Z","is_done":false}}`,
		},
		{
			name:                 "Invalid ID",
			idParam:              "invalid",
			mockBehavior:         func(r *servicemocks.MockTodoItem, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid ID"}`,
		},
		{
			name: "Service Error",
			mockBehavior: func(r *servicemocks.MockTodoItem, id int) {
				r.EXPECT().GetAll().Return(nil, errors.New("Service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id, err := strconv.Atoi(test.idParam)
			if err != nil {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Invalid ID"}`))
				return
			}

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			test.mockBehavior(mockTodoItem, id)

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			r := gin.New()
			r.GET("/api/todo/:id", handler.getTodoItemById)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/todo/"+test.idParam, nil)

			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d; got %d", test.expectedStatusCode, w.Code)
			}

			if w.Body.String() != test.expectedResponseBody {
				t.Errorf("expected response body %q; got %q", test.expectedResponseBody, w.Body.String())
			}
		})
	}

}
func TestUpdateTodoItemHandler(t *testing.T) {
	tests := []struct {
		name                 string
		idParam              string
		inputBody            string
		mockBehavior         func(r *servicemocks.MockTodoItem, id int, input todoListSber.UpdateItemInput)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Success",
			idParam:   "1",
			inputBody: `{"title": "Updated Task", "description": "Updated Description", "date": "2024-06-07T20:00:00Z", "is_done": true}`,
			mockBehavior: func(r *servicemocks.MockTodoItem, id int, input todoListSber.UpdateItemInput) {
				title := "Updated Task"
				description := "Updated Description"
				isDone := true
				date := time.Date(2024, time.June, 7, 20, 0, 0, 0, time.UTC)

				input = todoListSber.UpdateItemInput{
					Title:       &title,
					Description: &description,
					IsDone:      &isDone,
					Date:        &date,
				}
				r.EXPECT().Update(id, input).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid ID",
			idParam:              "invalid",
			inputBody:            `{}`,
			mockBehavior:         func(r *servicemocks.MockTodoItem, id int, input todoListSber.UpdateItemInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid ID"}`,
		},
		{
			name:                 "Invalid Input Body",
			idParam:              "1",
			inputBody:            `{}`,
			mockBehavior:         func(r *servicemocks.MockTodoItem, id int, input todoListSber.UpdateItemInput) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"EOF"}`,
		},
		{
			name:      "Service Error",
			idParam:   "1",
			inputBody: `{"title": "Updated Task", "description": "Updated Description", "date": "2024-06-07T20:00:00Z", "is_done": true}`,
			mockBehavior: func(r *servicemocks.MockTodoItem, id int, input todoListSber.UpdateItemInput) {
				title := "Updated Task"
				description := "Updated Description"
				isDone := true
				date := time.Date(2024, time.June, 7, 20, 0, 0, 0, time.UTC)

				input = todoListSber.UpdateItemInput{
					Title:       &title,
					Description: &description,
					IsDone:      &isDone,
					Date:        &date,
				}
				r.EXPECT().Update(id, input).Return(errors.New("Service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id, err := strconv.Atoi(test.idParam)
			if err != nil {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Invalid ID"}`))
				return
			}

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			test.mockBehavior(mockTodoItem, id, todoListSber.UpdateItemInput{})

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			r := gin.New()
			r.PUT("/api/todo/:id", handler.updateTodoItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/todo/"+test.idParam, bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d; got %d", test.expectedStatusCode, w.Code)
			}

			if w.Body.String() != test.expectedResponseBody {
				t.Errorf("expected response body %q; got %q", test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
func TestDeleteTodoItemHandler(t *testing.T) {
	tests := []struct {
		name                 string
		idParam              string
		mockBehavior         func(r *servicemocks.MockTodoItem, id int)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Success",
			idParam: "1",
			mockBehavior: func(r *servicemocks.MockTodoItem, id int) {
				r.EXPECT().Delete(id).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid ID",
			idParam:              "invalid",
			mockBehavior:         func(r *servicemocks.MockTodoItem, id int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"invalid\": invalid syntax"}`,
		},
		{
			name:    "Service Error",
			idParam: "1",
			mockBehavior: func(r *servicemocks.MockTodoItem, id int) {
				r.EXPECT().Delete(id).Return(errors.New("Service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id, err := strconv.Atoi(test.idParam)
			if err != nil {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Invalid ID"}`))
				return
			}

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			test.mockBehavior(mockTodoItem, id)

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			r := gin.New()
			r.DELETE("/api/todo/:id", handler.deleteTodoItem)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/todo/"+test.idParam, nil)

			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d; got %d", test.expectedStatusCode, w.Code)
			}

			if w.Body.String() != test.expectedResponseBody {
				t.Errorf("expected response body %q; got %q", test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
func TestGetDoneTodoItemsHandler(t *testing.T) {
	tests := []struct {
		name                 string
		queryParams          string
		mockBehavior         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Success",
			queryParams: "?date=2024-06-08&limit=10&offset=0",
			mockBehavior: func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {
				expectedDate, _ := time.Parse("2006-01-02", "2024-06-08")
				expectedTodos := []todoListSber.TodoItem{
					{
						Id:          1,
						Title:       "task 1",
						Description: "description 1",
						Date:        expectedDate,
						IsDone:      true,
					},
					{
						Id:          2,
						Title:       "task 2",
						Description: "description 2",
						Date:        expectedDate,
						IsDone:      true,
					},
				}
				r.EXPECT().GetDoneTodoItems(&expectedDate, 10, 0).Return(expectedTodos, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"title":"task 1","description":"description 1","date":"2024-06-08T00:00:00Z","is_done":true},{"id":2,"title":"task 2","description":"description 2","date":"2024-06-08T00:00:00Z","is_done":true}]`,
		},
		{
			name:                 "Invalid Date",
			queryParams:          "?date=invalid-date&limit=10&offset=0",
			mockBehavior:         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid date format"}`,
		},
		{
			name:                 "Invalid Limit",
			queryParams:          "?date=2024-06-08&limit=-1&offset=0",
			mockBehavior:         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid limit"}`,
		},
		{
			name:                 "Invalid Offset",
			queryParams:          "?date=2024-06-08&limit=10&offset=-1",
			mockBehavior:         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid offset"}`,
		},
		{
			name:        "Service Error",
			queryParams: "?date=2024-06-08&limit=10&offset=0",
			mockBehavior: func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {
				expectedDate, _ := time.Parse("2006-01-02", "2024-06-08")
				r.EXPECT().GetDoneTodoItems(&expectedDate, 10, 0).Return(nil, errors.New("Service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			if test.name != "Invalid Date" && test.name != "Invalid Limit" && test.name != "Invalid Offset" {

				dateStr := "2024-06-08"
				if test.name == "Invalid Date" {
					dateStr = "invalid-date"
				}
				date, _ := time.Parse("2006-01-02", dateStr)
				limit, _ := strconv.Atoi("10")
				offset, _ := strconv.Atoi("0")

				test.mockBehavior(mockTodoItem, &date, limit, offset)
			} else {
				test.mockBehavior(mockTodoItem, nil, 0, 0)
			}

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			r := gin.New()
			r.GET("/api/todo/done", handler.GetDoneTodoItems)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/todo/done"+test.queryParams, nil)

			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d; got %d", test.expectedStatusCode, w.Code)
			}

			if w.Body.String() != test.expectedResponseBody {
				t.Errorf("expected response body %q; got %q", test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
func TestGetUndoneTodoItemsHandler(t *testing.T) {
	tests := []struct {
		name                 string
		queryParams          string
		mockBehavior         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Success",
			queryParams: "?date=2024-06-08&limit=10&offset=0",
			mockBehavior: func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {
				expectedDate, _ := time.Parse("2006-01-02", "2024-06-08")
				expectedTodos := []todoListSber.TodoItem{
					{
						Id:          1,
						Title:       "task 1",
						Description: "description 1",
						Date:        expectedDate,
						IsDone:      false,
					},
					{
						Id:          2,
						Title:       "task 2",
						Description: "description 2",
						Date:        expectedDate,
						IsDone:      false,
					},
				}
				r.EXPECT().GetDoneTodoItems(&expectedDate, 10, 0).Return(expectedTodos, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"title":"task 1","description":"description 1","date":"2024-06-08T00:00:00Z","is_done":false},{"id":2,"title":"task 2","description":"description 2","date":"2024-06-08T00:00:00Z","is_done":false}]`,
		},
		{
			name:                 "Invalid Date",
			queryParams:          "?date=invalid-date&limit=10&offset=0",
			mockBehavior:         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid date format"}`,
		},
		{
			name:                 "Invalid Limit",
			queryParams:          "?date=2024-06-08&limit=-1&offset=0",
			mockBehavior:         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid limit"}`,
		},
		{
			name:                 "Invalid Offset",
			queryParams:          "?date=2024-06-08&limit=10&offset=-1",
			mockBehavior:         func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"Invalid offset"}`,
		},
		{
			name:        "Service Error",
			queryParams: "?date=2024-06-08&limit=10&offset=0",
			mockBehavior: func(r *servicemocks.MockTodoItem, date *time.Time, limit int, offset int) {
				expectedDate, _ := time.Parse("2006-01-02", "2024-06-08")
				r.EXPECT().GetDoneTodoItems(&expectedDate, 10, 0).Return(nil, errors.New("Service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"Service error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTodoItem := servicemocks.NewMockTodoItem(ctrl)
			if test.name != "Invalid Date" && test.name != "Invalid Limit" && test.name != "Invalid Offset" {

				dateStr := "2024-06-08"
				if test.name == "Invalid Date" {
					dateStr = "invalid-date"
				}
				date, _ := time.Parse("2006-01-02", dateStr)
				limit, _ := strconv.Atoi("10")
				offset, _ := strconv.Atoi("0")

				test.mockBehavior(mockTodoItem, &date, limit, offset)
			} else {
				test.mockBehavior(mockTodoItem, nil, 0, 0)
			}

			services := &service.Service{TodoItem: mockTodoItem}
			handler := Handler{services}

			r := gin.New()
			r.GET("/api/todo/done", handler.GetDoneTodoItems)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/todo/done"+test.queryParams, nil)

			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status %d; got %d", test.expectedStatusCode, w.Code)
			}

			if w.Body.String() != test.expectedResponseBody {
				t.Errorf("expected response body %q; got %q", test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
