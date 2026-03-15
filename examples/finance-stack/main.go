// Package main demonstrates a Finance Department Stack built on Lattice.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/latticehq/latticesdk/coordination"
	"github.com/latticehq/latticesdk/stack"
)

func main() {
	ctx := context.Background()

	s, err := stack.New(stack.Config{
		RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
		APIKey:       os.Getenv("LATTICE_API_KEY"),
		StackName:    "finance-stack",
		StackVersion: "0.1.0",
	})
	if err != nil {
		log.Fatalf("failed to create stack: %v", err)
	}

	// Check budget before proceeding.
	me, err := s.Identity.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	orgs, err := s.Identity.ListOrganizations(ctx)
	if err != nil {
		log.Fatalf("failed to list organizations: %v", err)
	}

	for _, org := range orgs {
		hasQuota, err := s.Budget.CheckBudget(ctx, org.ID, me.ID)
		if err != nil {
			log.Printf("budget check failed for org %s: %v", org.Name, err)
			continue
		}
		fmt.Printf("Org %s: budget_available=%v\n", org.Name, hasQuota)
	}

	// Publish a cost event for other stacks.
	err = s.Coordination.Publish(ctx, coordination.PublishRequest{
		Topic: "finance.cost-alert",
		Payload: map[string]interface{}{
			"message": "Monthly budget at 80% utilization",
			"user":    me.Username,
		},
	})
	if err != nil {
		log.Printf("failed to publish coordination event: %v", err)
	}

	// Run the stack.
	if err := s.Run(ctx); err != nil {
		log.Fatalf("stack exited: %v", err)
	}
}
