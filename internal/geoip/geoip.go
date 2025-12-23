package geoip

import (
	"net"
	"sync"
)

// IPInfo holds geolocation information for an IP
type IPInfo struct {
	CountryCode string
	Org         string
}

var (
	cache     = make(map[string]IPInfo)
	cacheMu   sync.RWMutex
	lookupSvc LookupService
	once      sync.Once
)

// LookupService defines the interface for IP geolocation
type LookupService interface {
	Lookup(ip string) IPInfo
}

// Initialize sets up the geolocation service
func Initialize() {
	once.Do(func() {
		lookupSvc = &ipAPIService{}
	})
}

// GetIPInfo returns geolocation info for an IP
// Returns empty IPInfo for local/private IPs or on error
func GetIPInfo(ip string) IPInfo {
	if ip == "" || ip == "*" {
		return IPInfo{}
	}

	// Check if it's a local/private IP
	if IsLocalOrPrivate(ip) {
		return IPInfo{}
	}

	// Check cache first
	cacheMu.RLock()
	if info, ok := cache[ip]; ok {
		cacheMu.RUnlock()
		return info
	}
	cacheMu.RUnlock()

	// Initialize if needed
	Initialize()

	// Lookup info
	info := lookupSvc.Lookup(ip)

	// Cache the result
	cacheMu.Lock()
	cache[ip] = info
	cacheMu.Unlock()

	return info
}

// GetCountryCode returns the ISO 3166-1 alpha-2 country code for an IP
func GetCountryCode(ip string) string {
	return GetIPInfo(ip).CountryCode
}

// GetOrg returns the organization for an IP
func GetOrg(ip string) string {
	return GetIPInfo(ip).Org
}

// IsLocalOrPrivate checks if an IP is local, loopback, or in private ranges
func IsLocalOrPrivate(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true // invalid IP, treat as local
	}

	// Loopback
	if ip.IsLoopback() {
		return true
	}

	// Link-local
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	// Unspecified (0.0.0.0 or ::)
	if ip.IsUnspecified() {
		return true
	}

	// Private ranges (RFC 1918 for IPv4, RFC 4193 for IPv6)
	if ip.IsPrivate() {
		return true
	}

	// Additional checks for special ranges
	// 100.64.0.0/10 (Carrier-grade NAT)
	cgnat := net.IPNet{
		IP:   net.ParseIP("100.64.0.0"),
		Mask: net.CIDRMask(10, 32),
	}
	if cgnat.Contains(ip) {
		return true
	}

	return false
}

// CountryFlag returns the flag emoji for a country code
func CountryFlag(countryCode string) string {
	if countryCode == "" || len(countryCode) != 2 {
		return "  " // two spaces for alignment
	}

	// Convert country code to regional indicator symbols
	// A = U+1F1E6, B = U+1F1E7, etc.
	code := []rune(countryCode)
	r1 := rune(0x1F1E6) + (rune(code[0]) - 'A')
	r2 := rune(0x1F1E6) + (rune(code[1]) - 'A')

	return string([]rune{r1, r2})
}

// GetFlag returns the flag emoji for an IP address
// Returns empty/spaces for local/private IPs
func GetFlag(ip string) string {
	return CountryFlag(GetIPInfo(ip).CountryCode)
}
