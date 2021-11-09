package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
	"schemes": [
		"http",
		"https"
	],
	"tags": [
        {
          "name": "Assets",
          "description": "Everything about Assets",
          "externalDocs": {
            "description": "Find out more",
            "url": "http://comdex.one"
          }
        }
	],
    "paths": {
        "/comdex/asset/v1beta1/assets/{assetId}": {
          "get": {
            "description": "Unique identifier of an asset.",
            "consumes": [
              "text/plain"
            ],
            "produces": [
              "application/json"
            ],
            "tags": [
              "Assets"
            ],
            "summary": "Search for an asset by Asset ID",
            "parameters": [
              {
                "type": "string",
                "description": "Asset ID",
                "name": "assetId",
                "in": "path",
                "required": true
              }
            ],
            "responses": {
              "200": {
                "description": "Message for a successful search.",
                "schema": {
                  "$ref": "#/definitions/asset.queryResponse.success"
                }
              },
              "default": {
                "description": "Message for an unexpected error.",
                "schema": {
                  "$ref": "#/definitions/asset.queryResponse.failed"
                }
              }
            }
          }
        }
      },
    "definitions": {
        "asset.queryResponse.success": {
          "type": "object",
          "properties": {
            "asset": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string"
                },
                "name": {
                  "type": "string"
                },
                "denom": {
                  "type": "string"
                },
                "decimals": {
                  "type": "string"
                }
              }
            }
          }
        },
        "asset.queryResponse.failed": {
          "type": "object",
          "properties": {
            "code": {
              "type": "integer"
            },
            "message": {
              "type": "string"
            },
            "details": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/helpers.Mappable"
              }
            }
          }
        },
        "helpers.Mappable": {
          "type": "object"
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
	Version:     "0.1.0",
	Host:        "localhost:1317",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Comdex Swagger Documentation",
	Description: "API Documentation of Comdex custom modules",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
