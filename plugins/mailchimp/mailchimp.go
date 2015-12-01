package mailchimp

import (
	"time"

	"github.com/influxdb/telegraf/plugins"
)

type MailChimp struct {
	api *ChimpAPI

	ApiKey string
}

var sampleConfig = `
  # MailChimp API key
  api_key = ""
`

func (m *MailChimp) SampleConfig() string {
	return sampleConfig
}

func (m *MailChimp) Description() string {
	return "Gathers metrics from the /reports MailChimp API"
}

func (m *MailChimp) Gather(acc plugins.Accumulator) error {
	if m.api == nil {
		m.api = NewChimp(m.ApiKey)
	}
	m.api.Debug = false

	reports, err := m.api.GetReports()
	if err != nil {
		return err
	}
	now := time.Now()

	for _, report := range reports.Reports {
		tags := make(map[string]string)
		tags["id"] = report.ID
		tags["campaign_title"] = report.CampaignTitle
		acc.Add("emails_sent", report.EmailsSent, tags, now)
		acc.Add("abuse_reports", report.AbuseReports, tags, now)
		acc.Add("unsubscribed", report.Unsubscribed, tags, now)
		acc.Add("hard_bounces", report.Bounces.HardBounces, tags, now)
		acc.Add("soft_bounces", report.Bounces.SoftBounces, tags, now)
		acc.Add("syntax_errors", report.Bounces.SyntaxErrors, tags, now)
		acc.Add("forwards_count", report.Forwards.ForwardsCount, tags, now)
		acc.Add("forwards_opens", report.Forwards.ForwardsOpens, tags, now)
		acc.Add("opens_total", report.Opens.OpensTotal, tags, now)
		acc.Add("unique_opens", report.Opens.UniqueOpens, tags, now)
		acc.Add("open_rate", report.Opens.OpenRate, tags, now)
		acc.Add("last_open", report.Opens.LastOpen, tags, now)
		acc.Add("clicks_total", report.Clicks.ClicksTotal, tags, now)
		acc.Add("unique_clicks", report.Clicks.UniqueClicks, tags, now)
		acc.Add("unique_subscriber_clicks", report.Clicks.UniqueSubscriberClicks, tags, now)
		acc.Add("click_rate", report.Clicks.ClickRate, tags, now)
		acc.Add("last_click", report.Clicks.LastClick, tags, now)
		acc.Add("facebook_recipient_likes", report.FacebookLikes.RecipientLikes, tags, now)
		acc.Add("facebook_unique_likes", report.FacebookLikes.UniqueLikes, tags, now)
		acc.Add("facebook_likes", report.FacebookLikes.FacebookLikes, tags, now)
		acc.Add("industry_type", report.IndustryStats.Type, tags, now)
		acc.Add("industry_open_rate", report.IndustryStats.OpenRate, tags, now)
		acc.Add("industry_click_rate", report.IndustryStats.ClickRate, tags, now)
		acc.Add("industry_bounce_rate", report.IndustryStats.BounceRate, tags, now)
		acc.Add("industry_unopen_rate", report.IndustryStats.UnopenRate, tags, now)
		acc.Add("industry_unsub_rate", report.IndustryStats.UnsubRate, tags, now)
		acc.Add("industry_abuse_rate", report.IndustryStats.AbuseRate, tags, now)
		acc.Add("list_stats_sub_rate", report.ListStats.SubRate, tags, now)
		acc.Add("list_stats_unsub_rate", report.ListStats.UnsubRate, tags, now)
		acc.Add("list_stats_open_rate", report.ListStats.OpenRate, tags, now)
		acc.Add("list_stats_click_rate", report.ListStats.ClickRate, tags, now)

		// TODO(cam) parse and add TimeSeries data here too
	}

	return nil
}

func init() {
	plugins.Add("mailchimp", func() plugins.Plugin {
		return &MailChimp{}
	})
}
