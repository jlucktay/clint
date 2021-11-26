package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get working directory: %v", err)
		return
	}

	opts := &git.PlainOpenOptions{
		DetectDotGit: true,
	}

	repo, err := git.PlainOpenWithOptions(wd, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open git repo at '%s': %v", wd, err)
		return
	}

	tags, err := repo.Tags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read tags from git repo: %v", err)
		return
	}

	if err := tags.ForEach(func(ref *plumbing.Reference) error {
		fmt.Printf("hash: '%s'\nname: '%s'\n", ref.Hash(), ref.Name())

		obj, err := repo.TagObject(ref.Hash())
		switch err {
		case nil:
			// Tag object may be present
			fallthrough

		case plumbing.ErrObjectNotFound:
			// Not an annotated tag

			if obj == nil {
				fmt.Println()
				return nil
			}

		default:
			// Some other error
			return err
		}

		fmt.Printf("message: '%s'\ntagger: '%s' <'%s'> @ '%s'\n\n",
			strings.TrimSpace(obj.Message), obj.Tagger.Name, obj.Tagger.Email, obj.Tagger.When.Format(time.RFC3339))

		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "could not iterate over tags from git repo: %v", err)
		return
	}
}
