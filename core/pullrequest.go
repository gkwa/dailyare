package core

type PullRequest struct {
	Merged bool   `json:"merged"`
	Title  string `json:"title"`
}

type PRService interface {
	GetPRStatus(url string) (bool, error)
}

type githubPRService struct {
	client GithubClient
}

func NewGithubPRService(client GithubClient) PRService {
	return &githubPRService{client: client}
}

func (s *githubPRService) GetPRStatus(url string) (bool, error) {
	apiURL := formatGithubURL(url)
	var pr PullRequest
	err := s.client.Get(apiURL, &pr)
	if err != nil {
		return false, err
	}
	return pr.Merged, nil
}
