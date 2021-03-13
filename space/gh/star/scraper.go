package star

import (
	"context"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	"github.com/milosgajdos/netscrape/pkg/attrs"
	"github.com/milosgajdos/netscrape/pkg/space"
	"github.com/milosgajdos/netscrape/pkg/space/entity"
	"github.com/milosgajdos/netscrape/pkg/space/plan"
	"github.com/milosgajdos/netscrape/pkg/space/resource"
	"github.com/milosgajdos/netscrape/pkg/space/top"
	"github.com/milosgajdos/netscrape/pkg/uuid"
)

const (
	// workers is default number of workers
	workers = 5
	// paging is default GitHub API paging
	paging = 50
)

// scraper scrapes GitHub stars API
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

// params groups resource parameters
type params struct {
	name    string
	group   string
	version string
	kind    string
	ns      bool
}

func resourceParams() []params {
	return []params{
		{name: ownerRes, group: ownerGroup, version: version, kind: resKind, ns: true},
		{name: repoRes, group: repoGroup, version: version, kind: resKind, ns: true},
		{name: topicRes, group: topicGroup, version: version, kind: resKind, ns: true},
		{name: langRes, group: langGroup, version: version, kind: resKind, ns: false},
	}
}

// Plan creates a new space.Plan and adds GitHub stars resources to it.
func (s *scraper) Plan(ctx context.Context, o space.Origin) (space.Plan, error) {
	plan, err := plan.New(o)
	if err != nil {
		return nil, err
	}

	params := resourceParams()

	for _, p := range params {
		r, err := resource.New(p.name, p.group, p.version, p.kind, p.ns)
		if err != nil {
			return nil, err
		}

		if err := plan.Add(ctx, r); err != nil {
			return nil, err
		}
	}

	return plan, nil
}

// fetchRepos fetches GitHub repos into reposChan.
// Fetching can be stopped by closing done channel.
func (s *scraper) fetchRepos(ctx context.Context, reposChan chan<- []*github.StarredRepository, done <-chan struct{}) error {
	defer close(reposChan)

	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: s.opts.Paging},
	}

	for {
		repos, resp, err := s.gh.Activity.ListStarred(ctx, s.opts.User, opts)
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

// addEntities creates space entities from r with given names in namespace ns and adds them to top.
func (s *scraper) addEntities(ctx context.Context, top space.Top, names []string, ns string, r space.Resource) ([]space.Entity, error) {
	entities := make([]space.Entity, len(names))

	for i, name := range names {
		// NOTE: we are setting the uid to the name of the entity
		// this is so we avoid duplicating topics with the same name
		uid, err := uuid.NewFromString(strings.ToLower(name + "-" + r.Name()))
		if err != nil {
			return nil, err
		}

		entities[i], err = entity.New(strings.ToLower(name), ns, r, entity.WithUID(uid))
		if err != nil {
			return nil, err
		}

		if err := top.Add(ctx, entities[i]); err != nil {
			return nil, err
		}
	}

	return entities, nil
}

// mapRepos reads GH repos from respoChan and adds them to topology top.
// mapRepos also adds repo owner, topics and langs to top, too.
// Before the repos are added to top, several links are created:
// * owner of the repo is linked to the repo
// * repo topics and langs are added to top
// * repo is linked to all the topics and langs
func (s *scraper) mapRepos(ctx context.Context, reposChan <-chan []*github.StarredRepository, top space.Top, resMap map[string]space.Resource) error {
	for repos := range reposChan {
		// NOTE: we are only iterating over the repos resources
		// since owners, topics and langs are merely adjacent nodes of repos
		// and do not have any API endpoint for querying them further
		for _, repo := range repos {
			a, err := attrs.New()
			if err != nil {
				return err
			}
			a.Set("starred_at", repo.StarredAt.Format(dateTime))
			a.Set("git_url", repo.Repository.GetURL())

			uid, err := uuid.NewFromString(*repo.Repository.NodeID)
			if err != nil {
				return err
			}

			repoEnt, err := entity.New(*repo.Repository.Name, ns, resMap[repoRes], entity.WithUID(uid), entity.WithAttrs(a))
			if err != nil {
				return err
			}

			if err := top.Add(ctx, repoEnt); err != nil {
				return err
			}

			owner := *repo.Repository.Owner.Login
			if repo.Repository.Organization != nil {
				owner = *repo.Repository.Organization.Login
			}

			ownerUID, err := uuid.NewFromString(owner + "-" + ownerRes)
			if err != nil {
				return err
			}

			ownerEnt, err := entity.New(owner, ns, resMap[ownerRes], entity.WithUID(ownerUID))
			if err != nil {
				return err
			}

			if err := top.Add(ctx, ownerEnt); err != nil {
				return err
			}

			a, err = attrs.New()
			if err != nil {
				return err
			}
			a.Set(attrs.Relation, ownerRel)
			a.Set(attrs.DOTLabel, ownerRel)

			if err := top.Link(ctx, ownerEnt.UID(), repoEnt.UID(), space.WithAttrs(a), space.WithMerge()); err != nil {
				return err
			}

			topics, err := s.addEntities(ctx, top, repo.Repository.Topics, ns, resMap[topicRes])
			if err != nil {
				return err
			}

			for _, topic := range topics {
				a, err := attrs.New()
				if err != nil {
					return err
				}
				a.Set(attrs.Relation, topicRel)
				a.Set(attrs.DOTLabel, ownerRel)

				if err := top.Link(ctx, repoEnt.UID(), topic.UID(), space.WithAttrs(a), space.WithMerge()); err != nil {
					return err
				}
			}

			if repo.Repository.Language != nil {
				langs, err := s.addEntities(ctx, top, []string{*repo.Repository.Language}, ns, resMap[langRes])
				if err != nil {
					return err
				}

				for _, lang := range langs {
					a, err := attrs.New()
					if err != nil {
						return err
					}
					a.Set(attrs.Relation, langRel)
					a.Set(attrs.DOTLabel, ownerRel)

					if err := top.Link(ctx, repoEnt.UID(), lang.UID(), space.WithAttrs(a), space.WithMerge()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// Map builds a map of GH stars space topology and returns it.
// It returns error if any of the API calls fails with error.
func (s *scraper) Map(ctx context.Context, p space.Plan) (space.Top, error) {
	top, err := top.New()
	if err != nil {
		return nil, err
	}

	rx, err := p.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	resMap := make(map[string]space.Resource)

	for _, r := range rx {
		resMap[r.Name()] = r
	}

	reposChan := make(chan []*github.StarredRepository, s.opts.Workers)
	errChan := make(chan error)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// launch repo processing workers
	// these are building the graph
	for i := 0; i < s.opts.Workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case errChan <- s.mapRepos(ctx, reposChan, top, resMap):
			case <-done:
			case <-ctx.Done():
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case errChan <- s.fetchRepos(ctx, reposChan, done):
		case <-ctx.Done():
		}
	}()

	select {
	case err = <-errChan:
	case <-ctx.Done():
	}

	close(done)
	wg.Wait()

	if err != nil {
		return nil, err
	}

	return top, nil
}
