package executable

import (
	"os"
	"path"
	"regexp"
	"strings"
)

var name string
var stripped string
var env string
var remote string
var config string
var extensions = []string{"app", "bat", "com", "exe", "sh"}

// Name returns the canonical name of the running executable. For example, if
// invocated as `/usr/local/bin/hephy apps`, this returns "hephy"
func Name() string {
	return name
}

// Env returns the name of the environment variable that is used to specify
// a profile for the CLI. This is guaranteed to always be uppercase and include
// no characters other than alphanumerics and underscores.
// Example: "HEPHY_PROFILE"
func Env() string {
	return env
}

// Remote returns the name of the default remote for a given invocation of the
// CLI. This is guaranteed to always be lowercase and include no characters
// other than alphanumerics and hyphens. Example: "hephy"
func Remote() string {
	return remote
}

// Config returns the name of the CLI's configuration directory. This is
// guaranteed to always be lowercase and include no characters other than
// alphanumerics, hyphens, and a single leading period. Example: ".hephy"
func Config() string {
	return config
}

// Render returns a copy of the given string that has had all template sentinels
// replaced with their corresponding values. The following sentinels are
// recognized:
//
// * {{.Name}} - the same as `Name()`
// * {{.Env}} - the same as `Config()`
// * {{.Remote}} - the same as `Remote()`
// * {{.Config}} - the same as `Config()`
func Render(input string) string {
	output := gsub(input, "{{.Name}}", name)
	output = gsub(output, "{{.Env}}", env)
	output = gsub(output, "{{.Remote}}", remote)
	output = gsub(output, "{{.Config}}", config)

	return output
}

func gsub(input string, sentinel string, value string) string {
	return strings.Replace(input, sentinel, value, -1)
}

func init() {
	name = path.Base(os.Args[0])
	stripped = strip(name)
	env = normalizedEnv(stripped)
	remote = normalizedRemote(stripped)
	config = normalizedConfig(stripped)
}

func normalizedEnv(input string) string {
	normalizer := regexp.MustCompile(`[^A-Z0-9_]`)

	return normalizer.ReplaceAllString(strings.ToUpper(input), "_") + "_PROFILE"
}

func normalizedRemote(input string) string {
	normalizer := regexp.MustCompile(`[^a-z0-9-]`)

	return normalizer.ReplaceAllString(strings.ToLower(input), "-")
}

func normalizedConfig(input string) string {
	normalizer := regexp.MustCompile(`[^a-z0-9-]`)

	return "." + normalizer.ReplaceAllString(strings.ToLower(input), "-")
}

func hasExt(input []string, ext string) bool {
	for _, p := range input {
		if p == ext {
			return true
		}
	}

	return false
}

func removeExt(input string, ext string) string {
	filtered := make([]string, 0)
	parts := strings.Split(input, ".")

	for _, part := range parts {
		if part == ext {
			continue
		}

		filtered = append(filtered, part)
	}

	return strings.Join(filtered, ".")
}

func strip(input string) string {
	output := strings.ToLower(input)

	for _, ext := range extensions {
		output = removeExt(output, ext)
	}

	return output
}
