// Package main demonstrates an HR Department Stack built on Lattice.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/latticehq/latticesdk/stack"
	"github.com/latticehq/latticesdk/types"
)

func main() {
	ctx := context.Background()

	s, err := stack.New(stack.Config{
		RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
		APIKey:       os.Getenv("LATTICE_API_KEY"),
		StackName:    "hr-stack",
		StackVersion: "0.1.0",
	})
	if err != nil {
		log.Fatalf("failed to create stack: %v", err)
	}

	// Verify we're authenticated.
	me, err := s.Identity.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}
	fmt.Printf("Authenticated as: %s (%s)\n", me.Username, me.Email)

	// List agents using the hr-agent template.
	agents, count, err := s.Identity.ListAgents(ctx, &types.AgentFilter{
		Template: "hr-agent",
	})
	if err != nil {
		log.Fatalf("failed to list agents: %v", err)
	}
	fmt.Printf("Found %d HR agents (total: %d)\n", len(agents), count)

	for _, a := range agents {
		// Check if we can update each agent.
		allowed, err := s.Authz.Can(ctx, types.RBACActionUpdate, types.RBACResourceAgent, a.ID.String())
		if err != nil {
			log.Printf("authz check failed for agent %s: %v", a.Name, err)
			continue
		}
		fmt.Printf("  - %s (healthy=%v, can_update=%v)\n", a.Name, a.Health.Healthy, allowed)
	}

	// Run the stack (heartbeats + graceful shutdown).
	if err := s.Run(ctx); err != nil {
		log.Fatalf("stack exited: %v", err)
	}
}
