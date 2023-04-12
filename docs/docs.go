// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth": {
            "get": {
                "description": "check is user authorised",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Auth",
                "responses": {
                    "200": {
                        "description": "success auth"
                    },
                    "401": {
                        "description": "failed auth",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/create_csrf": {
            "get": {
                "description": "Get CSRF token",
                "tags": [
                    "auth"
                ],
                "summary": "GetCSRF",
                "responses": {
                    "200": {
                        "description": "success create csrf"
                    },
                    "401": {
                        "description": "failed get user",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/folder/{slug}": {
            "get": {
                "description": "List of folder messages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "GetFolderMessages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "FolderSlug",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success get list of folder messages",
                        "schema": {
                            "$ref": "#/definitions/models.FolderResponse"
                        }
                    },
                    "400": {
                        "description": "invalid url address",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "folder not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/folders/": {
            "get": {
                "description": "List of outgoing messages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "GetFolders",
                "responses": {
                    "200": {
                        "description": "success get list of outgoing messages",
                        "schema": {
                            "$ref": "#/definitions/models.FoldersResponse"
                        }
                    },
                    "400": {
                        "description": "failed to get user",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/logout": {
            "delete": {
                "description": "check is user authorised",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout",
                "responses": {
                    "200": {
                        "description": "success logout"
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/message/send": {
            "post": {
                "description": "Message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "SendMessage",
                "responses": {
                    "200": {
                        "description": "success send message",
                        "schema": {
                            "$ref": "#/definitions/models.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "no valid emails",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "403": {
                        "description": "invalid form",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "message not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/message/{id}": {
            "get": {
                "description": "Message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "GetMessage",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success get messages",
                        "schema": {
                            "$ref": "#/definitions/models.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "invalid url address",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "message not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/message/{id}/read": {
            "post": {
                "description": "Message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "ReadMessage",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success read message",
                        "schema": {
                            "$ref": "#/definitions/models.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "invalid url address",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "message not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/message/{id}/unread": {
            "post": {
                "description": "Message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messages"
                ],
                "summary": "UnreadMessage",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success unread message",
                        "schema": {
                            "$ref": "#/definitions/models.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "invalid url address",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "message not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/signin": {
            "post": {
                "description": "user sign in",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SignIn",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.FormLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success sign in",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponse"
                        }
                    },
                    "401": {
                        "description": "wrong password",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "403": {
                        "description": "invalid form",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "description": "user sign up",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SignUp",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.FormSignUp"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "user created",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponse"
                        }
                    },
                    "401": {
                        "description": "invalid login",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "403": {
                        "description": "password too short",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "409": {
                        "description": "user already exists",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/user": {
            "delete": {
                "description": "delete user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete",
                "responses": {
                    "200": {
                        "description": "success delete user"
                    },
                    "400": {
                        "description": "failed to get user",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/user/avatar": {
            "get": {
                "description": "get user avatar",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "GetAvatar",
                "parameters": [
                    {
                        "type": "string",
                        "description": "email",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success get user avatar",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    },
                    "400": {
                        "description": "no bucket",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            },
            "put": {
                "description": "edit user avatar",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "EditAvatar",
                "responses": {
                    "200": {
                        "description": "success edit user avatar"
                    },
                    "400": {
                        "description": "unsupported content type",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "403": {
                        "description": "invalid form",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/user/info": {
            "get": {
                "description": "get info about request creator",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "GetPersonalInfo",
                "responses": {
                    "200": {
                        "description": "success get user info",
                        "schema": {
                            "$ref": "#/definitions/models.UserInfo"
                        }
                    },
                    "401": {
                        "description": "failed get user",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            },
            "put": {
                "description": "edit info about user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "EditInfo",
                "responses": {
                    "200": {
                        "description": "success edit user info",
                        "schema": {
                            "$ref": "#/definitions/models.EditUserInfoResponse"
                        }
                    },
                    "401": {
                        "description": "failed to get user",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "403": {
                        "description": "invalid form",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/user/info/{email}": {
            "get": {
                "description": "get info about user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "GetInfo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "email",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success get user info",
                        "schema": {
                            "$ref": "#/definitions/models.UserInfo"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        },
        "/user/pw": {
            "put": {
                "description": "edit password about user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "EditPw",
                "responses": {
                    "200": {
                        "description": "success edit user password"
                    },
                    "400": {
                        "description": "failed to get user",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "403": {
                        "description": "invalid form",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "404": {
                        "description": "user not found",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/errors.JSONError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.JSONError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.AuthResponse": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                }
            }
        },
        "models.EditUserInfoResponse": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "models.Folder": {
            "type": "object",
            "properties": {
                "folder_id": {
                    "type": "integer"
                },
                "folder_slug": {
                    "type": "string"
                },
                "messages_count": {
                    "type": "integer"
                },
                "messages_unseen": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.FolderResponse": {
            "type": "object",
            "properties": {
                "folder": {
                    "$ref": "#/definitions/models.Folder"
                },
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MessageInfo"
                    }
                }
            }
        },
        "models.FoldersResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "folders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Folder"
                    }
                }
            }
        },
        "models.FormLogin": {
            "type": "object",
            "required": [
                "login",
                "password",
                "remember"
            ],
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "remember": {
                    "type": "boolean"
                }
            }
        },
        "models.FormSignUp": {
            "type": "object",
            "required": [
                "firstName",
                "lastName",
                "login",
                "password",
                "repeatPw"
            ],
            "properties": {
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "repeatPw": {
                    "type": "string"
                }
            }
        },
        "models.MessageInfo": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted": {
                    "type": "boolean"
                },
                "favorite": {
                    "type": "boolean"
                },
                "from_user_id": {
                    "$ref": "#/definitions/models.UserInfo"
                },
                "message_id": {
                    "type": "integer"
                },
                "recipients": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UserInfo"
                    }
                },
                "reply_to": {
                    "$ref": "#/definitions/models.MessageInfo"
                },
                "seen": {
                    "type": "boolean"
                },
                "text": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "$ref": "#/definitions/models.MessageInfo"
                }
            }
        },
        "models.UserInfo": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8001",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "MailBox Swagger API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
