package sourcegraph

import (
	"fmt"
	"regexp"
)

type File struct {
	Repository string
	Path       string
	Contents   string
	TotalLines int
	Size       int
}

func (c client) GetFile(repository, path string, lines int) (File, error) {
	var q struct {
		Search struct {
			Results struct {
				Results []struct {
					FileMatch struct {
						File struct {
							Path       string
							Name       string
							Content    string `graphql:"content(startLine:0, endLine: $totalLines)"`
							TotalLines int
							ByteSize   int
						}
						Repository struct {
							Name string
						}
					} `graphql:"... on FileMatch"`
				}
			}
		} `graphql:"search(query: $q, patternType: regexp)"`
	} // query($q:  String!, $totalLines: Int!)

	err := c.gql.Query(c.ctx, &q, map[string]interface{}{
		"q":          fmt.Sprintf("repo:%s path:%s type:file", regexp.QuoteMeta(repository), regexp.QuoteMeta(path)),
		"totalLines": lines,
	})
	if err != nil {
		return File{}, err
	}

	if len(q.Search.Results.Results) != 1 {
		return File{}, fmt.Errorf("file not found")
	}

	return File{
		Repository: repository,
		Path:       path,
		Contents:   q.Search.Results.Results[0].FileMatch.File.Content,
		TotalLines: q.Search.Results.Results[0].FileMatch.File.TotalLines,
		Size:       q.Search.Results.Results[0].FileMatch.File.ByteSize,
	}, nil
}
