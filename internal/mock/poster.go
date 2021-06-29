package mock

import (
	"io"
	"net/http"
	"net/url"

	"github.com/nestoroprysk/TelegramBots/internal/telegramclient"
)

// Poster is a mock poster.
type Poster struct {
	f func(url string, data url.Values) (resp *http.Response, err error)
}

var _ telegramclient.Poster = &Poster{}

// PosterOption defines the poster.
type PosterOption func(Poster) Poster

// NewPoster creates a poster.
func NewPoster(opts ...PosterOption) Poster {
	result := Poster{
		f: func(_ string, _ url.Values) (*http.Response, error) {
			rc := NewReadCloser(nil)
			return &http.Response{
				Body: &rc,
			}, nil
		},
	}

	for _, o := range opts {
		result = o(result)
	}

	return result
}

// PostForm posts and returns a response.
func (p Poster) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return p.f(url, data)
}

// PostFormError results in PostForm returning the input error.
func PostFormError(err error) PosterOption {
	return func(p Poster) Poster {
		return Poster{
			f: func(_ string, _ url.Values) (*http.Response, error) {
				return nil, err
			},
		}
	}
}

// PostFormStatusCode results in PostForm returning the response status code from input.
func PostFormStatusCode(statusCode int) PosterOption {
	return func(p Poster) Poster {
		result, _ := p.f("", nil)
		return Poster{
			f: func(_ string, _ url.Values) (*http.Response, error) {
				result.StatusCode = statusCode
				return result, nil
			},
		}
	}
}

// PostFormBody results in PostForm returning the response body from input.
func PostFormBody(body io.ReadCloser) PosterOption {
	return func(p Poster) Poster {
		result, _ := p.f("", nil)
		return Poster{
			f: func(_ string, _ url.Values) (*http.Response, error) {
				result.Body = body
				return result, nil
			},
		}
	}
}
