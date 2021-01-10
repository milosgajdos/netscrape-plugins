package star

import (
	"context"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	"github.com/milosgajdos/netscrape/pkg/metadata"
	"github.com/milosgajdos/netscrape/pkg/query"
	"github.com/milosgajdos/netscrape/pkg/query/base"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/object"
	"github.com/milosgajdos/netscrape/pkg/space/plan"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/space/top"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	// version is GitHub API version
	version = "v3"
	// topicRel is topic relation
	topicRel = "topic"
	// langRel is language relation
	langRel = "lang"
	// resKind is default space.Resource kind
	resKind = "starred"
	// globalNs is global namespace
	globalNs = "global"
	// workers is default number of workers
	workers = 5
	// pagins is default paging
	paging = 50
)

type resources struct {
	repo  space.Resource
	topic space.Resource
	lang  space.Resource
}

type scraper struct {
	// gh is GitHub API client.
	gh *github.Client
	// opts are scraper options
	opts Options
}

// NewScraper creates a new GitHub star repository scraper and returns it.
func NewScraper(gh *github.Client, opts ...Option) (*scraper, error) {
	copts := Options{}
	for _, apply := range opts {
		apply(&copts)
	}

	if copts.Workers <= 0 {
		copts.Workers = workers
	}

	if copts.Paging <= 0 {
		copts.Paging = paging
	}

	return &scraper{
		gh:   gh,
		opts: copts,
	}, nil
}

// Plan builds GH stars space plan and returns it.
func (g *scraper) Plan(ctx context.Context, o space.Origin) (space.Plan, error) {
	p, err := plan.New(o)
	if err != nil {
		return nil, err
	}

	repo, err := resource.New("repo", "repos", version, resKind, true)
	if err != nil {
		return nil, err
	}

	if err := p.Add(ctx, repo, space.AddOptions{}); err != nil {
		return nil, err
	}

	lang, err := resource.New("lang", "langs", version, resKind, false)
	if err != nil {
		return nil, err
	}

	if err := p.Add(ctx, lang, space.AddOptions{}); err != nil {
		return nil, err
	}

	topic, err := resource.New("topic", "topics", version, resKind, false)
	if err != nil {
		return nil, err
	}

	if err := p.Add(ctx, topic, space.AddOptions{}); err != nil {
		return nil, err
	}

	return p, nil
}

func (g *scraper) mapObjects(ctx context.Context, t *top.Top, to uuid.UID, res space.Resource, rel string, names []string) ([]space.Object, error) {
	objects := make([]space.Object, len(names))

	for i, name := range names {
		linkMd, err := metadata.New()
		if err != nil {
			return nil, err
		}
		linkMd.Set("relation", rel)

		// NOTE: we are setting uid to the name of the object
		// this is so we avoid duplicating topics with the same name
		uid, err := uuid.NewFromString(strings.ToLower(name + "-" + res.Name()))
		if err != nil {
			return nil, err
		}

		objects[i], err = object.New(uid, strings.ToLower(name), globalNs, res)
		if err != nil {
			return nil, err
		}

		if err := objects[i].Link(to, space.LinkOptions{Metadata: linkMd}); err != nil {
			return nil, err
		}

		if err := t.Add(ctx, objects[i], space.AddOptions{MergeLinks: true}); err != nil {
			return nil, err
		}
	}

	return objects, nil
}

func (g *scraper) fetchRepos(ctx context.Context, reposChan chan<- []*github.StarredRepository, done <-chan struct{}) error {
	defer close(reposChan)

	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: g.opts.Paging},
	}

	for {
		repos, resp, err := g.gh.Activity.ListStarred(ctx, g.opts.User, opts)
		if err != nil {
			return err
		}

		select {
		case reposChan <- repos:
		case <-ctx.Done():
			return nil
		case <-done:
			return nil
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil
}

func (g *scraper) mapRepos(ctx context.Context, reposChan <-chan []*github.StarredRepository, t *top.Top, res resources) error {
	for repos := range reposChan {
		// NOTE: we are only iterating over the repos resources
		// since topics and langs are merely adjacent nodes of repos objects
		// and do not have any API endpoint for querying them further
		for _, repo := range repos {
			m, err := metadata.New()
			if err != nil {
				return err
			}
			m.Set("starred_at", repo.StarredAt)
			m.Set("git_url", repo.Repository.GetURL)

			owner := *repo.Repository.Owner.Login
			if repo.Repository.Organization != nil {
				owner = *repo.Repository.Organization.Login
			}

			uid, err := uuid.NewFromString(*repo.Repository.NodeID)
			if err != nil {
				return err
			}

			topicObjects, err := g.mapObjects(ctx, t, uid, res.topic, topicRel, repo.Repository.Topics)
			if err != nil {
				return err
			}

			var langObjects []space.Object
			if repo.Repository.Language != nil {
				var err error
				langObjects, err = g.mapObjects(ctx, t, uid, res.lang, langRel, []string{*repo.Repository.Language})
				if err != nil {
					return err
				}
			}

			obj, err := object.New(uid, *repo.Repository.Name, owner, res.repo, object.Metadata(m))
			if err != nil {
				return err
			}

			for _, o := range topicObjects {
				md, err := metadata.New()
				if err != nil {
					return err
				}
				md.Set("relation", topicRel)

				if err := obj.Link(o.UID(), space.LinkOptions{Metadata: md}); err != nil {
					return err
				}
			}

			for _, o := range langObjects {
				md, err := metadata.New()
				if err != nil {
					return err
				}
				md.Set("relation", langRel)

				if err := obj.Link(o.UID(), space.LinkOptions{Metadata: md}); err != nil {
					return err
				}
			}

			if err := t.Add(ctx, obj, space.AddOptions{MergeLinks: true}); err != nil {
				return err
			}
		}
	}

	return nil
}

func getResource(ctx context.Context, p space.Plan, name, version string) (space.Resource, error) {
	q := base.Build().
		Add(query.Name(name), query.StringEqFunc(name)).
		Add(query.Version(version), query.StringEqFunc(version))

	rx, err := p.Get(ctx, q)
	if err != nil {
		return nil, err
	}

	return rx[0], nil
}

// Map builds a map of GH stars space topology and returns it.
// It returns error if any of the API calls fails with error.
func (g *scraper) Map(ctx context.Context, p space.Plan) (space.Top, error) {
	t, err := top.New(p)
	if err != nil {
		return nil, err
	}

	repoRes, err := getResource(ctx, p, "repo", version)
	if err != nil {
		return nil, err
	}

	topicRes, err := getResource(ctx, p, "topic", version)
	if err != nil {
		return nil, err
	}

	langRes, err := getResource(ctx, p, "lang", version)
	if err != nil {
		return nil, err
	}

	res := resources{
		repo:  repoRes,
		topic: topicRes,
		lang:  langRes,
	}

	reposChan := make(chan []*github.StarredRepository, g.opts.Workers)
	errChan := make(chan error)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// launch repo processing workers
	// these are building the graph
	for i := 0; i < g.opts.Workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case errChan <- g.mapRepos(ctx, reposChan, t, res):
			case <-done:
			case <-ctx.Done():
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- g.fetchRepos(ctx, reposChan, done)
	}()

	select {
	case err = <-errChan:
		close(done)
	case <-done:
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return t, nil
}
