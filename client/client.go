package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type Client struct {
	ServerAddr     string
	Protocol       string
	Username       string
	MessageChan    chan string
	conn           net.Conn
	gui            *gocui.Gui
}

func (c *Client) Init() error {
	var err error
	c.conn, err = net.Dial(c.Protocol, c.ServerAddr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Start() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your name: ")
	scanner.Scan()
	c.Username = scanner.Text()
	c.conn.Write([]byte(c.Username + " has joined the chat!\n"))

	// initialize GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	c.gui = g
	g.SetManagerFunc(c.layout)

	// Handle keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		panic(err)
	}

	if err := g.SetKeybinding("inputView", gocui.KeyEnter, gocui.ModNone, c.sendMessage); err != nil {
		panic(err)
	}

	go c.readMessages()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}

	defer c.conn.Close()

	// for scanner.Scan() {
	//   h, m, _ := time.Now().Clock()
	//   timestamp := fmt.Sprintf("%d:%d", h, m)
	//   msg := fmt.Sprintf("[%s] %s: %s", timestamp, c.Username, scanner.Text()+"\n")
	//   c.conn.Write([]byte(msg))
	// }
}

func (c *Client) sendMessage(g *gocui.Gui, v *gocui.View) error {
	text := v.Buffer()
	if len(text) == 0 {
		return nil
	}
	text = strings.TrimSpace(text)
	h, m, _ := time.Now().Clock()
	timestamp := fmt.Sprintf("%d:%d", h, m)
	broadcastMsg := fmt.Sprintf("[%s] %s: %s", timestamp, c.Username, text+"\n")
  selfRenderedMsg := fmt.Sprintf("[%s] You: %s", timestamp, text+"\n")
	if _, err := c.conn.Write([]byte(broadcastMsg)); err != nil {
		return err
	}

	// render self written message ourselves since client message excluded from server broadcast
	c.gui.Update(c.renderMessageCallback(selfRenderedMsg))

	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

// dont want to immediately invoke writing to view
func (c *Client) renderMessageCallback(msg string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		v, err := g.View("chatView")
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(v, msg)
		return nil
	}
}

func (c *Client) readMessages() {
	reader := bufio.NewReader(c.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Lost connection to server.")
		}

		c.gui.Update(c.renderMessageCallback(msg))
	}
}

func (c *Client) quit(g *gocui.Gui, v *gocui.View) error {
	exitMsg := fmt.Sprintf("%s has left the chat.", c.Username)
	if _, err := c.conn.Write([]byte(exitMsg)); err != nil {
		return err
	}
	c.conn.Close()
	return gocui.ErrQuit
}
