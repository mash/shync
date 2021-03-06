package shopify

import (
	"context"
	"fmt"
	"net/url"

	"github.com/chromedp/chromedp"
	"github.com/mash/shync/log"
)

type Client struct {
	ctx    context.Context
	cancel func()

	Location string
}

// head=false means headless
func NewClient(head bool) *Client {
	ctx, cancel := ChromeContext(head)
	return &Client{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Client) Close() error {
	c.cancel()
	return nil
}

// head=false means headless
func ChromeContext(head bool) (context.Context, func()) {
	var ctx context.Context
	var cancel func()
	if !head {
		// headless
		ctx, cancel = chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Debugfn))
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-extensions", false),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.Flag("mute-audio", false),
		)

		allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)

		ctx_, cancel2 := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Debugfn))
		ctx = ctx_
		cancel = func() {
			cancel1()
			cancel2()
		}
	}
	return ctx, cancel
}

func (c *Client) Login(shop, username, password string) error {
	submit := `//button[@type='submit']`
	selectStore := `//a[contains(@href,'`+shop+`')]`
	emailInput := `//input[@name='account[email]']`
	passwordInput := `//input[@name='account[password]']`
	linkSettings := `//a[@href='/admin/settings']`

	actions := chromedp.Tasks{
		chromedp.Navigate(`https://accounts.shopify.com/store-login?new_store_login=true`),

		chromedp.WaitVisible(emailInput),
		chromedp.SendKeys(emailInput, username),
		chromedp.Submit(submit),

		chromedp.WaitVisible(passwordInput),
		chromedp.SendKeys(passwordInput, password),
		chromedp.Submit(submit),

		chromedp.WaitVisible(selectStore),
		chromedp.Click(selectStore),

		chromedp.WaitVisible(linkSettings),
		chromedp.Location(&c.Location),
	}
	if err := chromedp.Run(c.ctx, actions); err != nil {
		return fmt.Errorf("Login: %w", err)
	}
	return nil
}

// id must be one of Templates.
// *body will be set to the email template body.
func (c *Client) FetchEmailTemplate(id string, subject, body *string) error {
	next, err := c.URL(id)
	if err != nil {
		return fmt.Errorf("FetchEmailTemplate: %w", err)
	}

	input := `//input[@name='email_template[title]']`
	textarea := `//textarea[@name='email_template[body_html]']`
	actions := chromedp.Tasks{
		chromedp.Navigate(next),
		chromedp.WaitVisible(textarea),
		chromedp.Value(input, subject),
		chromedp.Value(textarea, body),
	}
	if err := chromedp.Run(c.ctx, actions); err != nil {
		return fmt.Errorf("FetchEmailTemplate: %w", err)
	}
	return nil
}

// id must be one of Templates.
// *body will be set to the email template body.
func (c *Client) UpdateEmailTemplate(id, subject, body string) error {
	if subject == "" && body == "" {
		log.Infof("nothing to do for %s", id)
		return nil
	}

	next, err := c.URL(id)
	if err != nil {
		return fmt.Errorf("UpdateEmailTemplate: %w", err)
	}

	input := `//input[@name='email_template[title]']`
	textarea := `//textarea[@name='email_template[body_html]']`
	form := `//form[@id='edit_email_template']`

	var currentSubject, currentBody string

	actions := chromedp.Tasks{
		chromedp.Navigate(next),
		chromedp.WaitReady(form),
		chromedp.Value(input, &currentSubject),
		chromedp.Value(textarea, &currentBody),
	}
	if err := chromedp.Run(c.ctx, actions); err != nil {
		return fmt.Errorf("UpdateEmailTemplate: %w", err)
	}

	if currentSubject == subject && currentBody == body {
		log.Infof("no change found in %s", id)
		return nil
	}

	actions = chromedp.Tasks{}
	if currentSubject != subject && subject != "" {
		log.Infof("pushing %s subject", id)
		actions = append(actions, chromedp.SetValue(input, subject))
	}
	if currentBody != body && body != "" {
		log.Infof("pushing %s body", id)
		actions = append(actions, chromedp.SetValue(textarea, body))
	}
	actions = append(actions,
		chromedp.Submit(form),
		chromedp.WaitVisible(form))
	if err := chromedp.Run(c.ctx, actions); err != nil {
		return fmt.Errorf("UpdateEmailTemplate: %w", err)
	}
	return nil
}

func (c *Client) URL(id string) (string, error) {
	path, err := url.Parse(fmt.Sprintf("/admin/email_templates/%s/edit", id))
	if err != nil {
		return "", err
	}
	base, err := url.Parse(c.Location)
	if err != nil {
		return "", err
	}
	next := base.ResolveReference(path)
	return next.String(), nil
}
