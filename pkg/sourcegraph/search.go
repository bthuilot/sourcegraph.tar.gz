package sourcegraph

import "log"

type FileSearchResult struct {
	Repository string
	Path       string
	TotalLines int
}

func (c client) SearchFiles(query string) ([]FileSearchResult, error) {
	var q struct {
		Search struct {
			Results struct {
				Results []struct {
					FileMatch struct {
						File struct {
							Path        string
							Name        string
							IsDirectory bool
							Binary      bool
							TotalLines  int
						}
						Repository struct {
							Name string
						}
					} `graphql:"... on FileMatch"`
				}
			}
		} `graphql:"search(query: $query, patternType: regexp)"`
	}

	err := c.gql.Query(c.ctx, &q, map[string]interface{}{
		"query": query,
	})
	if err != nil {
		return nil, err
	}

	results := make([]FileSearchResult, len(q.Search.Results.Results))
	seen := make(map[string]struct{})
	for i, r := range q.Search.Results.Results {
		if r.FileMatch.File.IsDirectory || r.FileMatch.File.Binary {
			log.Printf("Skipping %s because it is a directory or binary file", r.FileMatch.File.Path)
			continue
		}
		id := r.FileMatch.Repository.Name + "/" + r.FileMatch.File.Path
		if _, ok := seen[id]; ok {
			log.Printf("Skipping %s because it is a duplicate", r.FileMatch.File.Path)
			continue
		}
		seen[id] = struct{}{}
		log.Printf("match for %s/%s", r.FileMatch.Repository.Name, r.FileMatch.File.Path)
		results[i] = FileSearchResult{
			Repository: r.FileMatch.Repository.Name,
			Path:       r.FileMatch.File.Path,
			TotalLines: r.FileMatch.File.TotalLines,
		}
	}
	return results, nil
}
