// Package main demonstrates a Legal Department Stack built on Lattice.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/latticehq/latticesdk/audit"
	"github.com/latticehq/latticesdk/stack"
	"github.com/latticehq/latticesdk/types"
)

func main() {
	ctx := context.Background()

	s, err := stack.New(stack.Config{
		RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
		APIKey:       os.Getenv("LATTICE_API_KEY"),
		StackName:    "legal-stack",
		StackVersion: "0.1.0",
	})
	if err != nil {
		log.Fatalf("failed to create stack: %v", err)
	}

	// Query recent audit logs for compliance review.
	logs, count, err := s.Audit.Query(ctx, &types.AuditLogFilter{
		SearchQuery: "resource_type:agent action:delete",
	}, nil)
	if err != nil {
		log.Fatalf("failed to query audit logs: %v", err)
	}
	fmt.Printf("Found %d audit events matching filter (total: %d)\n", len(logs), count)

	for _, l := range logs {
		fmt.Printf("  [%s] %s %s on %s (%s)\n",
			l.Time.Format("2006-01-02 15:04"),
			l.Action,
			l.ResourceType,
			l.ResourceTarget,
			l.Description,
		)
	}

	// Emit a custom audit event for contract review.
	err = s.Audit.Emit(ctx, audit.EmitEventRequest{
		Action:       types.AuditActionWrite,
		ResourceType: "contract_review",
	})
	if err != nil {
		log.Printf("failed to emit audit event: %v", err)
	}

	// Run the stack.
	if err := s.Run(ctx); err != nil {
		log.Fatalf("stack exited: %v", err)
	}
}
