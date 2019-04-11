package registry

import "time"

type tagsResponse struct {
	Tags []string `json:"tags"`
}

func (registry *Registry) Tags(repository string) (tags []string, err error) {
	url := registry.url("/v2/%s/tags/list", repository)

	var response tagsResponse
	for {
		registry.Logf("registry.tags url=%s repository=%s", url, repository)
		url, err = registry.getPaginatedJson(url, &response)
		switch err {
		case ErrNoMorePages:
			tags = append(tags, response.Tags...)
			return tags, nil
		case nil:
			tags = append(tags, response.Tags...)
			continue
		default:
			return nil, err
		}
	}
}

type acrTagsResponse struct {
	Registry  string   `json:"registry"`
	ImageName string   `json:"imageName"`
	Tags      []acrTag `json:"tags"`
}

type acrTag struct {
	Name           string    `json:"name"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
	Digest         string    `json:"digest"`
	CreatedTime    time.Time `json:"createdTime"`
	Signed         bool      `json:"signed"`
}

type changeableAttributes struct {
	DeleteEnabled bool `json:"deleteEnabled"`
	WriteEnabled  bool `json:"writeEnabled"`
	ReadEnabled   bool `json:"readEnabled"`
	ListEnabled   bool `json:"listEnabled"`
}

func (registry *Registry) AcrTags(repository string) (tags *acrTagsResponse, err error) {
	url := registry.url("/acr/v1/%s/_tags", repository)

	var response acrTagsResponse
	registry.Logf("registry.tags url=%s repository=%s", url, repository)
	err = registry.getJson(url, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
