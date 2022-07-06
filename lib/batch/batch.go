package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func numberOfInteration(numberOfsubs int64, poolOfRequest int64) (iterationFirst int64, iterationSecond int64, remainderOfdivision int64, err error) {

	// check input vars
	if numberOfsubs <= 0 || poolOfRequest <= 0 {
		return 0, 0, 0, fmt.Errorf("input parameters should more than zero")
	} else {
		if poolOfRequest >= numberOfsubs {
			return 1, numberOfsubs, 0, nil
		} else {
			return numberOfsubs / poolOfRequest, poolOfRequest, numberOfsubs % poolOfRequest, nil
		}
	}
}

func getBatch(n int64, pool int64) (res []user) {
	var concurIter, poolIter int64
	var numberUser int64
	//var resIdx int64

	// check input vars
	iterationNum, iterationPool, remainderOfDivision, err := numberOfInteration(n, pool)
	if err == nil {
		fmt.Errorf("Input parameters has been error:%w", err)
	}
	//fmt.Printf("iteration number: %v\n", iterationNum)
	//fmt.Printf("iteration pool: %v\n", iterationPool)
	//fmt.Printf("the remainder of the division: %v\n", remainderOfDivision)
	// add the last iteration with not completelly pool (the remainder of the division)
	if remainderOfDivision > 0 {
		iterationNum++
	}
	// goroutin for collect of result by channel
	// new conext background
	ctx := context.Background()
	// start goroutines with collect users
	// at first cancurency should be implemented
	//randSource := rand.New(rand.NewSource((concurIter + 1) * (poolIter + 1)))
	numberUser = 0
	for concurIter = 0; concurIter < iterationNum; concurIter++ {
		// check there is the remainder of the division
		//fmt.Printf("concurIter: %v\n", concurIter)
		if concurIter == iterationNum-1 {
			if remainderOfDivision > 0 {
				iterationPool = remainderOfDivision
			}
		}
		// add context group
		userChan := make(chan user, 1)
		errGroup, _ := errgroup.WithContext(ctx)
		for poolIter = 0; poolIter < iterationPool; poolIter++ {
			// start of goroutine
			// time mesuares
			//timeSince := time.Now()
			errGroup.Go(func() (err error) {
				var mtx sync.Mutex

				mtx.Lock()
				// get user
				usr := getOne(numberUser)
				mtx.Unlock()
				userChan <- usr
				return nil
			})
			// read chan
			oneUserFromChan := <-userChan
			//fmt.Printf("time diff: %v ms\n", time.Since(timeSince).Milliseconds())
			res = append(res, oneUserFromChan)
			// number user increase
			numberUser++
		}
	}
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
//	fmt.Printf("It wants time for execution: %v", wantTime)
//	userRes = getBatch(n, pool)
//	fmt.Printf("It real spend time on execution: %v\n", time.Since(start).Milliseconds())
//	for resIdx, oneRes := range userRes {
//		fmt.Printf("result, number %v, user ID: %v\n", resIdx, oneRes)
//	}
//}
