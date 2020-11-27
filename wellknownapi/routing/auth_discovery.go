package routing

import (
	"net/http"
	"net/url"

	"github.com/matrix-org/dendrite/internal/httputil"
	"github.com/matrix-org/util"
)

// Almost all IANA registered metadatas
// See https://www.iana.org/assignments/oauth-parameters/oauth-parameters.xhtml#authorization-server-metadata
type discoveryDocument struct {
	Issuer                                             string   `json:"issuer"`
	AuthorizationEndpoint                              string   `json:"authorization_endpoint,omitempty"`
	TokenEndpoint                                      string   `json:"token_endpoint,omitempty"`
	JWKSUri                                            string   `json:"jwks_uri,omitempty"`
	RegistrationEndpoint                               string   `json:"registration_endpoint,omitempty"`
	ScopesSupported                                    []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported                             []string `json:"response_types_supported,omitempty"`
	ResponseModesSupported                             []string `json:"response_modes_supported,omitempty"`
	GrantTypesSupported                                []string `json:"grant_types_supported,omitempty"`
	TokenEndpointAuthMethodsSupported                  []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported         []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	ServiceDocumentation                               string   `json:"service_documentation,omitempty"`
	UILocalesSupported                                 []string `json:"ui_locales_supported,omitempty"`
	OpPolicyURI                                        string   `json:"op_policy_uri,omitempty"`
	OpTosURI                                           string   `json:"op_tos_uri,omitempty"`
	RevocationEndpoint                                 string   `json:"revocation_endpoint,omitempty"`
	RevocationEndpointAuthMethodsSupported             []string `json:"revocation_endpoint_auth_methods_supported,omitempty"`
	RevocationEndpointAuthSigningAlgValuesSupported    []string `json:"revocation_endpoint_auth_signing_alg_values_supported,omitempty"`
	IntrospectionEndpoint                              string   `json:"introspection_endpoint,omitempty"`
	IntrospectionEndpointAuthMethodsSupported          []string `json:"introspection_endpoint_auth_methods_supported,omitempty"`
	IntrospectionEndpointAuthSigningAlgValuesSupported []string `json:"introspection_endpoint_auth_signing_alg_values_supported,omitempty"`
	CodeChallengeMethodsSupported                      []string `json:"code_challenge_methods_supported,omitempty"`
	SignedMetadata                                     string   `json:"signed_metadata,omitempty"`
	DeviceAuthorizationEndpoint                        string   `json:"device_authorization_endpoint,omitempty"`
	TLSClientCertificateBoundAccessTokens              bool     `json:"tls_client_certificate_bound_access_tokens,omitempty"`
	MTLSEndpointAliases                                []string `json:"mtls_endpoint_aliases,omitempty"`
	UserinfoEndpoint                                   string   `json:"userinfo_endpoint,omitempty"`
	AcrValuesSupported                                 []string `json:"acr_values_supported,omitempty"`
	SubjectTypesSupported                              []string `json:"subject_types_supported,omitempty"`
	IDTokenSigningAlgValuesSupported                   []string `json:"id_token_signing_alg_values_supported,omitempty"`
	IDTokenEncryptionAlgValuesSupported                []string `json:"id_token_encryption_alg_values_supported,omitempty"`
	IDTokenEncryptionEncValuesSupported                []string `json:"id_token_encryption_enc_values_supported,omitempty"`
	UserinfoSigningAlgValuesSupported                  []string `json:"userinfo_signing_alg_values_supported,omitempty"`
	UserinfoEncryptionAlgValuesSupported               []string `json:"userinfo_encryption_alg_values_supported,omitempty"`
	UserinfoEncryptionEncValuesSupported               []string `json:"userinfo_encryption_enc_values_supported,omitempty"`
	RequestObjectSigningAlgValuesSupported             []string `json:"request_object_signing_alg_values_supported,omitempty"`
	RequestObjectEncryptionAlgValuesSupported          []string `json:"request_object_encryption_alg_values_supported,omitempty"`
	RequestObjectEncryptionEncValuesSupported          []string `json:"request_object_encryption_enc_values_supported,omitempty"`
	DisplayValuesSupported                             []string `json:"display_values_supported,omitempty"`
	ClaimTypesSupported                                []string `json:"claim_types_supported,omitempty"`
	ClaimsSupported                                    []string `json:"claims_supported,omitempty"`
	ClaimsLocalesSupported                             []string `json:"claims_locales_supported,omitempty"`
	ClaimsParameterSupported                           bool     `json:"claims_parameter_supported"`
	RequestParameterSupported                          bool     `json:"request_parameter_supported"`
	RequestUriParameterSupported                       bool     `json:"request_uri_parameter_supported"`
	RequireRequestUriRegistration                      bool     `json:"require_request_uri_registration"`
	RequireSignedRequestObject                         bool     `json:"require_signed_request_object"`
}

func relative(u *url.URL, p string) *url.URL {
	parsed, err := url.Parse(p)
	if err != nil {
		panic(err)
	}
	return u.ResolveReference(parsed)
}

func documentForPrefix(prefix *url.URL) discoveryDocument {
	opPrefix := relative(prefix, "."+httputil.PublicAuthPathPrefix)

	return discoveryDocument{
		Issuer:                            prefix.String(),
		AuthorizationEndpoint:             relative(opPrefix, "./auth").String(),
		TokenEndpoint:                     relative(opPrefix, "./token").String(),
		RegistrationEndpoint:              relative(opPrefix, "./clients/register").String(),
		JWKSUri:                           relative(prefix, "./.well-known/jwks.json").String(),
		ResponseTypesSupported:            []string{"code"},
		ResponseModesSupported:            []string{"query", "fragment", "form_post"},
		GrantTypesSupported:               []string{"authorization_code", "refresh_token"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic", "none"},
		CodeChallengeMethodsSupported:     []string{"none", "S256"},
		UserinfoEndpoint:                  relative(opPrefix, "./userinfo").String(),
		SubjectTypesSupported:             []string{"public"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		UserinfoSigningAlgValuesSupported: []string{"RS256"},
		DisplayValuesSupported:            []string{"page"},
		ClaimsSupported:                   []string{"openid"},
		ClaimsParameterSupported:          false,
		RequestParameterSupported:         false,
		RequestUriParameterSupported:      false,
		RequireRequestUriRegistration:     false,
		RequireSignedRequestObject:        false,
	}
}

func OpenIDConfiguration(req *http.Request) util.JSONResponse {
	// TODO: proper detection and/or issuer from configuration
	prefix := relative(req.URL, "../../")
	prefix.Host = req.Host
	if req.TLS == nil {
		prefix.Scheme = "http"
	} else {
		prefix.Scheme = "https"
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: documentForPrefix(prefix),
	}
}
