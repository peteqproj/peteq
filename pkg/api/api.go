package api

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type (
	// Resource represents set of endpoints
	Resource struct {
		// Path, including /
		Path        string
		Midderwares []gin.HandlerFunc
		Endpoints   []Endpoint
		Subresource []Resource
	}

	// Endpoint is one endpoint
	Endpoint struct {
		// HTTP verb
		Verb string
		// Path, including /
		Path        string
		Midderwares []gin.HandlerFunc
		Handler     gin.HandlerFunc
	}

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
	}
)

// NewRejectedCommandResponse build CommandResponse that is rejected with reason
func NewRejectedCommandResponse(reason string) CommandResponse {
	return CommandResponse{
		Status: "rejected",
		Reason: stringPtr(reason),
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

// WrapCommandAPI wrap command handler
func WrapCommandAPI(handler func(ctx context.Context, body io.ReadCloser) CommandResponse) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp := handler(c.Request.Context(), c.Request.Body)
		status := 200
		if resp.Reason != nil {
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
	return nil
}

// stringPtr return point to string
func stringPtr(str string) *string {
	return &str
}
