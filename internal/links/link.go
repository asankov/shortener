package links

type Link struct {
	ID      string
	URL     string
	Metrics *Metrics
}

type Metrics struct {
	Clicks int
}
