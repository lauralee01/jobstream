package fetcher

// Fetcher is an interface that represents a job provider (e.g. LinkedIn, Indeed).
// We use an interface here so that the JobService can fetch jobs from MANY sources
// without knowing the details of how each source works.
type Fetcher interface {
	// Name returns the name of the source (e.g. "LinkedIn")
	Name() string

	// Fetch returns a list of jobs from the source.
	// TODO: Define the method signature
	// Fetch() ([]domain.Job, error)
}

// TODO: Create a MockFetcher that implements this interface.
// A MockFetcher is useful for testing without making real network requests.
// type MockFetcher struct {}
// func (m *MockFetcher) Name() string { return "Mock" }
// func (m *MockFetcher) Fetch() ([]domain.Job, error) {
//    return []domain.Job{{Title: "Software Engineer Mock"}}, nil
// }
