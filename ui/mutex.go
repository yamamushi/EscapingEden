package ui

import "github.com/yamamushi/EscapingEden/logging"

func (c *Console) LockMainMutex(caller string) {
	c.Log.Println(logging.LogDebug, "Locking main Console mutex from: ", caller)
	c.mutex.Lock()
}

func (c *Console) UnlockMainMutex(caller string) {
	c.Log.Println(logging.LogDebug, "Unlocking main Console mutex from: ", caller)
	c.mutex.Unlock()
}
