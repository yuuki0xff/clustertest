package rpc

import (
	"context"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc"
	"github.com/yuuki0xff/clustertest/models"
	"os"
	"time"
)

type Client struct {
	client jsonrpc.RPCClient
}

func NewClient() (*Client, error) {
	addr := os.Getenv("CLUSTERTEST_SERVER")
	return &Client{
		client: jsonrpc.NewClient(addr),
	}, nil
}

func (c *Client) Create(task models.Task) (models.TaskID, error) {
	var id string
	err := c.call(&id, "run_task", task.SpecData())
	if err != nil {
		return nil, err
	}
	return &TaskID{id}, err
}
func (c *Client) Inspect(id models.TaskID) (models.TaskDetail, error) {
	var err error
	future := &Future{}
	future.id = id
	future.ready, err = c.isTaskReady(id)
	if err != nil {
		return nil, err
	}
	future.result, err = c.taskResult(id)
	if err != nil {
		return nil, err
	}
	return future, nil
}
func (c *Client) Wait(id models.TaskID, ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			ready, err := c.isTaskReady(id)
			if err != nil || !ready {
				continue
			}
			return nil
		case <-ctx.Done():
			return errors.Errorf("timeout")
		}
	}
}
func (c *Client) Cancel(id models.TaskID) error {
	panic("not impl")
}
func (c *Client) Delete(id models.TaskID) error {
	panic("not impl")
}
func (c *Client) call(out interface{}, method string, args ...interface{}) error {
	return c.client.CallFor(out, method, args)
}
func (c *Client) isTaskReady(id models.TaskID) (bool, error) {
	var ready bool
	err := c.call(&ready, "is_ready_task", id.String())
	return ready, err
}
func (c *Client) taskResult(id models.TaskID) (*Result, error) {
	result := &Result{}
	err := c.call(&result, "get_task_result", id.String())
	return result, err
}