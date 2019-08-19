package vcsurl

import (
	"net/url"
	"regexp"
	"strings"
)

// IsGitHub returns true if the supplied URL belongs to GitHub.
func IsGitHub(url *url.URL) bool {
	return url.Host == "github.com" || url.Host == "raw.githubusercontent.com"
}

// IsBitBucket returns true if the supplied URL belongs to BitBucket.
func IsBitBucket(url *url.URL) bool {
	return url.Host == "bitbucket.org"
}

// IsGitLab returns true if the supplied URL belongs to BitBucket.
func IsGitLab(url *url.URL) bool {
	return url.Host == "gitlab.com"
}

// IsAccount returns true if the supplied URL points to the root page of an org or user account.
func IsAccount(url *url.URL) bool {
	if url.Host == "github.com" {
		if ok, _ := regexp.MatchString("^/[^/]+$", url.Path); ok {
			return true
		}
	}
	if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/?$", url.Path); ok {
			return true
		}
	}
	if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^/[^/]+/?$", url.Path); ok {
			return true
		}
	}
	return false
}

// IsRepo returns true if the supplied URL points to the root page of a repository.
func IsRepo(url *url.URL) bool {
	if url.Host == "github.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+$", url.Path); ok {
			return true
		}
	}
	if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+$", url.Path); ok {
			return true
		}
	}
	if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+){2,}/?$", url.Path); ok {
			if ok, _ := regexp.MatchString("/(blob|raw)/", url.Path); !ok {
				return true
			}
		}
	}
	return false
}

// IsFile returns true if the supplied URL points to a file in non-raw mode.
func IsFile(url *url.URL) bool {
	if url.Host == "github.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/blob/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/src/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+)+/blob/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	return false
}

// IsRawFile returns true if the supplied URL points to a raw file.
func IsRawFile(url *url.URL) bool {
	if url.Host == "raw.githubusercontent.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/raw/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+)+/raw/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	return false
}

// GetRawFile returns the raw URL corresponding to the supplied file URL.
// In case the supplied file URL is already raw, it will be returned back.
func GetRawFile(url *url.URL) *url.URL {
	if IsRawFile(url) {
		return url
	}
	if url.Host == "github.com" {
		re := regexp.MustCompile("^https://github.com/([^/]+)/([^/]+)/blob/(.+)$")
		url, _ := url.Parse(re.ReplaceAllString(url.String(), "https://raw.githubusercontent.com/$1/$2/$3"))
		return url
	}
	if url.Host == "bitbucket.org" {
		re := regexp.MustCompile("^https://bitbucket.org/([^/]+)/([^/]+)/src/(.+)$")
		url, _ := url.Parse(re.ReplaceAllString(url.String(), "https://bitbucket.org/$1/$2/raw/$3"))
		return url
	}
	if IsGitLab(url) {
		re := regexp.MustCompile("^(https://.+?(?:/[^/]+)+)/blob/([^/]+/.+)$")
		url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/raw/$2"))
		return url
	}
	return nil
}

// IsRawRoot returns true if the supplied URL is the root for raw files.
func IsRawRoot(url *url.URL) bool {
	if url.Host == "raw.githubusercontent.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/[^/]+/$", url.Path); ok {
			return true
		}
	}
	if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/raw/[^/]+/$", url.Path); ok {
			return true
		}
	}
	if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+)+/raw/[^/]+/$", url.Path); ok {
			return true
		}
	}
	return false
}

// GetRawRoot returns the URL of the raw repository root containing the supplied file.
func GetRawRoot(url *url.URL) *url.URL {
	if IsFile(url) || IsRawFile(url) {
		url = GetRawFile(url)
		if url.Host == "raw.githubusercontent.com" {
			re := regexp.MustCompile("^(https://raw.githubusercontent.com/[^/]+/[^/]+/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/"))
			return url
		}
		if url.Host == "bitbucket.org" {
			re := regexp.MustCompile("^(https://bitbucket.org/[^/]+/[^/]+/raw/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/"))
			return url
		}
		if IsGitLab(url) {
			re := regexp.MustCompile("^(https://.+?(?:/[^/]+)+/raw/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/"))
			return url
		}
	}
	return nil
}

// GetRepo returns the URL of the main page of the repository (i.e. not raw nor git)
func GetRepo(url *url.URL) *url.URL {
	if IsRepo(url) {
		url.Path = strings.TrimSuffix(url.Path, ".git")
		return url
	}
	if IsFile(url) || IsRawFile(url) || IsRawRoot(url) {
		if IsFile(url) {
			url = GetRawFile(url)
		}
		if url.Host == "raw.githubusercontent.com" {
			re := regexp.MustCompile("^https://raw.githubusercontent.com/([^/]+/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "https://github.com/$1"))
			return url
		}
		if url.Host == "bitbucket.org" {
			re := regexp.MustCompile("^(https://bitbucket.org/[^/]+/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1"))
			return url
		}
		if IsGitLab(url) {
			re := regexp.MustCompile("^(https://.+?(?:/[^/]+)+?)/raw/.*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1"))
			return url
		}
	}
	return nil
}
