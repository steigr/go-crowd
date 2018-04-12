package crowd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	MATCH_MODE_EXACTLY_MATCHES = "EXACTLY_MATCHES"

	RESTRICTION_TYPE_PROPERTY_SEARCH_RESTRICTION = "property-search-restriction"

	PROPERTY_STRING = "STRING"
)

type SearchProperty struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Group represents a group in Crowd
type Search struct {
	MatchMode       string         `json:"match-mode"`
	Value           string         `json:"value"`
	RestrictionType string         `json:"restriction-type"`
	Property        SearchProperty `json:"property"`
}

type SearchResultLink struct {
	Reference string `json:"href"`
	Relation  string `json:"rel"`
}

type SearchResults struct {
	Link SearchResultLink `json:"link"`
	Name string           `json:"name"`
}

type GroupSearchResult struct {
	Expand string          `json:"expand"`
	Groups []SearchResults `json:"groups"`
}

type UserSearchResult struct {
	Expand string          `json:"expand"`
	Users  []SearchResults `json:"users"`
}

type SearchResult struct {
	Group GroupSearchResult
	User  UserSearchResult
}

// GetUser retrieves user information
func (c *Crowd) Search(value, match_mode, property, property_type string) (SearchResult, error) {
	r := SearchResult{}
	s := Search{
		MatchMode:       match_mode,
		Value:           value,
		RestrictionType: RESTRICTION_TYPE_PROPERTY_SEARCH_RESTRICTION,
		Property: SearchProperty{
			Name: property,
			Type: property_type,
		},
	}

	v := url.Values{}
	entityType := "user"
	v.Set("entity-type", entityType)
	v.Set("max-results", "1000")
	v.Set("start-index", "0")

	url := c.url + "/rest/usermanagement/1/search?" + v.Encode()
	c.Client.Jar = c.cookies
	b, _ := json.Marshal(&s)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return r, err
	}
	req.SetBasicAuth(c.user, c.passwd)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	// case 404:
	// 	return u, fmt.Errorf("user not found")
	case 200:
		// fall through switch without returning
	default:
		return r, fmt.Errorf("request failed: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return r, err
	}
	if entityType == "user" {
		err = json.Unmarshal(body, &r.User)
		if err != nil {
			return r, err
		}
	} else if entityType == "group" {
		err = json.Unmarshal(body, &r.Group)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}
