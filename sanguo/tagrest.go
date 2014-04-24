package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"sanguo/tagging"
)

func main() {
	api := tagging.Api{}
	api.Init()

	handler := rest.ResourceHandler{EnableRelaxedContentType: true}
	handler.SetRoutes(
		rest.RouteObjectMethod("GET", "/api/resumes/p/:start", &api, "ListResumes"),
		rest.RouteObjectMethod("GET", "/api/resumes/:id", &api, "GetResumeById"),

		rest.RouteObjectMethod("POST", "/api/tags", &api, "SaveTag"),
		rest.RouteObjectMethod("DELETE", "/api/tags/:id", &api, "DeleteTagById"),

		rest.RouteObjectMethod("GET", "/api/tagformats", &api, "ListTagFormats"),
		rest.RouteObjectMethod("POST", "/api/tagformats", &api, "SaveTagFormat"),
		rest.RouteObjectMethod("DELETE", "/api/tagformats/:tag", &api, "DeleteTagFormatByTag"),
	)
	http.ListenAndServe(":8080", &handler)
}
