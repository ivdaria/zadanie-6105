//go:build tools
// +build tools

package tools

//go:generate oapi-codegen -config config.yaml openapi.yml
