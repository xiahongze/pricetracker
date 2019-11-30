package trackers

// Tracker is the type future implementation should follow
type Tracker func(url, xpath *string) (string, error)
