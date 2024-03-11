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

package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/schollz/progressbar/v3"
	"github.com/xanzy/go-gitlab"
	"io"
	"sync"
)

func Clone(app Application, projects []*gitlab.Project) error {
	if len(projects) == 0 {
		return fmt.Errorf("no projects found for group: %d", *app.Group)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(projects))
	concurrencyCh := make(chan struct{}, *app.Concurrency)
	doneCh := make(chan struct{})

	for _, project := range projects {
		wg.Add(1)

		go func(project *gitlab.Project) {
			defer wg.Done()

			concurrencyCh <- struct{}{}
			defer func() {
				<-concurrencyCh
			}()

			bar := progressbar.DefaultBytes(
				-1,
				"cloning "+project.Name,
			)

			fmt.Printf("Cloning project: %s\n", project.Name)
			path := fmt.Sprintf("%s/%s", project.Namespace.Name, project.Name)

			_, err := git.PlainClone(path, false, &git.CloneOptions{
				Auth: &http.BasicAuth{
					Username: *app.Auth.Username,
					Password: *app.Auth.Token,
				},
				URL:      project.HTTPURLToRepo,
				Progress: io.Writer(bar),
			})

			if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
				mu.Lock()
				errCh <- err
				mu.Unlock()
			}
		}(project)
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(doneCh) // Signal that all goroutines have finished
	}()

	<-doneCh

	var errs []error
	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered errors during cloning: %v", errs)
	}

	return nil
}
