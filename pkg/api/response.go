package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CommandResponse unified response for any command
	CommandResponse struct {
		// Status can be accepted or rejected
		Status string `json:"status"`
		// Reason exist when command is rejected
		Reason *string `json:"reason"`
		// Type the type of the related entity
		Type string `json:"type"`
		// ID the id of the related entity
		ID string `json:"id"`
		// Data to pass additional data
		Data interface{} `json:"data"`
	}
)

// NewRejectedCommandResponse build CommandResponse that is rejected with reason
func NewRejectedCommandResponse(err error) CommandResponse {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		msg := []string{}
		for _, e := range errs {
			msg = append(msg, fmt.Sprintf("Error: %s %s", e.Field(), e.ActualTag()))
		}
		return CommandResponse{
			Status: "rejected",
			Reason: stringPtr(strings.Join(msg, " | ")),
		}
	}
	return CommandResponse{
		Status: "rejected",
		Reason: stringPtr(err.Error()),
	}
}

// NewAcceptedCommandResponse build CommandResponse that is rejected with reason
func NewAcceptedCommandResponse(resource string, id string) CommandResponse {
	return CommandResponse{
		Status: "accepted",
		Type:   resource,
		ID:     id,
	}
}

// NewAcceptedCommandResponseWithData build CommandResponse that is rejected with reason
func NewAcceptedCommandResponseWithData(resource string, id string, data interface{}) CommandResponse {
	return CommandResponse{
		Status: "accepted",
		Type:   resource,
		ID:     id,
		Data:   data,
	}
}

// WrapCommandAPI wrap command handler
func WrapCommandAPI(handler func(ctx context.Context, body io.ReadCloser) CommandResponse, logger logger.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := handler(c.Request.Context(), c.Request.Body)
		status := 200
		if resp.Reason != nil {
			logger.Info(*resp.Reason)
			status = 400
		}
		c.JSON(status, resp)
	}
}

// UnmarshalInto unmarshal request body into struct
func UnmarshalInto(body io.ReadCloser, into interface{}) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(bytes), into); err != nil {
		return err
	}
	err = validator.New().Struct(into)
	return err
}

// stringPtr return point to string
func stringPtr(str string) *string {
	return &str
}
