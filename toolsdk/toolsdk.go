// Package toolsdk provides AI tool definitions and typed handler infrastructure
// for building agent tools that interact with the Lattice Runtime API.
package toolsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/latticehq/latticesdk/client"
)

// ToolDefinition describes a tool's schema for AI model consumption.
type ToolDefinition struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Parameters  map[string]Parameter `json:"parameters,omitempty"`
}

// Parameter describes a single tool parameter.
type Parameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// HandlerFunc is a typed handler for a tool.
type HandlerFunc[Arg any, Ret any] func(ctx context.Context, deps Deps, arg Arg) (Ret, error)

// Tool represents a typed AI tool with compile-time safe handler.
type Tool[Arg any, Ret any] struct {
	ToolDef            ToolDefinition
	Handler            HandlerFunc[Arg, Ret]
	UserClientOptional bool
}

// GenericHandlerFunc is a type-erased handler that accepts and returns raw JSON.
type GenericHandlerFunc func(ctx context.Context, deps Deps, arg json.RawMessage) (json.RawMessage, error)

// GenericTool is a type-erased tool for runtime use (e.g., tool registries).
type GenericTool struct {
	ToolDef            ToolDefinition
	Handler            GenericHandlerFunc
	UserClientOptional bool
}

// Deps provides dependencies for tool handlers.
type Deps struct {
	Client     *client.Client
	ReportTask func(ReportTaskArgs) error
}

// DepsOption configures optional Deps fields.
type DepsOption func(*Deps)

// WithTaskReporter sets the task reporter callback on Deps.
func WithTaskReporter(fn func(ReportTaskArgs) error) DepsOption {
	return func(d *Deps) {
		d.ReportTask = fn
	}
}

// NewDeps creates a new Deps with the given client and options.
func NewDeps(c *client.Client, opts ...DepsOption) Deps {
	d := Deps{Client: c}
	for _, opt := range opts {
		opt(&d)
	}
	return d
}

// Generic converts a typed tool to a generic tool by wrapping its handler
// with JSON marshaling/unmarshaling.
func (t Tool[Arg, Ret]) Generic() GenericTool {
	return GenericTool{
		ToolDef:            t.ToolDef,
		UserClientOptional: t.UserClientOptional,
		Handler: func(ctx context.Context, deps Deps, rawArg json.RawMessage) (json.RawMessage, error) {
			var arg Arg
			if len(rawArg) > 0 {
				if err := json.Unmarshal(rawArg, &arg); err != nil {
					return nil, fmt.Errorf("unmarshal tool argument: %w", err)
				}
			}
			ret, err := t.Handler(ctx, deps, arg)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(ret)
			if err != nil {
				return nil, fmt.Errorf("marshal tool result: %w", err)
			}
			return data, nil
		},
	}
}

// WithRecover wraps a GenericHandlerFunc with panic recovery.
// If the handler panics, the panic is caught and returned as an error.
func WithRecover(h GenericHandlerFunc) GenericHandlerFunc {
	return func(ctx context.Context, deps Deps, arg json.RawMessage) (result json.RawMessage, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("tool handler panicked: %v\n%s", r, debug.Stack())
			}
		}()
		return h(ctx, deps, arg)
	}
}

// AllTools returns all available tools as generic tools.
func AllTools() []GenericTool {
	return allTools
}

// ToolByName returns a generic tool by its name. If no tool matches,
// the second return value is false.
func ToolByName(name string) (GenericTool, bool) {
	for _, t := range allTools {
		if t.ToolDef.Name == name {
			return t, true
		}
	}
	return GenericTool{}, false
}
