package kvsclient

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

type Request struct {
	Action string `json:"action"`
	Domain string `json:"domain,omitempty"`
	Key    string `json:"key,omitempty"`
	SLKey  string `json:"slkey,omitempty"`
	Value  string `json:"value,omitempty"`
	MinKey string `json:"min_key,omitempty"`
	MaxKey string `json:"max_key,omitempty"`
}

type Response struct {
	Status  string   `json:"status"`
	Message string   `json:"message,omitempty"`
	Value   string   `json:"value,omitempty"`
	Values  []string `json:"values,omitempty"`
}

func NewClient(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) sendRequest(req Request) (Response, error) {
	err := c.conn.WriteJSON(req)
	if err != nil {
		return Response{}, err
	}

	var resp Response
	err = c.conn.ReadJSON(&resp)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}

func (c *Client) CreateDomain(domain string) error {
	req := Request{Action: "create_domain", Domain: domain}
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	if resp.Status != "success" {
		return fmt.Errorf("failed to create domain: %s", resp.Message)
	}
	return nil
}

func (c *Client) SetString(domain, key, value string) error {
	req := Request{Action: "set_string", Domain: domain, Key: key, Value: value}
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	if resp.Status != "success" {
		return fmt.Errorf("failed to set string: %s", resp.Message)
	}
	return nil
}

func (c *Client) GetString(domain, key string) (string, error) {
	req := Request{Action: "get_string", Domain: domain, Key: key}
	resp, err := c.sendRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Status != "success" {
		return "", fmt.Errorf("failed to get string: %s", resp.Message)
	}
	return resp.Value, nil
}

func (c *Client) InsertToSkipList(domain, slkey, key, value string) error {
	req := Request{Action: "insert_skiplist", Domain: domain, SLKey: slkey, Key: key, Value: value}
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	if resp.Status != "success" {
		return fmt.Errorf("failed to insert to skip list: %s", resp.Message)
	}
	return nil
}

func (c *Client) DeleteFromSkipList(domain, slkey, key string) error {
	req := Request{Action: "delete_skiplist", Domain: domain, SLKey: slkey, Key: key}
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	if resp.Status != "success" {
		return fmt.Errorf("failed to delete from skip list: %s", resp.Message)
	}
	return nil
}

func (c *Client) DeleteRangeFromSkipList(domain, slkey, minKey, maxKey string) error {
	req := Request{Action: "delete_range_skiplist", Domain: domain, SLKey: slkey, MinKey: minKey, MaxKey: maxKey}
	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}
	if resp.Status != "success" {
		return fmt.Errorf("failed to delete range from skip list: %s", resp.Message)
	}
	return nil
}

func (c *Client) SearchInSkipList(domain, slkey, key string) (string, error) {
	req := Request{Action: "search_skiplist", Domain: domain, SLKey: slkey, Key: key}
	resp, err := c.sendRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Status != "success" {
		return "", fmt.Errorf("failed to search in skip list: %s", resp.Message)
	}
	return resp.Value, nil
}

func (c *Client) RankInSkipList(domain, slkey, key string) (string, error) {
	req := Request{Action: "rank_skiplist", Domain: domain, SLKey: slkey, Key: key}
	resp, err := c.sendRequest(req)
	if err != nil {
		return "", err
	}
	if resp.Status != "success" {
		return "", fmt.Errorf("failed to get rank: %s", resp.Message)
	}
	return resp.Value, nil
}

func (c *Client) Increment(domain, key string) error {
	req := Request{Action: "increment", Domain: domain, Key: key}
	resp, err := c.sendRequest(req)
	if (err != nil) {
		return err
	}
	if (resp.Status != "success") {
		return fmt.Errorf("failed to increment: %s", resp.Message)
	}
	return nil
}

func (c *Client) Decrement(domain, key string) error {
	req := Request{Action: "decrement", Domain: domain, Key: key}
	resp, err := c.sendRequest(req)
	if (err != nil) {
		return err
	}
	if (resp.Status != "success") {
		return fmt.Errorf("failed to decrement: %s", resp.Message)
	}
	return nil
}