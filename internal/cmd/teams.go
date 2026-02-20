package cmd

import (
	"context"
	"errors"

	"github.com/yborunov/umami-cli/internal/client"
	"github.com/yborunov/umami-cli/internal/out"
)

type TeamsCmd struct {
	List     TeamsListCmd     `cmd:"" help:"List teams"`
	Websites TeamsWebsitesCmd `cmd:"" help:"List websites for a team"`
}

type TeamsListCmd struct{}

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type teamsListResponse struct {
	Data []Team `json:"data"`
}

func (c *TeamsListCmd) Run(ctx *Context) error {
	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	var resp teamsListResponse
	_, err = api.Do(context.Background(), "GET", "/teams", nil, &resp, true)
	if err != nil {
		return err
	}

	if ctx.JSON {
		return out.PrintJSON(resp.Data)
	}

	if len(resp.Data) == 0 {
		out.Printf("No teams found.\n")
		return nil
	}

	for _, t := range resp.Data {
		out.Printf("%s\t%s\n", t.ID, t.Name)
	}
	return nil
}

type TeamsWebsitesCmd struct {
	TeamID string `arg:"" name:"team-id" help:"Team ID"`
}

func (c *TeamsWebsitesCmd) Run(ctx *Context) error {
	if c.TeamID == "" {
		return errors.New("team-id is required")
	}

	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	var resp websitesListResponse
	path := "/teams/" + c.TeamID + "/websites"
	_, err = api.Do(context.Background(), "GET", path, nil, &resp, true)
	if err != nil {
		return err
	}

	if ctx.JSON {
		return out.PrintJSON(resp.Data)
	}

	if len(resp.Data) == 0 {
		out.Printf("No websites found for team %s.\n", c.TeamID)
		return nil
	}

	for _, w := range resp.Data {
		out.Printf("%s\t%s\t%s\n", w.ID, w.Name, w.Domain)
	}
	return nil
}
