package edenutil

import (
	"fmt"
	"time"
)

// These are a collection of utilities for working with Time as it relates to the world of Eden.

// EdenMonth is an int type that represents a month in the world of Eden.
// An Eden Calendar Month is broken into 6 weeks, with a week being broken into 12 days.
type EdenMonth int64

// An Eden Calendar Year is broken into 18 months.
const (
	EdenMonthCapris EdenMonth = iota
	EdenMonthYuena
	EdenMonthGabriel
	EdenMonthZebulon
	EdenMonthYul
	EdenMonthEden
	EdenMonthMikael
	EdenMonthLeonis
	EdenMonthAvrila
	EdenMonthSedrah
	EdenMonthBonafu
	EdenMonthVenu
	EdenMonthTavros
	EdenMonthSycorax
	EdenMonthKarkat
	EdenMonthTerezi
	EdenMonthVriska
	EdenMonthEquius
)

func (em EdenMonth) String() string {
	switch em {
	case EdenMonthCapris:
		return "Capris"
	case EdenMonthYuena:
		return "Yuena"
	case EdenMonthGabriel:
		return "Gabriel"
	case EdenMonthZebulon:
		return "Zebulon"
	case EdenMonthYul:
		return "Yul"
	case EdenMonthEden:
		return "Eden"
	case EdenMonthMikael:
		return "Mikael"
	case EdenMonthLeonis:
		return "Leonis"
	case EdenMonthAvrila:
		return "Avrila"
	case EdenMonthSedrah:
		return "Sedrah"
	case EdenMonthBonafu:
		return "Bonafu"
	case EdenMonthVenu:
		return "Venu"
	case EdenMonthTavros:
		return "Tavros"
	case EdenMonthSycorax:
		return "Sycorax"
	case EdenMonthKarkat:
		return "Karkat"
	case EdenMonthTerezi:
		return "Terezi"
	case EdenMonthVriska:
		return "Vriska"
	case EdenMonthEquius:
		return "Equius"
	default:
		return ""
	}
}

// EdenDay is an int type that represents a day in the world of Eden.
// An Eden Day is broken into 27 hours, with an hour being broken into 60 minutes, a minute 60 seconds, etc.
type EdenDay int64

const (
	EdenDaySandu EdenDay = iota
	EdenDayMoudu
	EdenDayTudu
	EdenDayWendu
	EdenDayThurdu
	EdenDayFradu
	EdenDaySadu
	EdenDayLandu
	EdenDayZedu
	EdenDayKadu
	EdenDayVedu
	EdenDayBedu
)

func (ed EdenDay) String() string {
	switch ed {
	case EdenDaySandu:
		return "Sandu"
	case EdenDayMoudu:
		return "Moudu"
	case EdenDayTudu:
		return "Tudu"
	case EdenDayWendu:
		return "Wendu"
	case EdenDayThurdu:
		return "Thurdu"
	case EdenDayFradu:
		return "Fradu"
	case EdenDaySadu:
		return "Sadu"
	case EdenDayLandu:
		return "Landu"
	case EdenDayZedu:
		return "Zedu"
	case EdenDayKadu:
		return "Kadu"
	case EdenDayVedu:
		return "Vedu"
	case EdenDayBedu:
		return "Bedu"
	default:
		return ""
	}
}

// EdenTime is a type that is "similar", a term used loosely, to time.Time, but as relates to the world of Eden.
// The creation time of the Eden Universe begins at 12:00AM January 1st, 2021 UTC.
// It is relative to the home capital city of Eden, Freeport.
type EdenTime struct {
}

