package proto

import "strings"

// FullyQualifiedServiceName returns the fully-qualified service name.
func FullyQualifiedServiceName(pkg, svc string) string {
	var s []string
	if pkg != "" {
		s = append(s, pkg)
	}

	if svc != "" {
		s = append(s, svc)
	}

	return strings.Join(s, ".")
}

// FullyQualifiedMessageName returns the fully-qualified message name.
func FullyQualifiedMessageName(pkg, msg string) string {
	var s []string
	if pkg != "" {
		s = append(s, pkg)
	}

	if msg != "" {
		s = append(s, msg)
	}

	return strings.Join(s, ".")
}

// ParseFullyQualifiedServiceName returns the package and service name from a fully-qualified service name.
func ParseFullyQualifiedServiceName(fqsn string) (string, string) {
	i := strings.LastIndex(fqsn, ".")
	if i == -1 {
		return "", fqsn
	}
	return fqsn[:i], fqsn[i+1:]
}
