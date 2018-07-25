package core

const (
	NamePolicy = "Policy"
	NameScheme = "Scheme"
	NameRemote = "Remote"
	NameHost   = "Host"
	NameURL    = "URL"
	NameProxy  = "Proxy"
)

var RecordFormat = []string{
	NameScheme,
	NameHost,
	NameRemote,
	NamePolicy,
	NameProxy,
	NameURL,
}
