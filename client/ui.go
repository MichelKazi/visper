package client

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (c *Client) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if chatView, err := g.SetView("chatView", 0, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		chatView.Title = fmt.Sprintf("Chatting on %s as %s", c.ServerAddr, c.Username)
		chatView.Autoscroll = true
		chatView.Wrap = true
	}

	if inputView, err := g.SetView("inputView", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		inputView.Title = "Enter your message here"
		inputView.Editable = true
		inputView.Wrap = true

    if _, err := g.SetCurrentView("inputView"); err != nil {
      return err
    }
	}

	return nil
}
