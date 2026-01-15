package functions

func saveHistory(msg string) {
	mu.Lock()
	history = append(history, msg)
	mu.Unlock()
}
