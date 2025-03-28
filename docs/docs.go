// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/auth/admin/register": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Registers a new admin with the given email, password, and name. Validates the input, checks password match and creates the user in the database.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Register a new admin",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "User registration details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterUserDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User successfully registered",
                        "schema": {
                            "$ref": "#/definitions/handlers.AuthResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Invalid input or user already exists",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Authenticates a user by validating the email and password, and generates an authentication token if successful.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User login details",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Authentication successful, token returned",
                        "schema": {
                            "$ref": "#/definitions/handlers.AuthResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "401": {
                        "description": "Incorrect email or password",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Registers a new user with the given email, password, and name. Validates the input, checks password match, assigns a default role, and creates the user in the database.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterUserDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User successfully registered",
                        "schema": {
                            "$ref": "#/definitions/handlers.AuthResponseDTO"
                        }
                    },
                    "400": {
                        "description": "Invalid input or user already exists",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            }
        },
        "/inventory": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves a list of inventory items from the database. Optionally filters the results to show items with low stock, based on the query parameter ` + "`" + `lowStock` + "`" + `.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inventory"
                ],
                "summary": "Get a list of inventory items",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Filter by low stock (quantity \u003c= LOW_STOCK_THRESHOLD)",
                        "name": "lowStock",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of inventory items",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.InventoryItem"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Creates a new item in the inventory with the provided details, validates input, and handles errors such as duplicate items.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inventory"
                ],
                "summary": "Create a new inventory item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Inventory item details",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.InventoryItem"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Item successfully created",
                        "schema": {
                            "$ref": "#/definitions/models.InventoryItem"
                        }
                    },
                    "400": {
                        "description": "Invalid input or item already exists",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            }
        },
        "/inventory/restock": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves the restock history for inventory items. Optionally filters the history based on the ` + "`" + `itemId` + "`" + ` query parameter.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inventory",
                    "Restock"
                ],
                "summary": "Get restock history for inventory items",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Filter restock history by item ID",
                        "name": "itemId",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of restock history records",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.RestockHistoryDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "404": {
                        "description": "Item not found",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            }
        },
        "/inventory/{itemID}/restock": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Restocks an inventory item with the provided quantity, checks if restock quota is reached, and handles errors such as invalid itemID, quota exceeded, or database issues.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inventory"
                ],
                "summary": "Restock an inventory item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Inventory Item ID",
                        "name": "itemID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Restock details",
                        "name": "restock",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RestockItemDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Item successfully restocked",
                        "schema": {
                            "$ref": "#/definitions/models.Restock"
                        }
                    },
                    "400": {
                        "description": "Invalid itemID or invalid input",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "404": {
                        "description": "Item not found",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "409": {
                        "description": "Item already exists",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "429": {
                        "description": "Item quota reached",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/api.RequestError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.RequestError": {
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
        "handlers.AuthResponseDTO": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/handlers.UserResponseDTO"
                }
            }
        },
        "handlers.LoginDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handlers.RegisterUserDTO": {
            "type": "object",
            "properties": {
                "confirmPassword": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handlers.RestockHistoryDTO": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "itemID": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                }
            }
        },
        "handlers.RestockItemDTO": {
            "type": "object",
            "properties": {
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "handlers.UserResponseDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "models.InventoryItem": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "lastRestock": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "models.Restock": {
            "type": "object",
            "properties": {
                "inventoryItem": {
                    "$ref": "#/definitions/models.InventoryItem"
                },
                "name": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "Go challenge.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
