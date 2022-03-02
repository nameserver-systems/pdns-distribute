/*PowerDNS Authoritative HTTP API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.13
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
//nolint:golint
package modelpowerdns

// Returned when the server encounters an error. Either in client input or internally.
type ModelError struct {
	// A human readable error message
	Error string `json:"error"`
	// Optional array of multiple errors encountered during processing
	Errors []string `json:"errors,omitempty"`
}
