package vcsurl

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

var gitlabDomains = make(map[string]bool)

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
	if url.Host == "gitlab.com" {
		return true
	}

	// Did we already validate this URL as a GitLab site?
	if _, seen := gitlabDomains[url.Host]; seen {
		return true
	}

	// Detect GitLab running on custom domains by performing a HTTP request and looking
	// for the _gitlab_session cookie.
	url2, _ := url.Parse(url.String())
	url2.Path = "/api"
	resp, err := http.Get(url2.String())
	if err == nil {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "_gitlab_session" {
				gitlabDomains[url.Host] = true
				return true
			}
		}
	}

	return false
}

// IsHttpRepo returns true if the supplied URL points to an HTTP git repository.
func IsHttpRepo(url *url.URL) bool {
	refsUrl, _ := url.Parse(url.String())
	refsUrl.Path = path.Join(refsUrl.Path, "/info/refs")

	query := url.Query()
	query.Set("service", "git-upload-pack")
	refsUrl.RawQuery = query.Encode()

	resp, err := http.Head(refsUrl.String())
	if err == nil {
		return resp.StatusCode >= 200 && resp.StatusCode <= 299
	}

	return false
}

// IsAccount returns true if the supplied URL points to the root page of an org or user account.
func IsAccount(url *url.URL) bool {
	if url.Host == "github.com" {
		if ok, _ := regexp.MatchString("^/[^/]+$", url.Path); ok {
			return true
		}
	} else if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/?$", url.Path); ok {
			return true
		}
	} else if IsGitLab(url) {
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
	} else if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+$", url.Path); ok {
			return true
		}
	} else if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+){2,}/?$", url.Path); ok {
			if ok, _ := regexp.MatchString("/(blob|raw)/", url.Path); !ok {
				return true
			}
		}
	} else if IsHttpRepo(url) {
		return true
	}
	return false
}

// IsFile returns true if the supplied URL points to a file in non-raw mode.
func IsFile(url *url.URL) bool {
	if url.Host == "github.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/blob/[^/]+/.+$", url.Path); ok {
			return true
		}
	} else if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/src/[^/]+/.+$", url.Path); ok {
			return true
		}
	} else if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+)+/(-/)?blob/[^/]+/.+$", url.Path); ok {
			return true
		}
	}
	return false
}

// IsRawFile returns true if the supplied URL points to a raw file.
func IsRawFile(url *url.URL) bool {
	if url.Host == "github.com" {
		return false
	} else if url.Host == "raw.githubusercontent.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/[^/]+/.+$", url.Path); ok {
			return true
		}
	} else if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/raw/[^/]+/.+$", url.Path); ok {
			return true
		}
	} else if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+)+/(-/)?raw/[^/]+/.+$", url.Path); ok {
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
	} else if url.Host == "bitbucket.org" {
		re := regexp.MustCompile("^https://bitbucket.org/([^/]+)/([^/]+)/src/(.+)$")
		url, _ := url.Parse(re.ReplaceAllString(url.String(), "https://bitbucket.org/$1/$2/raw/$3"))
		return url
	} else if IsGitLab(url) {
		re := regexp.MustCompile("^(https?://.+?(?:/[^/]+)+?)/(?:-/)?blob/([^/]+/.+)$")
		url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/raw/$2"))
		return url
	}
	return nil
}

// IsRawRoot returns true if the supplied URL is the root for raw files.
func IsRawRoot(url *url.URL) bool {
	if url.Host == "github.com" {
		return false
	} else if url.Host == "raw.githubusercontent.com" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/[^/]+/?$", url.Path); ok {
			return true
		}
	} else if url.Host == "bitbucket.org" {
		if ok, _ := regexp.MatchString("^/[^/]+/[^/]+/raw/[^/]+/?$", url.Path); ok {
			return true
		}
	} else if IsGitLab(url) {
		if ok, _ := regexp.MatchString("^(/[^/]+)+/(-/)?raw/[^/]+/?$", url.Path); ok {
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
		} else if url.Host == "bitbucket.org" {
			re := regexp.MustCompile("^(https://bitbucket.org/[^/]+/[^/]+/raw/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/"))
			return url
		} else if IsGitLab(url) {
			re := regexp.MustCompile("^(https?://.+?(?:/[^/]+)+/(-/)?raw/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1/"))
			return url
		}
	} else {
		if url.Host == "github.com" {
			//strip .git if present
			url.Path = strings.TrimSuffix(url.Path, ".git")

			re := regexp.MustCompile("^https://github.com/([^/]+)/([^/]+)$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "https://raw.githubusercontent.com/$1/$2/master/"))
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
		} else if url.Host == "bitbucket.org" {
			re := regexp.MustCompile("^(https://bitbucket.org/[^/]+/[^/]+).*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1"))
			return url
		} else if IsGitLab(url) {
			re := regexp.MustCompile("^(https?://.+?(?:/[^/]+)+?)/(-/)?raw/.*$")
			url, _ := url.Parse(re.ReplaceAllString(url.String(), "$1"))
			return url
		}
	}
	return nil
}
