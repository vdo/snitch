package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/karol-broda/snitch/internal/collector"
	"github.com/karol-broda/snitch/internal/geoip"
)

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 2 {
		return s[:max]
	}
	return s[:max-1] + SymbolEllipsis
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func sortFieldLabel(f collector.SortField) string {
	switch f {
	case collector.SortByLport:
		return "port"
	case collector.SortByProcess:
		return "proc"
	case collector.SortByPID:
		return "pid"
	case collector.SortByState:
		return "state"
	case collector.SortByProto:
		return "proto"
	default:
		return "port"
	}
}

func formatRemote(addr string, port int) string {
	if addr == "" || addr == "*" || port == 0 {
		return "-"
	}
	return fmt.Sprintf("%s:%d", addr, port)
}

// getConnectionGeoIP returns the IP to use for geolocation lookup
// Prefers remote address, falls back to local if remote is local/private
func getConnectionGeoIP(c collector.Connection) string {
	// First try remote address
	if c.Raddr != "" && c.Raddr != "*" && !geoip.IsLocalOrPrivate(c.Raddr) {
		return c.Raddr
	}

	// Fall back to local address if remote is local/private
	if c.Laddr != "" && c.Laddr != "*" && !geoip.IsLocalOrPrivate(c.Laddr) {
		return c.Laddr
	}

	return ""
}

// getConnectionFlag returns the country flag emoji for a connection
func getConnectionFlag(c collector.Connection) string {
	ip := getConnectionGeoIP(c)
	if ip == "" {
		return "  "
	}
	return geoip.GetFlag(ip)
}

// getConnectionOrg returns the organization for a connection
func getConnectionOrg(c collector.Connection) string {
	ip := getConnectionGeoIP(c)
	if ip == "" {
		return ""
	}
	return geoip.GetOrg(ip)
}

