package dockerlib

// EnsureImageProgressDetail is an object to help unmarshall JSON returned from Docker during a pull.
type EnsureImageProgressDetail struct {
	Current int
	Total   int
}

// EnsureImageProgress is an object to unmarshall JSON returned from Docker during a pull.
type EnsureImageProgress struct {
	Status         string
	ProgressDetail EnsureImageProgressDetail
	Progress       string
	ID             string
}

func (p EnsureImageProgress) String() string {
	if len(p.ID) > 0 {
		return p.ID + " " + p.Status + " " + p.Progress
	} else {
		return p.Status
	}
}
