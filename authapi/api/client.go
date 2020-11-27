package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/ory/fosite"
	"golang.org/x/text/language"
)

type ApplicationType string

const (
	WebApplicationType    ApplicationType = "web"
	NativeApplicationType                 = "native"
)

type ResponseType string

const (
	CodeResponseType ResponseType = "code"
)

type GrantType string

const (
	AuthorizationCodeGrantType GrantType = "authorization_code"
	RefreshTokenGrantType                = "refresh_token"
)

type SubjectType string

const (
	PublicSubjectType   SubjectType = "public"
	PairwiseSubjectType             = "pairwise"
)

type I18nURI map[language.Tag]*url.URL
type I18nString map[language.Tag]string

func langTagToMetadata(key string, tag language.Tag) string {
	if tag == language.Und {
		return key
	}
	return fmt.Sprintf("%s#%s", key, tag.String())
}

func parseMetadata(metadata string) (string, language.Tag, error) {
	parts := strings.SplitN(metadata, "#", 2)
	meta := parts[0]
	tag := language.Und
	if len(parts) == 2 {
		var err error
		tag, err = language.Parse(parts[1])
		if err != nil {
			return "", tag, err
		}
	}

	return meta, tag, nil
}

type Client struct {
	ID              string
	RedirectURIs    []*url.URL
	ResponseTypes   []ResponseType
	GrantTypes      []GrantType
	ApplicationType ApplicationType
	Contacts        []string
	ClientName      I18nString
	LogoURI         I18nURI
	ClientURI       I18nURI
	PolicyURI       I18nURI
	TOS_URI         I18nURI
	JWKS_URI        *url.URL
	SubjectType     SubjectType
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetHashedSecret() []byte {
	return []byte{}
}

func (c *Client) GetRedirectURIs() []string {
	uris := []string{}
	for _, uri := range c.RedirectURIs {
		uris = append(uris, uri.String())
	}
	return uris
}

func (c *Client) GetGrantTypes() fosite.Arguments {
	types := []string{}
	for _, gt := range c.GrantTypes {
		types = append(types, string(gt))
	}
	return types
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	types := []string{}
	for _, rt := range c.ResponseTypes {
		types = append(types, string(rt))
	}
	return types
}

func (c *Client) GetScopes() fosite.Arguments {
	return []string{
		"openid",
	}
}

func (c *Client) IsPublic() bool {
	return true
}

func (c *Client) GetAudience() fosite.Arguments {
	return []string{
		c.GetID(),
	}
}

func (c *Client) GetResponseModes() []fosite.ResponseModeType {
	return []fosite.ResponseModeType{
		fosite.ResponseModeDefault,
		fosite.ResponseModeFormPost,
		fosite.ResponseModeFragment,
		fosite.ResponseModeQuery,
	}
}

func (uris I18nURI) Populate(prefix string, output map[string]interface{}) {
	for language, uri := range uris {
		output[langTagToMetadata(prefix, language)] = uri.String()
	}
}

func (strings I18nString) Populate(prefix string, output map[string]interface{}) {
	for language, value := range strings {
		output[langTagToMetadata(prefix, language)] = value
	}
}

func NewClient(id string) *Client {
	return &Client{
		ID:              id,
		RedirectURIs:    []*url.URL{},
		ResponseTypes:   []ResponseType{},
		GrantTypes:      []GrantType{},
		ApplicationType: WebApplicationType,
		Contacts:        []string{},
		ClientName:      I18nString{},
		LogoURI:         I18nURI{},
		ClientURI:       I18nURI{},
		PolicyURI:       I18nURI{},
		TOS_URI:         I18nURI{},
		JWKS_URI:        nil,
		SubjectType:     PublicSubjectType,
	}
}

func (c *Client) Serialize() map[string]interface{} {
	out := map[string]interface{}{
		"client_id":        c.ID,
		"scope":            c.GetScopes(),
		"grant_types":      c.GetGrantTypes(),
		"application_type": c.ApplicationType,
		"redirect_uris":    c.GetRedirectURIs(),
		"response_types":   c.GetResponseTypes(),
	}

	if len(c.Contacts) != 0 {
		out["contacts"] = c.Contacts
	}

	c.ClientName.Populate("client_name", out)
	c.ClientURI.Populate("client_uri", out)
	c.LogoURI.Populate("logo_uri", out)
	c.TOS_URI.Populate("tos_uri", out)
	c.PolicyURI.Populate("policy_uri", out)
	if c.JWKS_URI != nil {
		out["jwks_uri"] = c.JWKS_URI.String()
	}

	return out
}

func (c *Client) Fill(values map[string]interface{}) error {
	for key, value := range values {
		key, tag, err := parseMetadata(key)
		if err != nil {
			return err
		}

		switch key {
		case "scope":
			continue

		case "grant_types":
			if tag != language.Und {
				return fmt.Errorf("unexpected language tag %s on metadata %s", tag.String(), key)
			}

			values, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("invalid type for metadata %s", key)
			}

			for _, value := range values {
				value, ok := value.(string)
				if !ok {
					return fmt.Errorf("invalid type for metadata %s", key)
				}

				switch value {
				case string(AuthorizationCodeGrantType):
					c.GrantTypes = append(c.GrantTypes, AuthorizationCodeGrantType)
				case string(RefreshTokenGrantType):
					c.GrantTypes = append(c.GrantTypes, RefreshTokenGrantType)
				default:
					return fmt.Errorf("unsupported grant type %s", value)
				}
			}

		case "response_types":
			if tag != language.Und {
				return fmt.Errorf("unexpected language tag %s on metadata %s", tag.String(), key)
			}

			values, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("invalid type for metadata %s", key)
			}

			for _, value := range values {
				value, ok := value.(string)
				if !ok {
					return fmt.Errorf("invalid type for metadata %s", key)
				}

				switch value {
				case string(CodeResponseType):
					c.ResponseTypes = append(c.ResponseTypes, CodeResponseType)
				default:
					return fmt.Errorf("unsupported response type type %s", value)
				}
			}

		case "application_type":
			if tag != language.Und {
				return fmt.Errorf("unexpected language tag %s on metadata %s", tag.String(), key)
			}

			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s", key)
			}

			switch value {
			case string(WebApplicationType):
				c.ApplicationType = WebApplicationType
			case string(NativeApplicationType):
				c.ApplicationType = NativeApplicationType
			default:
				return fmt.Errorf("unsupported application type %s", value)
			}

		case "redirect_uris":
			if tag != language.Und {
				return fmt.Errorf("unexpected language tag %s on metadata %s", tag.String(), key)
			}

			values, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("invalid type for metadata %s", key)
			}

			for _, value := range values {
				value, ok := value.(string)
				if !ok {
					return fmt.Errorf("invalid type for metadata %s", key)
				}

				uri, err := url.Parse(value)
				if err != nil {
					return fmt.Errorf("invalid redirect uri %s: %v", value, err)
				}
				c.RedirectURIs = append(c.RedirectURIs, uri)
			}

		case "contacts":
			if tag != language.Und {
				return fmt.Errorf("unexpected language tag %s on metadata %s", tag.String(), key)
			}

			values, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("invalid type for metadata %s", key)
			}

			for _, value := range values {
				value, ok := value.(string)
				if !ok {
					return fmt.Errorf("invalid type for metadata %s", key)
				}

				c.Contacts = append(c.Contacts, value)
			}

		case "client_name":
			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s (%s)", key, tag.String())
			}

			c.ClientName[tag] = value

		case "client_uri":
			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s (%s)", key, tag.String())
			}

			uri, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf("invalid client uri (%s) %s: %v", tag.String(), value, err)
			}

			c.ClientURI[tag] = uri

		case "logo_uri":
			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s (%s)", key, tag.String())
			}

			uri, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf("invalid logo uri (%s) %s: %v", tag.String(), value, err)
			}

			c.LogoURI[tag] = uri

		case "tos_uri":
			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s (%s)", key, tag.String())
			}

			uri, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf("invalid tos uri (%s) %s: %v", tag.String(), value, err)
			}

			c.TOS_URI[tag] = uri

		case "policy_uri":
			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s (%s)", key, tag.String())
			}

			uri, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf("invalid policy uri (%s) %s: %v", tag.String(), value, err)
			}

			c.PolicyURI[tag] = uri

		case "jwks_uri":
			if tag != language.Und {
				return fmt.Errorf("unexpected language tag %s on metadata %s", tag.String(), key)
			}

			value, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid type for metadata %s", key)
			}

			uri, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf("invalid jwks uri %s: %v", value, err)
			}

			c.JWKS_URI = uri

		default:
			return fmt.Errorf("unsupported metadata %s", key)
		}
	}

	return nil
}

func (c *Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Serialize())
}

func (c *Client) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	c.Fill(m)
	return nil
}
