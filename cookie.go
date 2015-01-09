package crowd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CookieConfig struct {
	XMLName struct{} `xml:"cookie-config"`
	Domain  string   `xml:"domain"`
	Secure  bool     `xml:"secure"`
	Name    string   `xml:"name"`
}

// Get setting for SSO cookie
func (c *Crowd) GetCookieConfig() (CookieConfig, error) {
	cc := CookieConfig{}

	client := http.Client{Jar: c.cookies}
	req, err := http.NewRequest("GET", c.url+"rest/usermanagement/1/config/cookie", nil)
	if err != nil {
		return cc, err
	}
	req.SetBasicAuth(c.user, c.passwd)
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml")
	resp, err := client.Do(req)
	if err != nil {
		return cc, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return cc, fmt.Errorf("Request failed: %s\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cc, err
	}

	err = xml.Unmarshal(body, &cc)
	if err != nil {
		return cc, err
	}

	return cc, nil
}