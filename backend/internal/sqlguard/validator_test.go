package sqlguard

import "testing"

func TestValidatorAllowsSelectAndAppliesLimit(t *testing.T) {
	result := NewValidator().Validate("postgres", "SELECT * FROM users")
	if !result.Valid || !result.Executable {
		t.Fatalf("expected valid select, got %+v", result)
	}
	if !result.LimitApplied || result.Limit == nil || *result.Limit != 1000 {
		t.Fatalf("expected default limit, got %+v", result)
	}
}

func TestValidatorRejectsDangerousStatements(t *testing.T) {
	validator := NewValidator()
	for _, sql := range []string{"INSERT INTO users VALUES (1)", "UPDATE users SET id=1", "DELETE FROM users", "DROP TABLE users", "SELECT 1; SELECT 2"} {
		result := validator.Validate("postgres", sql)
		if result.Valid || result.Executable {
			t.Fatalf("expected invalid for %q, got %+v", sql, result)
		}
	}
}

func TestMySQLValidatorParsesSelectAndAppliesLimit(t *testing.T) {
	result := NewValidator().Validate("mysql", "SELECT id, name FROM users")
	if !result.Valid || !result.Executable {
		t.Fatalf("expected valid mysql select, got %+v", result)
	}
	if !result.LimitApplied || result.Limit == nil || *result.Limit != 1000 {
		t.Fatalf("expected mysql default limit, got %+v", result)
	}
}

func TestMySQLValidatorRejectsInvalidSQL(t *testing.T) {
	result := NewValidator().Validate("mysql", "SELECT FROM")
	if result.Valid || result.Executable {
		t.Fatalf("expected invalid mysql syntax, got %+v", result)
	}
}
