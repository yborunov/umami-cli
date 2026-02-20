package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/yborunov/umami-cli/internal/client"
	"github.com/yborunov/umami-cli/internal/out"
)

type AnalyticsCmd struct {
	Active          AnalyticsActiveCmd          `cmd:"" help:"Active users"`
	EventsSeries    AnalyticsEventsSeriesCmd    `cmd:"" help:"Event series"`
	Metrics         AnalyticsMetricsCmd         `cmd:"" help:"Metrics"`
	MetricsExpanded AnalyticsMetricsExpandedCmd `cmd:"" help:"Expanded metrics"`
	Pageviews       AnalyticsPageviewsCmd       `cmd:"" help:"Pageviews"`
	Stats           AnalyticsStatsCmd           `cmd:"" help:"Summary stats"`
}

type TimeRange struct {
	StartAt int64 `help:"Start timestamp (ms since epoch)"`
	EndAt   int64 `help:"End timestamp (ms since epoch)"`
}

type Filters struct {
	Path       string `help:"Filter by URL path"`
	Referrer   string `help:"Filter by referrer"`
	Title      string `help:"Filter by page title"`
	Query      string `help:"Filter by query parameter"`
	Browser    string `help:"Filter by browser"`
	OS         string `help:"Filter by operating system"`
	Device     string `help:"Filter by device"`
	Country    string `help:"Filter by country"`
	Region     string `help:"Filter by region"`
	City       string `help:"Filter by city"`
	Hostname   string `help:"Filter by hostname"`
	Tag        string `help:"Filter by tag"`
	DistinctID string `help:"Filter by distinct ID"`
	Segment    string `help:"Filter by segment UUID"`
	Cohort     string `help:"Filter by cohort UUID"`
}

type AnalyticsActiveCmd struct {
	WebsiteID string `arg:"" name:"website-id" help:"Website ID"`
}

func (c *AnalyticsActiveCmd) Run(ctx *Context) error {
	if c.WebsiteID == "" {
		return errors.New("website-id is required")
	}

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	var resp any
	path := fmt.Sprintf("/websites/%s/active", c.WebsiteID)
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}
	return out.PrintJSON(resp)
}

type AnalyticsEventsSeriesCmd struct {
	WebsiteID string `arg:"" name:"website-id" help:"Website ID"`
	TimeRange
	Unit     string `help:"Time unit (year|month|day|hour|minute)"`
	Timezone string `help:"Timezone (e.g. America/Los_Angeles)"`
	Filters
}

func (c *AnalyticsEventsSeriesCmd) Run(ctx *Context) error {
	if err := validateWebsiteID(c.WebsiteID); err != nil {
		return err
	}
	startAt, endAt := normalizeRange(c.StartAt, c.EndAt)

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	q := buildQuery(startAt, endAt, c.Unit, c.Timezone, c.Filters, 0, 0, "")
	path := withQuery(fmt.Sprintf("/websites/%s/events/series", c.WebsiteID), q)

	var resp any
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}
	return out.PrintJSON(resp)
}

type AnalyticsMetricsCmd struct {
	WebsiteID string `arg:"" name:"website-id" help:"Website ID"`
	TimeRange
	Type   string `help:"Metric type (path|entry|exit|title|query|referrer|channel|domain|country|region|city|browser|os|device|language|screen|event|hostname|tag|distinctId)"`
	Limit  int    `help:"Number of rows returned (default 500)"`
	Offset int    `help:"Number of rows to skip (default 0)"`
	Filters
}

func (c *AnalyticsMetricsCmd) Run(ctx *Context) error {
	if err := validateWebsiteID(c.WebsiteID); err != nil {
		return err
	}
	if c.Type == "" {
		return errors.New("type is required")
	}
	startAt, endAt := normalizeRange(c.StartAt, c.EndAt)

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	q := buildQuery(startAt, endAt, "", "", c.Filters, c.Limit, c.Offset, c.Type)
	path := withQuery(fmt.Sprintf("/websites/%s/metrics", c.WebsiteID), q)

	var resp any
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}
	return out.PrintJSON(resp)
}

type AnalyticsMetricsExpandedCmd struct {
	WebsiteID string `arg:"" name:"website-id" help:"Website ID"`
	TimeRange
	Type   string `help:"Metric type (path|entry|exit|title|query|referrer|channel|domain|country|region|city|browser|os|device|language|screen|event|hostname|tag|distinctId)"`
	Limit  int    `help:"Number of rows returned (default 500)"`
	Offset int    `help:"Number of rows to skip (default 0)"`
	Filters
}

