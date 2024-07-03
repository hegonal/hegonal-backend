package utils

import "fmt"

import (
	"time"

	"github.com/godruoyi/go-snowflake"
)

func SnowFlakeInit() {
	snowflake.SetMachineID(1)
	snowflake.SetStartTime(time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC))
}

func GenerateId() string {
    orderId := snowflake.ID()
	return fmt.Sprint(orderId)
}


