package tools_test

import (
	"strings"
	"testing"

	"satset-app/internal/ai/tools"
)

func TestToolDefinitionsCountAndReadonly(t *testing.T) {
	defs := tools.GetToolDefinitions()

	if len(defs) != 13 {
		t.Fatalf("Expected exactly 13 tool definitions, got %d", len(defs))
	}

	expectedTools := map[string]bool{
		"get_business_summary":    false,
		"get_revenue_report":      false,
		"get_expense_report":      false,
		"get_orders_by_status":    false,
		"get_popular_services":    false,
		"get_unpaid_orders":       false,
		"get_overdue_orders":      false,
		"get_customer_info":       false,
		"get_recent_transactions": false,
		"get_cash_flow_summary":   false,
		"get_service_list":        false,
		"get_expense_categories":  false,
		"draft_customer_message":  false,
	}

	for _, d := range defs {
		if _, exists := expectedTools[d.Name]; !exists {
			t.Errorf("Unexpected tool name: %s", d.Name)
		}
		expectedTools[d.Name] = true

		if d.Description == "" {
			t.Errorf("Tool %s missing description", d.Name)
		}

		// Ensure tools do NOT have mutating keywords
		lower := strings.ToLower(d.Name)
		if strings.Contains(lower, "delete") || strings.Contains(lower, "update") || strings.Contains(lower, "create") || strings.Contains(lower, "cancel") || strings.Contains(lower, "refund") {
			t.Errorf("Tool %s violates read-only constraint", d.Name)
		}
	}

	for name, found := range expectedTools {
		if !found {
			t.Errorf("Missing expected tool: %s", name)
		}
	}
}
