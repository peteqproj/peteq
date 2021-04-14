// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Oleg Sucharevich"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "List",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "RestAPI"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/list.List"
                            }
                        }
                    }
                }
            }
        },
        "/api/project": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Project",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "RestAPI"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/project.Project"
                            }
                        }
                    }
                }
            }
        },
        "/api/project/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Project",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "RestAPI"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/project.Project"
                        }
                    }
                }
            }
        },
        "/api/task/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "RestAPI"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/task.Task"
                            }
                        }
                    }
                }
            }
        },
        "/api/task/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Task",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "RestAPI"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/task.Task"
                        }
                    }
                }
            }
        },
        "/c/list/moveTasks": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Move tasks from source to destination list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "List Command API"
                ],
                "parameters": [
                    {
                        "description": "Move tasks from source to destination list",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/MoveTasksRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/project/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project Command API"
                ],
                "parameters": [
                    {
                        "description": "Add tasks into project",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateProjectRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/task/complete": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Complete task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task Command API"
                ],
                "parameters": [
                    {
                        "description": "Complete Task Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/task.completeReopenTaskRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/task/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task Command API"
                ],
                "parameters": [
                    {
                        "description": "Create Task Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateTaskResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/task/delete": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task Command API"
                ],
                "parameters": [
                    {
                        "description": "Delete Task Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/task.deleteTaskRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/task/reopen": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Reopen task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task Command API"
                ],
                "parameters": [
                    {
                        "description": "Reopen Task Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/task.completeReopenTaskRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/task/update": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Task Command API"
                ],
                "parameters": [
                    {
                        "description": "Update Task Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/task.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/trigger/run": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Trigger run",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Trigger Command API"
                ],
                "parameters": [
                    {
                        "description": "Trigger run",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TriggerRunRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/user/login": {
            "post": {
                "description": "Login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Command API"
                ],
                "parameters": [
                    {
                        "description": "Login",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/LoginRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/c/user/register": {
            "post": {
                "description": "Register new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Command API"
                ],
                "parameters": [
                    {
                        "description": "Register new user",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RegistrationRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/CommandResponse"
                        }
                    }
                }
            }
        },
        "/q/backlog": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Backlog View",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "View"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/backlog.backlogView"
                        }
                    }
                }
            }
        },
        "/q/home": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Home View",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "View"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/home.homeView"
                        }
                    }
                }
            }
        },
        "/q/project": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Projects View",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "View"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/projects.projectsView"
                        }
                    }
                }
            }
        },
        "/q/project/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Single Project View",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "View"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/project.projectView"
                        }
                    }
                }
            }
        },
        "/q/triggers": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Triggers View",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "View"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/triggers.triggersView"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "AddTasksRequestBody": {
            "type": "object",
            "required": [
                "project",
                "tasks"
            ],
            "properties": {
                "project": {
                    "type": "string"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "CommandResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Data to pass additional data",
                    "type": "object"
                },
                "id": {
                    "description": "ID the id of the related entity",
                    "type": "string"
                },
                "reason": {
                    "description": "Reason exist when command is rejected",
                    "type": "string"
                },
                "status": {
                    "description": "Status can be accepted or rejected",
                    "type": "string"
                },
                "type": {
                    "description": "Type the type of the related entity",
                    "type": "string"
                }
            }
        },
        "CreateProjectRequestBody": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "color": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "imageUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "CreateTaskResponse": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "list": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "project": {
                    "type": "string"
                }
            }
        },
        "LoginRequestBody": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "MoveTasksRequestBody": {
            "type": "object",
            "required": [
                "tasks"
            ],
            "properties": {
                "destination": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "RegistrationRequestBody": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "TriggerRunRequestBody": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "data": {
                    "type": "object"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "backlog.backlogTask": {
            "type": "object",
            "properties": {
                "list": {
                    "$ref": "#/definitions/backlog.backlogTaskList"
                },
                "project": {
                    "$ref": "#/definitions/backlog.backlogTaskProject"
                },
                "task": {
                    "$ref": "#/definitions/task.Task"
                }
            }
        },
        "backlog.backlogTaskList": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "backlog.backlogTaskProject": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "backlog.backlogView": {
            "type": "object",
            "properties": {
                "lists": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/backlog.backlogTaskList"
                    }
                },
                "projects": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/backlog.backlogTaskProject"
                    }
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/backlog.backlogTask"
                    }
                }
            }
        },
        "home.homeList": {
            "type": "object",
            "properties": {
                "metadata": {
                    "$ref": "#/definitions/list.Metadata"
                },
                "spec": {
                    "$ref": "#/definitions/list.Spec"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/home.homeTask"
                    }
                }
            }
        },
        "home.homeTask": {
            "type": "object",
            "properties": {
                "project": {
                    "$ref": "#/definitions/repo.Resource"
                },
                "task": {
                    "$ref": "#/definitions/task.Task"
                }
            }
        },
        "home.homeView": {
            "type": "object",
            "properties": {
                "lists": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/home.homeList"
                    }
                }
            }
        },
        "list.List": {
            "type": "object",
            "properties": {
                "metadata": {
                    "$ref": "#/definitions/list.Metadata"
                },
                "spec": {
                    "$ref": "#/definitions/list.Spec"
                }
            }
        },
        "list.Metadata": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "list.Spec": {
            "type": "object",
            "properties": {
                "index": {
                    "type": "number"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "project.Metadata": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "project.Project": {
            "type": "object",
            "properties": {
                "metadata": {
                    "$ref": "#/definitions/project.Metadata"
                },
                "spec": {
                    "$ref": "#/definitions/project.Spec"
                }
            }
        },
        "project.Spec": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "imageUrl": {
                    "type": "string"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "project.projectView": {
            "type": "object",
            "properties": {
                "project": {
                    "$ref": "#/definitions/project.Project"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/task.Task"
                    }
                }
            }
        },
        "projects.populatedProject": {
            "type": "object",
            "properties": {
                "project": {
                    "$ref": "#/definitions/project.Project"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/task.Task"
                    }
                }
            }
        },
        "projects.projectsView": {
            "type": "object",
            "properties": {
                "projects": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/projects.populatedProject"
                    }
                }
            }
        },
        "repo.Metadata": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "repo.Resource": {
            "type": "object",
            "properties": {
                "metadata": {
                    "$ref": "#/definitions/repo.Metadata"
                },
                "spec": {
                    "type": "object"
                }
            }
        },
        "task.Metadata": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "task.Spec": {
            "type": "object",
            "properties": {
                "completed": {
                    "type": "boolean"
                }
            }
        },
        "task.Task": {
            "type": "object",
            "properties": {
                "metadata": {
                    "$ref": "#/definitions/task.Metadata"
                },
                "spec": {
                    "$ref": "#/definitions/task.Spec"
                }
            }
        },
        "task.completeReopenTaskRequestBody": {
            "type": "object",
            "properties": {
                "task": {
                    "type": "string"
                }
            }
        },
        "task.deleteTaskRequestBody": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "triggers.triggerExecutionHistoryItem": {
            "type": "object",
            "properties": {
                "manual": {
                    "type": "boolean"
                },
                "triggeredAt": {
                    "type": "string"
                }
            }
        },
        "triggers.triggerViewItem": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "history": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/triggers.triggerExecutionHistoryItem"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "spec": {
                    "type": "object"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "triggers.triggersView": {
            "type": "object",
            "properties": {
                "triggers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/triggers.triggerViewItem"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Peteq API",
	Description: "Peteq OpenAPI spec.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