// NanoSecond Returns the nanoseconds since the creation of the Eden Universe.
func (et EdenTime) NanoSecond() int64 {
	// Current nanoseconds since 12:00AM January 1st, 2021 UTC
	nanoSecond := time.Since(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Nanoseconds()
	return nanoSecond
}

// MicroSecond Returns the microseconds since the creation of the Eden Universe.
func (et EdenTime) MicroSecond() int64 {
	// Current microseconds since 12:00AM January 1st, 2021 UTC
	microSecond := time.Since(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds()
	return microSecond
}

// Second Returns the seconds since the creation of the Eden Universe.
func (et EdenTime) Second() int64 {
	// Current seconds since 12:00AM January 1st, 2021 UTC
	second := int64(time.Since(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Seconds())
	return second
}

// Minute Returns the minutes since the creation of the Eden Universe.
func (et EdenTime) Minute() int64 {
	// Current minutes since 12:00AM January 1st, 2021 UTC
	minute := int64(time.Since(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Minutes())
	return minute
}

// Hour Returns the hours since the creation of the Eden Universe.
func (et EdenTime) Hour() int64 {
	// Current hours since 12:00AM January 1st, 2021 UTC
	hour := int64(time.Since(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Hours())
	return hour
}

// Day returns the days since the creation of the Eden Universe
// A day is broken into 27 hours, with an hour being broken into 60 minutes, a minute 60 seconds, etc.
func (et EdenTime) Day() int64 {
	// Current days since 12:00AM January 1st, 2021 UTC
	day := int64(time.Since(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Hours() / 24)
	return day
}

// Week returns the weeks since the creation of the Eden Universe
// A week is broken into 12 Eden days, with an Eden day being broken into 27 hours, with an hour being broken
// into 60 minutes, a minute 60 seconds, etc.
func (et EdenTime) Week() int64 {
	// Current weeks since 12:00AM January 1st, 2021 UTC
	week := et.Day() / 12
	return week
}

// Month returns the months since the creation of the Eden Universe
// A month is broken into 6 weeks, with a week being broken into 12 days.
func (et EdenTime) Month() int64 {
	// Current months since 12:00AM January 1st, 2021 UTC
	month := et.Week() / 6
	return month
}

// Year returns the years since the creation of the Eden Universe
// A year is broken into 18 months.
func (et EdenTime) Year() int64 {
	// Current years since 12:00AM January 1st, 2021 UTC
	year := et.Month() / 18
	return year
}

// EdenTimestamp returns the current hour, minute and second formatted as HH:MM:SS
func (et EdenTime) EdenTimestamp() (second int64, minute int64, hour int64) {
	// Current hour, minute and second since 12:00AM January 1st, 2021 UTC
	second = et.Second()
	minute = (second / 60)
	hour = (minute / 60)
	return second % 60, minute % 60, hour % 27
}

// EdenMonth Returns the current EdenMonth
func (et EdenTime) EdenMonth() EdenMonth {
	// The current month of the Eden Calendar
	month := et.Month() % 18
	return EdenMonth(month)
}

// EdenDay Returns the current EdenDay
func (et EdenTime) EdenDay() EdenDay {
	// The current day of the Eden Calendar
	day := et.Day() % 12
	return EdenDay(day)
}

// EdenDayString returns the current EdenDay as a string
func (et EdenTime) EdenDayString() string {
	// The current day of the Eden Calendar
	day := et.EdenDay()
	return day.String()
}

// CurrentTimeString returns the current time formatted as HH:MM:SS
func (et EdenTime) CurrentTimeString() string {
	second, minute, hour := et.EdenTimestamp()
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

// TimeStamp String returns the current day time formatted as "<day> of the <week> Week, <month> in the Year <year> - HH:MM:SS"
func (et EdenTime) TimeStamp() string {
	// The current day of the Eden Calendar
	day := et.EdenDay().String()
	// The current week of the Eden Calendar
	week := et.Week() % 6
	// The current month of the Eden Calendar
	month := et.EdenMonth().String()
	// The current year of the Eden Calendar
	year := et.Year()
	// The current hour, minute and second since 12:00AM January 1st, 2021 UTC
	second, minute, hour := et.EdenTimestamp()

	return fmt.Sprintf("%s of Week %d in %s, Year %d - %02d:%02d:%02d", day, week, month, year, hour, minute, second)
}
