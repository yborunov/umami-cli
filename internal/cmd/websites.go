package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yborunov/umami-cli/internal/client"
	"github.com/yborunov/umami-cli/internal/out"
)

type WebsitesCmd struct {
	List WebsitesListCmd `cmd:"" help:"List websites"`
}

type WebsitesListCmd struct{}

type Website struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type websitesListResponse struct {
	Data []Website `json:"data"`
}

func (c *WebsitesListCmd) Run(ctx *Context) error {
	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	var resp websitesListResponse
	if debugEnabled() {
		out.Printf("debug: websites request method=GET path=/websites endpoint=%s\n", ctx.Config.Endpoint)
		status, body, err := api.DoRaw(context.Background(), "GET", "/websites", nil, true)
		if err != nil {
			out.Printf("debug: websites response status=%d error=%v\n", status, err)
			if len(body) > 0 {
				out.Printf("debug: websites response body=%s\n", truncateBody(body))
			}
			return err
		}
		out.Printf("debug: websites response status=%d\n", status)
		if len(body) > 0 {
			out.Printf("debug: websites response body=%s\n", truncateBody(body))
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("failed to parse websites response: %w", err)
		}
	} else {
		_, err = api.Do(context.Background(), "GET", "/websites", nil, &resp, true)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if ctx.JSON {
		return out.PrintJSON(resp.Data)
	}

	if len(resp.Data) == 0 {
		out.Printf("No websites found.\n")
		return nil
	}

	for _, w := range resp.Data {
		out.Printf("%s\t%s\t%s\n", w.ID, w.Name, w.Domain)
	}
	return nil
}

func truncateBody(body []byte) string {
	const max = 2048
	if len(body) <= max {
		return string(body)
	}
	return string(body[:max]) + "...(truncated)"
}
