/*
Copyright © 2024 Furkan Pehlivan <furkanpehlivan34@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package gitlab

import (
	"crypto/tls"
	"github.com/pehlicd/gsc/internal"
	"github.com/xanzy/go-gitlab"
	"net/http"
)

func NewClient(auth *internal.Auth) (*gitlab.Client, error) {
	if *auth.Insecure {
		return newInsecureClient(auth)
	}

	client, err := gitlab.NewClient(*auth.Token, gitlab.WithBaseURL(*auth.Host))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newInsecureClient(auth *internal.Auth) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(*auth.Token, gitlab.WithBaseURL(*auth.Host), gitlab.WithHTTPClient(insecureHTTPClient()))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func insecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
