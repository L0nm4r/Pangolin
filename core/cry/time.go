package cry

import "time"

const ErrorRange = 3 // 3s

func GenTimeStamp() int64 {
	return time.Now().Unix()
}

func CheckTimeStamp(timeStamp int64) bool {
	ct := time.Unix(timeStamp,0)
	return time.Now().Sub(ct) < ErrorRange* time.Second
}