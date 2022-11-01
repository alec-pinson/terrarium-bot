package main

import (
	"log"
	"time"
)

var (
	dayStartTime time.Time
	dayEndTime   time.Time
	sunriseTime  time.Time
	sunsetTime   time.Time
)

func GetSunriseTime() time.Time {
	dayStartTime, _ = time.Parse("15:04", c.Day.Start)
	sunriseTime = dayStartTime.Add(-c.Day.Sunrise)
	log.Println("Sunrise time: " + sunriseTime.Format("15:04"))
	return sunriseTime
}

func GetSunsetTime() time.Time {
	dayEndTime, _ = time.Parse("15:04", c.Day.End)
	sunsetTime = dayEndTime.Add(-c.Day.Sunset)
	log.Println("Sunset time: " + sunsetTime.Format("15:04"))
	return sunsetTime
}

func Sunset() bool {
	nowTime, _ := time.Parse("15:04", time.Now().Format("15:04"))

	if nowTime.After(sunsetTime) && nowTime.Before(dayEndTime) {
		return true
	} else if nowTime.Equal(sunsetTime) {
		return true
	} else {
		return false
	}
}

func Sunrise() bool {
	nowTime, _ := time.Parse("15:04", time.Now().Format("15:04"))

	if nowTime.After(sunriseTime) && nowTime.Before(dayStartTime) {
		return true
	} else if nowTime.Equal(sunriseTime) {
		return true
	} else {
		return false
	}
}

func DayTime() bool {
	nowTime, _ := time.Parse("15:04", time.Now().Format("15:04"))

	if nowTime.After(dayStartTime) && nowTime.Before(dayEndTime) {
		return true
	} else if nowTime.Equal(dayStartTime) {
		return true
	} else {
		return false
	}
}
