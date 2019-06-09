// Code generated by zenrpc; DO NOT EDIT.

package rpc

import (
	"context"
	"encoding/json"

	"github.com/semrush/zenrpc"
	"github.com/semrush/zenrpc/smd"
)

var RPC = struct {
	Server struct{ Run_Task, Is_Ready_Task, Get_Task_Result string }
}{
	Server: struct{ Run_Task, Is_Ready_Task, Get_Task_Result string }{
		Run_Task:        "run_task",
		Is_Ready_Task:   "is_ready_task",
		Get_Task_Result: "get_task_result",
	},
}

func (Server) SMD() smd.ServiceInfo {
	return smd.ServiceInfo{
		Description: ``,
		Methods: map[string]smd.Service{
			"Run_Task": {
				Description: ``,
				Parameters: []smd.JSONSchema{
					{
						Name:        "spec",
						Optional:    false,
						Description: ``,
						Type:        smd.Array,
						Items: map[string]string{
							"type": smd.Integer,
						},
					},
				},
				Returns: smd.JSONSchema{
					Description: ``,
					Optional:    false,
					Type:        smd.String,
				},
			},
			"Is_Ready_Task": {
				Description: ``,
				Parameters: []smd.JSONSchema{
					{
						Name:        "id",
						Optional:    false,
						Description: ``,
						Type:        smd.String,
					},
				},
				Returns: smd.JSONSchema{
					Description: ``,
					Optional:    false,
					Type:        smd.Boolean,
				},
			},
			"Get_Task_Result": {
				Description: ``,
				Parameters: []smd.JSONSchema{
					{
						Name:        "id",
						Optional:    false,
						Description: ``,
						Type:        smd.String,
					},
				},
				Returns: smd.JSONSchema{
					Description: ``,
					Optional:    false,
					Type:        smd.Object,
					Properties:  map[string]smd.Property{},
				},
			},
		},
	}
}

// Invoke is as generated code from zenrpc cmd
func (s Server) Invoke(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
	resp := zenrpc.Response{}
	var err error

	switch method {
	case RPC.Server.Run_Task:
		var args = struct {
			Spec []byte `json:"spec"`
		}{}

		if zenrpc.IsArray(params) {
			if params, err = zenrpc.ConvertToObject([]string{"spec"}, params); err != nil {
				return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
			}
		}

		if len(params) > 0 {
			if err := json.Unmarshal(params, &args); err != nil {
				return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
			}
		}

		resp.Set(s.Run_Task(args.Spec))

	case RPC.Server.Is_Ready_Task:
		var args = struct {
			Id string `json:"id"`
		}{}

		if zenrpc.IsArray(params) {
			if params, err = zenrpc.ConvertToObject([]string{"id"}, params); err != nil {
				return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
			}
		}

		if len(params) > 0 {
			if err := json.Unmarshal(params, &args); err != nil {
				return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
			}
		}

		resp.Set(s.Is_Ready_Task(args.Id))

	case RPC.Server.Get_Task_Result:
		var args = struct {
			Id string `json:"id"`
		}{}

		if zenrpc.IsArray(params) {
			if params, err = zenrpc.ConvertToObject([]string{"id"}, params); err != nil {
				return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
			}
		}

		if len(params) > 0 {
			if err := json.Unmarshal(params, &args); err != nil {
				return zenrpc.NewResponseError(nil, zenrpc.InvalidParams, "", err.Error())
			}
		}

		resp.Set(s.Get_Task_Result(args.Id))

	default:
		resp = zenrpc.NewResponseError(nil, zenrpc.MethodNotFound, "", nil)
	}

	return resp
}