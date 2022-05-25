package cry

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	ts := GenTimeStamp()
	//time.Sleep(3*time.Second)

	b := UInt64ToBytes(uint64(ts))
	tt := BytesToUInt64(b)

	fmt.Println(CheckTimeStamp(int64(tt)))
	time.Sleep(3*time.Second)
	fmt.Println(CheckTimeStamp(int64(tt)))
}