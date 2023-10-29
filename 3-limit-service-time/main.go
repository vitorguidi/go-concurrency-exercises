//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var timeAllowance = make(map[int]time.Duration)
var timeAllowanceMut = &sync.Mutex{}
var freeTierTime = time.Second * 10

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	//simplifying hypothesis: every goroutine will complete without being preempted
	//otherwise this gets out of hand, how can I keep track of in CPU time? o.O
	processingDone := make(chan time.Duration)

	timeAllowanceMut.Lock()
	timeSlice, isUserPresent := timeAllowance[u.ID]
	if !isUserPresent {
		timeAllowance[u.ID] = freeTierTime
		timeSlice = freeTierTime
	}
	timeAllowanceMut.Unlock()

	ticker := time.NewTicker(timeSlice)

	go func() {
		startTime := time.Now()
		process()
		processingDone <- time.Since(startTime)
	}()

	for {
		select {
		case <-ticker.C:
			timeAllowanceMut.Lock()
			timeSlice = timeAllowance[u.ID] - freeTierTime
			if timeSlice < 0*time.Second {
				timeSlice = 0
			}
			timeAllowance[u.ID] = timeSlice
			timeAllowanceMut.Unlock()
			//should return false if process gets throttled
			//premium (true) does not, free (false) gets
			shouldGetThrottled := !u.IsPremium
			return !shouldGetThrottled
		case processingDuration := <-processingDone:
			timeAllowanceMut.Lock()
			timeSlice = timeAllowance[u.ID]
			exceededFreeLimit := timeSlice < processingDuration
			if exceededFreeLimit {
				timeSlice = 0 * time.Second
			} else {
				timeSlice -= processingDuration
			}
			timeAllowance[u.ID] = timeSlice
			timeAllowanceMut.Unlock()
			shouldGetThrottled := exceededFreeLimit && !u.IsPremium
			//weird semantics. if someone should get throttled, we return false
			return !shouldGetThrottled
		}
	}

}

func main() {
	RunMockServer()
}
