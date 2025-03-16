package openapi

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func NewDefinition() *openapi3.T {
	return &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       "Gitea-Issues-To-Redmine-Sync",
			Description: "This API provides endpoints to sync Gitea issues with redmine.",
			Version:     "1.0.0",
			Contact: &openapi3.Contact{
				Name:  "aborgardt",
				URL:   "https://github.com/b1tray3r/gitea-issues-to-redmine",
				Email: "5030347+b1tray3r@users.noreply.github.com",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "localhost",
				URL:         "http://localhost:8085",
			},
		},
		Paths: openapi3.NewPaths(
			openapi3.WithPath(
				"/gitea/webhook",
				&openapi3.PathItem{
					Summary:     "Webhook endpoint for Gitea.",
					Description: "Takes in the webhook content from Gitea",
					Post: &openapi3.Operation{
						Tags:    []string{"Gitea"},
						Summary: "Endpoint filled by a Gitea wehook.",
						RequestBody: &openapi3.RequestBodyRef{
							Ref: "#/components/requestBodies/PostGiteaWebhookRequestBody",
						},
						Responses: openapi3.NewResponses(
							openapi3.WithStatus(
								http.StatusOK,
								&openapi3.ResponseRef{
									Ref: "#/components/responses/WebhookAcceptedResponse",
								},
							),
						),
					},
				},
			),
			openapi3.WithPath(
				"/health",
				&openapi3.PathItem{
					Summary:     "HealthCheck without authentication",
					Description: "Checks the availablility of the API server.",
					Get: &openapi3.Operation{
						Tags:    []string{"General"},
						Summary: "Get health status",
						Responses: openapi3.NewResponses(
							openapi3.WithStatus(
								http.StatusOK,
								&openapi3.ResponseRef{
									Ref: "#/components/responses/HealthCheckResponse",
								},
							),
						),
					},
				},
			),
		),
		Components: &openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{
				"basicAuth": &openapi3.SecuritySchemeRef{
					Value: openapi3.NewSecurityScheme().
						WithDescription("HTTP Basic authentication").
						WithType("http").
						WithScheme("basic"),
				},
			},
			Schemas: openapi3.Schemas{
				"Status": openapi3.NewSchemaRef("",
					openapi3.NewObjectSchema().
						WithProperty(
							"message",
							openapi3.NewStringSchema().WithDefault("Some message from the backend!"),
						),
				),
				"GiteaPayload": openapi3.NewSchemaRef("",
					openapi3.NewObjectSchema().
						WithProperty("number", openapi3.NewInt64Schema()).
						WithProperty(
							"action",
							openapi3.NewStringSchema(),
						).WithProperty(
						"issue",
						openapi3.NewObjectSchema().
							WithProperty("title", openapi3.NewStringSchema()).
							WithProperty("body", openapi3.NewStringSchema()).
							WithProperty("url", openapi3.NewStringSchema()).
							WithProperty("user", openapi3.NewObjectSchema().
								WithProperty("email", openapi3.NewStringSchema()),
							),
					).WithProperty(
						"repository",
						openapi3.NewObjectSchema().
							WithProperty("owner", openapi3.NewObjectSchema().
								WithProperty("login", openapi3.NewStringSchema()),
							).WithProperty(
							"name",
							openapi3.NewStringSchema(),
						),
					),
				),
			},
			RequestBodies: map[string]*openapi3.RequestBodyRef{
				"PostGiteaWebhookRequestBody": {
					Value: openapi3.NewRequestBody().
						WithDescription("Request body for Gitea webhook").
						WithContent(
							openapi3.NewContentWithJSONSchemaRef(
								&openapi3.SchemaRef{
									Ref: "#/components/schemas/GiteaPayload",
								},
							),
						),
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"ErrorResponse": {
					Value: openapi3.NewResponse().
						WithDescription("Something went wrong!").
						WithContent(
							openapi3.NewContentWithJSONSchemaRef(
								&openapi3.SchemaRef{
									Ref: "#/components/schemas/Status",
								},
							),
						),
				},
				"WebhookAcceptedResponse": {
					Value: openapi3.NewResponse().
						WithDescription("Server is running!").
						WithContent(
							openapi3.NewContentWithJSONSchemaRef(
								&openapi3.SchemaRef{
									Ref: "#/components/schemas/Status",
								},
							),
						),
				},
				"HealthCheckResponse": {
					Value: openapi3.NewResponse().
						WithDescription("Server is running!").
						WithContent(
							openapi3.NewContentWithJSONSchemaRef(
								&openapi3.SchemaRef{
									Ref: "#/components/schemas/Status",
								},
							),
						),
				},
			},
		},
	}
}
