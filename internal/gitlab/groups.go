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
	"errors"
	"github.com/xanzy/go-gitlab"
	"regexp"
)

func GetGroupProjects(app Application) ([]*gitlab.Project, error) {
	var groupsIDs []int
	groupsIDs = append(groupsIDs, *app.Group)

	if *app.Recursive {
		groupsIDs = append(groupsIDs, getSubgroups(app)...)
	}

	app.Log.Debug().Msgf("groups: %v", groupsIDs)

	var projects []*gitlab.Project

	for _, groupID := range groupsIDs {
		groupProjects, _, err := app.Client.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{
			Archived: gitlab.Ptr(false),
		})
		if err != nil {
			return nil, err
		}
		projects = append(projects, groupProjects...)
	}

	// Filter projects by matcher
	if *app.Matcher != "" {
		if err := validateMatcher(*app.Matcher); err != nil {
			return nil, errors.New("invalid regex matcher provided")
		}

		var filteredProjects []*gitlab.Project
		for _, project := range projects {
			if ok, _ := regexp.MatchString(*app.Matcher, project.Name); ok {
				filteredProjects = append(filteredProjects, project)
			}
		}

		projects = filteredProjects
	}

	return projects, nil
}

func getSubgroups(app Application) []int {
	var groups []int
	subGroups, _, err := app.Client.Groups.ListSubGroups(*app.Group, &gitlab.ListSubGroupsOptions{})
	if err != nil {
		app.Log.Error().Err(err).Msgf("failed to get subgroups for group %d", *app.Group)
	}

	for _, subGroup := range subGroups {
		groups = append(groups, subGroup.ID)
	}

	return groups
}

// validateMatcher validates the provided regex matcher
func validateMatcher(matcher string) error {
	_, err := regexp.Compile(matcher)
	if err != nil {
		return err
	}

	return nil
}
