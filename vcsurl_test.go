package vcsurl

import (
	"net/url"
	"reflect"
	"runtime"
	"testing"
)

func TestGitHub(t *testing.T) {
	url, _ := url.Parse("https://github.com/alranel")
	AssertEqual(t, IsGitHub(url), true)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsAccount(url), true)

	url, _ = url.Parse("https://github.com/alranel/go-vcsurl")
	AssertEqual(t, IsGitHub(url), true)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)

	url, _ = url.Parse("https://github.com/alranel/go-vcsurl/blob/master/README.md")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, GetRawFile(url).String(), "https://raw.githubusercontent.com/alranel/go-vcsurl/master/README.md")
	AssertEqual(t, GetRawRoot(url).String(), "https://raw.githubusercontent.com/alranel/go-vcsurl/master/")
}

func TestBitBucket(t *testing.T) {
	url, _ := url.Parse("https://bitbucket.org/Comune_Venezia/")
	AssertEqual(t, IsBitBucket(url), true)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsAccount(url), true)

	url, _ = url.Parse("https://bitbucket.org/Comune_Venezia/whistleblowing")
	AssertEqual(t, IsBitBucket(url), true)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)

	url, _ = url.Parse("https://bitbucket.org/Comune_Venezia/whistleblowing/src/master/LICENSE")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, GetRawFile(url).String(), "https://bitbucket.org/Comune_Venezia/whistleblowing/raw/master/LICENSE")
	AssertEqual(t, GetRawRoot(url).String(), "https://bitbucket.org/Comune_Venezia/whistleblowing/raw/master/")
}

func TestGitLab(t *testing.T) {
	url, _ := url.Parse("https://gitlab.com/gitlab-org")
	AssertEqual(t, IsGitLab(url), true)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsAccount(url), true)

	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com")
	AssertEqual(t, IsGitLab(url), true)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)

	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com/blob/master/LICENSE")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRaw(url), false)
	AssertEqual(t, GetRawFile(url).String(), "https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com/raw/master/LICENSE")
	AssertEqual(t, GetRawRoot(url).String(), "https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com/raw/master/")
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	//debug.PrintStack()
	_, fn, line, _ := runtime.Caller(1)
	t.Errorf("%s:%d: Received %v (type %v), expected %v (type %v)", fn, line, a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
