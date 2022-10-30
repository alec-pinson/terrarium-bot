package main

import (
	"log"
	"time"
)

func Sunset() bool {
	dNow, _ := time.Parse("15:04", time.Now().Format("15:04"))
	dEnd, _ := time.Parse("15:04", c.Day.End)
	dSunset := dEnd.Add(-c.Day.Sunset)
	log.Println("Sunset(): Sunset time: " + dSunset.Format("15:04"))

	if dNow.After(dSunset) && dNow.Before(dEnd) {
		return true
	} else if dNow.Equal(dSunset) {
		return true
	} else {
		return false
	}
}

func Sunrise() bool {
	dNow, _ := time.Parse("15:04", time.Now().Format("15:04"))
	dStart, _ := time.Parse("15:04", c.Day.Start)
	dSunrise := dStart.Add(-c.Day.Sunrise)
	log.Println("Sunrise(): Sunrise time: " + dSunrise.Format("15:04"))

	if dNow.After(dSunrise) && dNow.Before(dStart) {
		return true
	} else if dNow.Equal(dSunrise) {
		return true
	} else {
		return false
	}
}

func DayTime() bool {
	dNow, _ := time.Parse("15:04", time.Now().Format("15:04"))
	dStart, _ := time.Parse("15:04", c.Day.Start)
	dEnd, _ := time.Parse("15:04", c.Day.End)

	if dNow.After(dStart) && dNow.Before(dEnd) {
		return true
	} else if dNow.Equal(dStart) {
		return true
	} else {
		return false
	}
}
