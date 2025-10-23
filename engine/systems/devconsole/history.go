package devconsole

func (s *ConsoleState) NavigateHistory(direction int) {
    if len(s.History) == 0 {
        return
    }

    switch {
    case direction < 0: // Up
        if s.HistoryIdx == -1 {
            s.HistoryIdx = len(s.History) - 1
        } else if s.HistoryIdx > 0 {
            s.HistoryIdx--
        }
    case direction > 0: // Down
        if s.HistoryIdx == -1 {
            return
        }
        if s.HistoryIdx < len(s.History)-1 {
            s.HistoryIdx++
        } else {
            s.HistoryIdx = -1
            s.InputBuffer = ""
            return
        }
    }

    if s.HistoryIdx >= 0 && s.HistoryIdx < len(s.History) {
        s.InputBuffer = s.History[s.HistoryIdx]
    }
}

func (s *ConsoleState) PushHistory(cmd string) {
    if len(s.History) >= maxHistoryStored {
        s.History = s.History[1:]
    }
    s.History = append(s.History, cmd)
}

