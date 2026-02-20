package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/yborunov/umami-cli/internal/client"
	"github.com/yborunov/umami-cli/internal/out"
)

type AuthCmd struct {
	Login  AuthLoginCmd  `cmd:"" help:"Login with username and password"`
	Verify AuthVerifyCmd `cmd:"" help:"Verify stored token"`
}

type AuthLoginCmd struct {
	Username string `help:"Umami username" env:"UMAMI_USERNAME"`
	Password string `help:"Umami password" env:"UMAMI_PASSWORD"`
}

type loginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *AuthLoginCmd) Run(ctx *Context) error {
	if c.Username == "" || c.Password == "" {
		return errors.New("username and password are required")
	}

	api, err := client.New(ctx.Config.Endpoint, "")
	if err != nil {
		return err
	}

	req := loginRequest{Username: c.Username, Password: c.Password}
	resp := loginResponse{}
	_, err = api.Do(context.Background(), "POST", "/auth/login", req, &resp, false)
	if err != nil {
		return err
	}

	ctx.Config.Token = resp.Token
	if err := ctx.Config.Save(); err != nil {
		return err
	}

	if ctx.JSON {
		return out.PrintJSON(resp)
	}

	out.Printf("Logged in as %s (id %s). Token saved.\n", resp.User.Username, resp.User.ID)
	return nil
}

type AuthVerifyCmd struct{}

type verifyResponse struct {
	Valid bool   `json:"valid"`
	User  string `json:"user"`
	Role  string `json:"role"`
}

func (c *AuthVerifyCmd) Run(ctx *Context) error {
	api, err := client.New(ctx.Config.Endpoint, ctx.Config.Token)
	if err != nil {
		return err
	}

	var resp any
	if debugEnabled() {
		out.Printf("debug: verify request method=POST path=/auth/verify endpoint=%s\n", ctx.Config.Endpoint)
	}
	status, err := api.Do(context.Background(), "POST", "/auth/verify", nil, &resp, true)
	if debugEnabled() {
		if err != nil {
			out.Printf("debug: verify response status=%d error=%v\n", status, err)
		} else {
			out.Printf("debug: verify response status=%d\n", status)
		}
	}
	if err != nil {
		return err
	}

	if ctx.JSON {
		return out.PrintJSON(resp)
	}

	out.Printf("Token verified at %s.\n", time.Now().Format(time.RFC3339))
	return nil
}
