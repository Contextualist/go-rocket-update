package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type providerGithub struct {
	repoURL     string
	zipName     string
	zipProvider *providerZip
}

type GithubCommit struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

type GithubTag struct {
	Name       string       `json:"name"`
	ZipballURL string       `json:"zipball_url"`
	TarballURL string       `json:"tarball_url"`
	Commit     GithubCommit `json:"commit"`
	NodeID     string       `json:"node_id"`
}

type GithubTags struct {
	Tag []GithubTag
}

type repositoryInfo struct {
	RepositoryOwner string
	RepositoryName  string
}

// NewProviderGithub creates a new provider for local files
func NewProviderGithub(repoURL, zipName string) Provider {
	return &providerGithub{
		repoURL: repoURL,
		zipName: zipName,
	}
}

// getRepositoryInfo parses the github repository URL
func (c *providerGithub) repositoryInfo() (*repositoryInfo, error) {
	re := regexp.MustCompile(`github\.com/(.*?)/(.*?)$`)
	submatches := re.FindAllStringSubmatch(c.repoURL, 1)
	if len(submatches) < 1 {
		return nil, errors.New("Invalid github URL:" + c.repoURL)
	}
	return &repositoryInfo{
		RepositoryOwner: submatches[0][1],
		RepositoryName:  submatches[0][2],
	}, nil
}

// tagsURL get the tags URL for the github repository
func (c *providerGithub) tagsURL() (string, error) {
	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags",
		info.RepositoryOwner,
		info.RepositoryName,
	), nil
}

// zipURL get the zip URL
func (c *providerGithub) zipURL(tag string) (string, error) {
	info, err := c.repositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		info.RepositoryOwner,
		info.RepositoryName,
		tag,
		c.zipName,
	), nil
}

// getTags gets tags of the repository
func (c *providerGithub) getTags() (string, error) {
	tagsURL, err := c.tagsURL()
	if err != nil {
		return "", err
	}
	fmt.Println(tagsURL)
	response, err := http.Get(tagsURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var tags []GithubTag
	err = json.NewDecoder(response.Body).Decode(&tags)
	if err != nil {
		return "", err
	}
	fmt.Println(tags)
	return "", nil
}

// Open opens the provider
func (c *providerGithub) Open() error {
	_, err := c.getTags()
	return err
}

// Close closes the provider
func (c *providerGithub) Close() error {
	return nil
}

// GetLatestVersion gets the lastest version
func (c *providerGithub) GetLatestVersion() (string, error) {
	return "1.0", nil
}

// Walk walks all the files provided
func (c *providerGithub) Walk(walkFn WalkFunc) error {
	return nil
}

// Retrieve file relative to "provider" to destination
func (c *providerGithub) Retrieve(src string, dest string) error {
	return nil
}
