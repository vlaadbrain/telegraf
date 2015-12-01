package mailchimp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const (
	reports_summary_endpoint string = "/3.0/reports"
)

var mailchimp_datacenter = regexp.MustCompile("[a-z]+[0-9]+$")

type ChimpAPI struct {
	Transport http.RoundTripper
	endpoint  string
	Debug     bool
}

func NewChimp(apiKey string) *ChimpAPI {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = fmt.Sprintf("%s.api.mailchimp.com", mailchimp_datacenter.FindString(apiKey))
	u.Path = reports_summary_endpoint
	u.User = url.UserPassword("", apiKey)
	return &ChimpAPI{endpoint: u.String()}
}

type APIError struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Name   string `json:"name"`
	Err    string `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v: %v", e.Code, e.Err)
}

func (a *ChimpAPI) GetReports() (ReportsResponse, error) {
	var response ReportsResponse
	rawjson, err := runChimp(a)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(rawjson, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func runChimp(api *ChimpAPI) ([]byte, error) {
	if api.Debug {
		log.Printf("Request URL: %s", api.endpoint)
	}
	client := &http.Client{Transport: api.Transport}
	resp, err := client.Get(api.endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if api.Debug {
		log.Printf("Response Body:%s", string(body))
	}

	// TODO(cam) check response body for error
	// if err = chimpErrorCheck(body); err != nil {
	// 	return nil, err
	// }
	return body, nil
}

type ReportsResponse struct {
	Reports    []Report `json:"reports"`
	TotalItems int      `json:"total_items"`
}

type Report struct {
	ID            string `json:"id"`
	CampaignTitle string `json:"campaign_title"`
	Type          string `json:"type"`
	EmailsSent    int    `json:"emails_sent"`
	AbuseReports  int    `json:"abuse_reports"`
	Unsubscribed  int    `json:"unsubscribed"`
	SendTime      string `json:"send_time"`

	TimeSeries    []TimeSerie
	Bounces       Bounces       `json:"bounces"`
	Forwards      Forwards      `json:"forwards"`
	Opens         Opens         `json:"opens"`
	Clicks        Clicks        `json:"clicks"`
	FacebookLikes FacebookLikes `json:"facebook_likes"`
	IndustryStats IndustryStats `json:"industry_stats"`
	ListStats     ListStats     `json:"list_stats"`
}

type Bounces struct {
	HardBounces  int `json:"hard_bounces"`
	SoftBounces  int `json:"soft_bounces"`
	SyntaxErrors int `json:"syntax_errors"`
}

type Forwards struct {
	ForwardsCount int `json:"forwards_count"`
	ForwardsOpens int `json:"forwards_opens"`
}

type Opens struct {
	OpensTotal  int     `json:"opens_total"`
	UniqueOpens int     `json:"unique_opens"`
	OpenRate    float64 `json:"open_rate"`
	LastOpen    string  `json:"last_open"`
}

type Clicks struct {
	ClicksTotal            int     `json:"clicks_total"`
	UniqueClicks           int     `json:"unique_clicks"`
	UniqueSubscriberClicks int     `json:"unique_subscriber_clicks"`
	ClickRate              float64 `json:"click_rate"`
	LastClick              string  `json:"last_click"`
}

type FacebookLikes struct {
	RecipientLikes int `json:"recipient_likes"`
	UniqueLikes    int `json:"unique_likes"`
	FacebookLikes  int `json:"facebook_likes"`
}

type IndustryStats struct {
	Type       string  `json:"type"`
	OpenRate   float64 `json:"open_rate"`
	ClickRate  float64 `json:"click_rate"`
	BounceRate float64 `json:"bounce_rate"`
	UnopenRate float64 `json:"unopen_rate"`
	UnsubRate  float64 `json:"unsub_rate"`
	AbuseRate  float64 `json:"abuse_rate"`
}

type ListStats struct {
	SubRate   float64 `json:"sub_rate"`
	UnsubRate float64 `json:"unsub_rate"`
	OpenRate  float64 `json:"open_rate"`
	ClickRate float64 `json:"click_rate"`
}

type TimeSerie struct {
	TimeStamp       string `json:"timestamp"`
	EmailsSent      int    `json:"emails_sent"`
	UniqueOpens     int    `json:"unique_opens"`
	RecipientsClick int    `json:"recipients_click"`
}
