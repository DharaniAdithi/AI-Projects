package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"booking_system/booking_service/internal/dtos"
	"booking_system/booking_service/internal/errorcodes"
	"booking_system/booking_service/internal/utils"

	"github.com/stretchr/testify/assert"
)

// ---- WriteResponse Tests ----

func TestWriteResponse_Success(t *testing.T) {
	w := httptest.NewRecorder()
	payload := &dtos.APIResponse{
		Status:  "Success",
		Message: "OK",
	}

	utils.WriteResponse(w, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response dtos.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Success", response.Status)
}

func TestWriteResponse_StatusCodes(t *testing.T) {
	tests := []struct {
		name   string
		code   int
	}{
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
		{"Created", http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			utils.WriteResponse(w, tt.code, map[string]string{"msg": "test"})
			assert.Equal(t, tt.code, w.Code)
		})
	}
}

// ---- ValidateStruct Tests ----

type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,min=1"`
}

func TestValidateStruct_NoErrors(t *testing.T) {
	s := &TestStruct{
		Name:  "John",
		Email: "john@example.com",
		Age:   25,
	}

	errs := utils.ValidateStruct(s)
	assert.Empty(t, errs)
}

func TestValidateStruct_MissingRequired(t *testing.T) {
	s := &TestStruct{
		Name: "",  // required field missing
	}

	errs := utils.ValidateStruct(s)
	assert.NotEmpty(t, errs)

	// Check that Name field error is present
	fieldNames := make([]string, len(errs))
	for i, e := range errs {
		fieldNames[i] = e.Field
	}
	assert.Contains(t, fieldNames, "Name")
}

func TestValidateStruct_InvalidEmail(t *testing.T) {
	s := &TestStruct{
		Name:  "John",
		Email: "not-an-email",
		Age:   25,
	}

	errs := utils.ValidateStruct(s)
	assert.NotEmpty(t, errs)

	fieldNames := make([]string, len(errs))
	for i, e := range errs {
		fieldNames[i] = e.Field
	}
	assert.Contains(t, fieldNames, "Email")
}

func TestValidateStruct_ErrorCodeSet(t *testing.T) {
	s := &TestStruct{}  // all empty

	errs := utils.ValidateStruct(s)
	assert.NotEmpty(t, errs)
	for _, e := range errs {
		assert.NotEmpty(t, e.Code)
		assert.Equal(t, errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode], e.Code)
	}
}

// ---- MapErrorCode Tests ----

func TestMapErrorCode_NoErrors_ReturnsOK(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Status:  "Success",
		Message: "done",
		Errors:  []dtos.Error{},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusOK, result.Code)
}

func TestMapErrorCode_RecordNotFound_Returns404(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusNotFound, result.Code)
}

func TestMapErrorCode_ValidationError_Returns400(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusBadRequest, result.Code)
}

func TestMapErrorCode_UnexpectedError_Returns500(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusInternalServerError, result.Code)
}

func TestMapErrorCode_UnauthorizedError_Returns401(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.UnauthorizedErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusUnauthorized, result.Code)
}

func TestMapErrorCode_BadRequest_Returns400(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.BadRequestErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusBadRequest, result.Code)
}

func TestMapErrorCode_UnprocessableEntity_Returns422(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.UnprocessableEntityErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusUnprocessableEntity, result.Code)
}

func TestMapErrorCode_Forbidden_Returns403(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: errorcodes.ErrorCodeStatus[errorcodes.ForbiddenErrorCode]},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusForbidden, result.Code)
}

func TestMapErrorCode_UnknownCode_Returns500(t *testing.T) {
	apiResponse := &dtos.APIResponse{
		Errors: []dtos.Error{
			{Code: "UNKNOWN_CODE_XYZ"},
		},
	}

	result := utils.MapErrorCode(apiResponse)
	assert.Equal(t, http.StatusInternalServerError, result.Code)
}

