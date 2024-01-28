package cache

import "strings"

// ensureTrailingDot ensures that the domain ends with a dot.
func ensureTrailingDot(domain string) string {
	if !strings.HasSuffix(domain, ".") {
		return domain + "."
	}
	return domain
}

// isWildcard checks if the domain is a wildcard domain.
func isWildcard(domain string) bool {
	return strings.Contains(domain, "*")
}

// convertToSQLPattern converts a wildcard domain to an SQL-like pattern.
func convertToSQLPattern(domain string) string {
	return strings.Replace(domain, "*", "%", 1)
}