func (c *AnalyticsMetricsExpandedCmd) Run(ctx *Context) error {
	if err := validateWebsiteID(c.WebsiteID); err != nil {
		return err
	}
	if c.Type == "" {
		return errors.New("type is required")
	}
	startAt, endAt := normalizeRange(c.StartAt, c.EndAt)

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	q := buildQuery(startAt, endAt, "", "", c.Filters, c.Limit, c.Offset, c.Type)
	path := withQuery(fmt.Sprintf("/websites/%s/metrics/expanded", c.WebsiteID), q)

	var resp any
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}
	return out.PrintJSON(resp)
}

type AnalyticsPageviewsCmd struct {
	WebsiteID string `arg:"" name:"website-id" help:"Website ID"`
	TimeRange
	Unit     string `help:"Time unit (year|month|day|hour|minute)"`
	Timezone string `help:"Timezone (e.g. America/Los_Angeles)"`
	Compare  string `help:"Comparison value (prev|yoy)"`
	Filters
}

func (c *AnalyticsPageviewsCmd) Run(ctx *Context) error {
	if err := validateWebsiteID(c.WebsiteID); err != nil {
		return err
	}
	startAt, endAt := normalizeRange(c.StartAt, c.EndAt)

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	q := buildQuery(startAt, endAt, c.Unit, c.Timezone, c.Filters, 0, 0, "")
	if c.Compare != "" {
		q.Set("compare", c.Compare)
	}
	path := withQuery(fmt.Sprintf("/websites/%s/pageviews", c.WebsiteID), q)

	var resp any
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}
	return out.PrintJSON(resp)
}

type AnalyticsStatsCmd struct {
	WebsiteID string `arg:"" name:"website-id" help:"Website ID"`
	TimeRange
	Filters
}

func (c *AnalyticsStatsCmd) Run(ctx *Context) error {
	if err := validateWebsiteID(c.WebsiteID); err != nil {
		return err
	}
	startAt, endAt := normalizeRange(c.StartAt, c.EndAt)

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	q := buildQuery(startAt, endAt, "", "", c.Filters, 0, 0, "")
	path := withQuery(fmt.Sprintf("/websites/%s/stats", c.WebsiteID), q)

	var resp any
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}
	return out.PrintJSON(resp)
}

func validateWebsiteID(websiteID string) error {
	if websiteID == "" {
		return errors.New("website-id is required")
	}
	return nil
}

func normalizeRange(startAt, endAt int64) (int64, int64) {
	if startAt != 0 && endAt != 0 {
		return startAt, endAt
	}
	now := time.Now().UTC()
	end := now.UnixMilli()
	start := now.Add(-24 * time.Hour).UnixMilli()
	return start, end
}

func buildQuery(startAt, endAt int64, unit, timezone string, filters Filters, limit, offset int, metricType string) url.Values {
	q := url.Values{}
	if startAt != 0 {
		q.Set("startAt", strconv.FormatInt(startAt, 10))
	}
	if endAt != 0 {
		q.Set("endAt", strconv.FormatInt(endAt, 10))
	}
	if unit != "" {
		q.Set("unit", unit)
	}
	if timezone != "" {
		q.Set("timezone", timezone)
	}
	if metricType != "" {
		q.Set("type", metricType)
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}

	addFilter := func(key, value string) {
		if value != "" {
			q.Set(key, value)
		}
	}

	addFilter("path", filters.Path)
	addFilter("referrer", filters.Referrer)
	addFilter("title", filters.Title)
	addFilter("query", filters.Query)
	addFilter("browser", filters.Browser)
	addFilter("os", filters.OS)
	addFilter("device", filters.Device)
	addFilter("country", filters.Country)
	addFilter("region", filters.Region)
	addFilter("city", filters.City)
	addFilter("hostname", filters.Hostname)
	addFilter("tag", filters.Tag)
	addFilter("distinctId", filters.DistinctID)
	addFilter("segment", filters.Segment)
	addFilter("cohort", filters.Cohort)

	return q
}

func withQuery(path string, q url.Values) string {
	if len(q) == 0 {
		return path
	}
	return path + "?" + q.Encode()
}
