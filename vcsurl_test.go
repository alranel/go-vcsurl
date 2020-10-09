package vcsurl

import (
	"net/url"
	"reflect"
	"runtime"
	"testing"
)

func TestHttpRepo(t *testing.T) {
	url, _ := url.Parse("https://git.libreoffice.org/core")
	AssertEqual(t, IsGitHub(url), false)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)
}

func TestGitHub(t *testing.T) {
	url, _ := url.Parse("https://github.com/alranel")
	AssertEqual(t, IsGitHub(url), true)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsAccount(url), true)
	AssertNil(t, GetRepo(url))

	url, _ = url.Parse("https://github.com/alranel/go-vcsurl")
	AssertEqual(t, IsGitHub(url), true)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)
	AssertEqual(t, GetRepo(url).String(), url.String())
	AssertEqual(t, GetRawRoot(url).String(), "https://raw.githubusercontent.com/alranel/go-vcsurl/master/")

	url, _ = url.Parse("https://github.com/alranel/go-vcsurl.git")
	AssertEqual(t, GetRepo(url).String(), "https://github.com/alranel/go-vcsurl")
	AssertEqual(t, GetRawRoot(url).String(), "https://raw.githubusercontent.com/alranel/go-vcsurl/master/")

	url, _ = url.Parse("https://github.com/alranel/go-vcsurl/blob/master/README.md")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, GetRawFile(url).String(), "https://raw.githubusercontent.com/alranel/go-vcsurl/master/README.md")
	AssertEqual(t, GetRawRoot(url).String(), "https://raw.githubusercontent.com/alranel/go-vcsurl/master/")
	AssertEqual(t, IsRawRoot(GetRawRoot(url)), true)
	AssertEqual(t, GetRepo(url).String(), "https://github.com/alranel/go-vcsurl")
	AssertEqual(t, GetRepo(GetRawRoot(url)).String(), "https://github.com/alranel/go-vcsurl")

	url, _ = url.Parse("https://raw.githubusercontent.com/alranel/go-vcsurl/master") // no trailing slash
	AssertEqual(t, GetRepo(url).String(), "https://github.com/alranel/go-vcsurl")
}

func TestBitBucket(t *testing.T) {
	url, _ := url.Parse("https://bitbucket.org/Comune_Venezia/")
	AssertEqual(t, IsBitBucket(url), true)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsAccount(url), true)
	AssertNil(t, GetRepo(url))

	url, _ = url.Parse("https://bitbucket.org/Comune_Venezia/whistleblowing")
	AssertEqual(t, IsBitBucket(url), true)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)
	AssertEqual(t, GetRepo(url).String(), url.String())

	url, _ = url.Parse("https://bitbucket.org/Comune_Venezia/whistleblowing/src/master/LICENSE")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, GetRawFile(url).String(), "https://bitbucket.org/Comune_Venezia/whistleblowing/raw/master/LICENSE")
	AssertEqual(t, GetRawRoot(url).String(), "https://bitbucket.org/Comune_Venezia/whistleblowing/raw/master/")
	AssertEqual(t, IsRawRoot(GetRawRoot(url)), true)
	AssertEqual(t, GetRepo(url).String(), "https://bitbucket.org/Comune_Venezia/whistleblowing")
	AssertEqual(t, GetRepo(GetRawRoot(url)).String(), "https://bitbucket.org/Comune_Venezia/whistleblowing")
}

func TestGitLab(t *testing.T) {
	url, _ := url.Parse("https://gitlab.com/gitlab-org")
	AssertEqual(t, IsGitLab(url), true)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsAccount(url), true)
	AssertNil(t, GetRepo(url))

	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com")
	AssertEqual(t, IsGitLab(url), true)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, IsRepo(url), true)
	AssertEqual(t, IsAccount(url), false)
	AssertEqual(t, GetRepo(url).String(), url.String())

	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com/blob/master/LICENSE")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, GetRawFile(url).String(), "https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com/raw/master/LICENSE")
	AssertEqual(t, GetRawRoot(url).String(), "https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com/raw/master/")
	AssertEqual(t, IsRawRoot(GetRawRoot(url)), true)
	AssertEqual(t, GetRepo(url).String(), "https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com")
	AssertEqual(t, GetRepo(GetRawRoot(url)).String(), "https://gitlab.com/gitlab-org/gitlab-services/design.gitlab.com")

	// Self hosted GitLab
	url, _ = url.Parse("https://dev.funkwhale.audio/funkwhale/ansible")
	AssertEqual(t, IsGitLab(url), true)

	// Self hosted GitLab with HTTP URL and paths namespaced with '-'.
	url, _ = url.Parse("http://dev.funkwhale.audio/funkwhale/ansible/-/blob/master/README.md")
	AssertEqual(t, IsGitLab(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRawFile(url), false)
	AssertEqual(t, GetRawFile(url).String(), "http://dev.funkwhale.audio/funkwhale/ansible/raw/master/README.md")
	AssertEqual(t, GetRawRoot(url).String(), "http://dev.funkwhale.audio/funkwhale/ansible/raw/master/")
	AssertEqual(t, IsRawRoot(GetRawRoot(url)), true)
	AssertEqual(t, GetRepo(url).String(), "http://dev.funkwhale.audio/funkwhale/ansible")
	AssertEqual(t, GetRepo(GetRawRoot(url)).String(), "http://dev.funkwhale.audio/funkwhale/ansible")

	// New style raw paths
	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab/-/raw/master/README.md")
	AssertEqual(t, IsRawFile(url), true)
	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab/-/raw/master/")
	AssertEqual(t, IsRawRoot(url), true)
	url, _ = url.Parse("https://gitlab.com/gitlab-org/gitlab/")
	AssertEqual(t, IsRawRoot(url), false)

	// Old style raw paths
	url, _ = url.Parse("https://riuso.comune.salerno.it/root/simel_2/raw/master/publiccode.yml")
	AssertEqual(t, IsRawFile(url), true)
	url, _ = url.Parse("https://riuso.comune.salerno.it/root/simel_2/raw/master/")
	AssertEqual(t, IsRawRoot(url), true)
	url, _ = url.Parse("https://riuso.comune.salerno.it/")
	AssertEqual(t, IsRawRoot(url), false)

	url, _ = url.Parse("https://riuso.comune.salerno.it/root/simel_2/blob/master/publiccode.yml")
	AssertEqual(t, IsFile(url), true)
	AssertEqual(t, IsRepo(url), false)
	AssertEqual(t, IsRawFile(url), false)
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

// AssertNil checks if a value is nil
func AssertNil(t *testing.T, a interface{}) {
	if reflect.ValueOf(a).IsNil() {
		return
	}
	//debug.PrintStack()
	_, fn, line, _ := runtime.Caller(1)
	t.Errorf("%s:%d: Received %v (type %v), expected nil", fn, line, a, reflect.TypeOf(a))
}
