package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	var mxtRes sync.Mutex
	var waitGrp sync.WaitGroup
	var concurIter int64

	res = make([]user, n)
	chanLimit := make(chan struct{}, pool)

	for concurIter = 0; concurIter < n; concurIter++ {
		// add wait group
		waitGrp.Add(1)
		// write into chan empty struct
		chanLimit <- struct{}{}
		go func(numberUser int64) {
			defer waitGrp.Done()
			oneUser := getOne(numberUser)
			// write user datat to res slice
			mxtRes.Lock()
			res[numberUser] = oneUser
			mxtRes.Unlock()
			// clear channel :-)
			<-chanLimit
		}(concurIter)
	}
	waitGrp.Wait()
	return res
}

//func main() {
//	var userRes []user
//	var n, pool int64
//
//	n = 10
//	pool = 2
//
//	start := time.Now()
//	fmt.Printf("Current time: %v\n", start)
//	wantTime := start.Add(time.Duration(n/pool) * 100)
//	fmt.Printf("It wants time for execution: %v\n", wantTime)
//	userRes = getBatch(n, pool)
//	fmt.Printf("It real spend time on execution: %v\n", time.Since(start).Milliseconds())
//	for resIdx, oneRes := range userRes {
//		fmt.Printf("result, number %v, user ID: %v\n", resIdx, oneRes)
//	}
//}
