package state

// State holds information about the last poll against the API
// and approvals which are in process (i.e link is managing) when the application was terminated
type State struct {
	LastPollId int
}

func (s *State) CreateAndWrite() error {
	return nil
}
