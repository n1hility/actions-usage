package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v31/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

func trunc(str string, size int) string {
	if len(str) > size {
		return str[0:size]
	}
	return str
}

func printChar(char string) {
    fmt.Printf(char)
}

func formatDuration(duration time.Duration) string {
    days := duration / (24 * time.Hour)
    hours := (duration - (days * 24 * time.Hour)) / time.Hour
    mins := (duration - ((days * 24 + hours) * time.Hour)) / time.Minute
    
    if days > 0 {
        return fmt.Sprintf("%02d d %02d h", days, hours)
    } 

    return fmt.Sprintf("%02d h %02d m", hours, mins)
}

func readToken(name string) (string, error) {
    dir, err := homedir.Dir()
    if err != nil {
        return "", err
    }

    name = filepath.Join(dir, name)
    data, err := ioutil.ReadFile(name)
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(data)), err
}

func createClient(token string) *github.Client {
    ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return github.NewClient(tc)
}

func getAllRepositories(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error){
    opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
    }
    var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err !=  nil {
            return nil, err
        } 
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
    }
    
    return allRepos, nil
}

func getAllWorkflowRuns(ctx context.Context, client *github.Client, repos []*github.Repository, print func (string))  ([]*github.WorkflowRun, error) {
    var allRuns []*github.WorkflowRun
    var err error
	for _, repo := range repos {
        var retries,retries2 int
		print(".")
		allRuns, retries, err = getWorkflowRunsByStatus(context.Background(), client, repo.GetFullName(), "in_progress", allRuns)
		if err == nil {
            allRuns, retries2, err = getWorkflowRunsByStatus(context.Background(), client, repo.GetFullName(), "queued", allRuns)
            retries += retries2
        }
        if retries > 0 {
            print(strings.Repeat("e", retries))
        }
		if err != nil {
			return nil, err
		}
    }

    return allRuns, nil
}

func computeJobs(ctx context.Context, client *github.Client, org string, allRuns []*github.WorkflowRun, print func (string)) ([]string, int, int, int, error) {
    workflows := make(map[string]*github.Workflow)
    var tq, tr, tc int
    var results []string
    var err error
	for i, run := range allRuns {
		key := run.GetWorkflowURL()
		workflow := workflows[key]
		if workflow == nil {
			workflow, _, err = getWorkflow(ctx, client, key)
			if err != nil {
				return nil, 0, 0, 0, err
			}
			workflows[key] = workflow
        }
        queued, inProgress, completed, err := countWorkflowRunJobs(ctx, client, org, run.GetRepository().GetName(), run.GetID())
        if err != nil {
            return nil, 0, 0, 0, err
        }
        tq += queued; tr += inProgress; tc += completed;
        count := fmt.Sprintf("%3vq / %3vr / %3vc", queued, inProgress, completed)
        age := time.Now().Sub(run.GetCreatedAt().Time).Round(time.Minute)

        source := fmt.Sprintf("%v:%v", run.GetHeadRepository().GetFullName(), run.GetHeadBranch())
        results = append(results, fmt.Sprintf("%3v.  %10v  %18v  %20v  %13v  %-15v  %9v  %v\n", i+1, run.GetID(), count, trunc(workflow.GetName(), 20), 
                         trunc(run.GetEvent(),13), trunc(run.GetRepository().GetName(), 15), formatDuration(age), trunc(source, 40)))
        print(".")
    }

    return results, tq, tr, tc, nil
}

func getWorkflowRunsByStatus(ctx context.Context, client *github.Client, fullRepo string, status string, runs []*github.WorkflowRun) ([]*github.WorkflowRun, int, error) {
	opt := &github.ListWorkflowRunsOptions{
		Status:      status,
		ListOptions: github.ListOptions{PerPage: 100},
	}

    retries := 5
	for {
		listRuns, resp, err := listRepositoryWorkflowRuns(ctx, client, fullRepo, opt)
		if err != nil {
            if (retries > 0 && resp.StatusCode >= 500 && resp.StatusCode < 600) {
                retries--
                continue
            }
			return runs, 5 - retries, err
		}
		runs = append(runs, listRuns.WorkflowRuns...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return runs, 0, nil
}

func countWorkflowRunJobs(ctx context.Context, client *github.Client, owner string, repo string, runID int64) (int, int, int, error) {
	opt := &github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

    var jobs []*github.WorkflowJob
	for {
		listJobs, resp, err := client.Actions.ListWorkflowJobs(ctx, owner, repo, runID, opt)
		if err != nil {
			return 0, 0, 0, err
		}
		jobs = append(jobs, listJobs.Jobs...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
    }

    counts := make(map[string]int)
    for _, job := range jobs {
        counts[job.GetStatus()]++
    }

	return counts["queued"], counts["in_progress"], counts["completed"], nil
}


func main() {
    if len(os.Args) < 2 {
        fmt.Printf("Usage: %v [org]\n", filepath.Base(os.Args[0]))
        return
    }

    org := os.Args[1]
    
    token, err := readToken("actions-usage.tok")
    if err != nil {
        fmt.Printf("Token not found!\n")
        fmt.Printf("Before you can use this tool, you must create a read-only Github personal \n")
        fmt.Printf("access token and place it in your home directory under the name:\n\n")
        fmt.Printf("'actions-usage.tok'\n\n")
        fmt.Printf("https://github.com/settings/tokens\n\n")
        return
    }    

    client := createClient(token)

    allRepos, err := getAllRepositories(context.Background(), client, org)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Finding workflows running on all repositories on %v", org)    
	allRuns, err := getAllWorkflowRuns(context.Background(), client, allRepos, printChar)
	if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Println()

    fmt.Printf("Analyzing jobs")
    results, tq, tr, tc, err := computeJobs(context.Background(), client, org, allRuns, printChar)    
    fmt.Printf("\n\n%4v  %10v  %18v  %20v  %13v  %-15v  %-9v  %v\n", "", "wf id", "queue/  run / comp", "name", "event", "repo", "created", "source")
    fmt.Printf("%4v  %v\n", "", strings.Repeat("-", 137))
    for _, line := range results {
        fmt.Printf(line);
    }
    fmt.Printf("%4v  %10v  %3vq / %3vr / %3vc\n", "", "Total:", tq, tr, tc);
}