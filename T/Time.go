package T

// Time support package
import (
	"time"

	"github.com/kokizzu/rand"

	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
)

const ISO = `2006-01-02T15:04:05.999999`
const YMD_HM = `2006-01-02 15:04`
const YMD_HMS = `2006-01-02 15:04:05`
const YMD = `2006-01-02`
const FILE = `20060102_150405`
const HUMAN = `2-Jan-2006 15:04:05`
const HUMAN_DATE = `2 Jan 2006`
const YY = `06`
const YMDH = `20060102.15`
const YMDHM = `20060102.1504`
const HMS = `150405`

var EMPTY = time.Time{}

// ToIsoStr convert time to iso formatted time string
//
//	T.ToIsoStr(time.Now()) // "2016-03-17T10:04:50.6489"
func ToIsoStr(t time.Time) string {
	if t.Equal(EMPTY) {
		return ``
	}
	return t.Format(ISO)
}

// IsoStr current iso time
//
//	T.IsoStr() // "2016-03-17T10:07:56.418728"
func IsoStr() string {
	return time.Now().Format(ISO)
}

// ToDateStr convert time to iso date
//
//	T.ToDateStr(time.Now()) // output "2016-03-17"
func ToDateStr(t time.Time) string {
	if t.Equal(EMPTY) {
		return ``
	}
	return t.Format(YMD)
}

// DateStr current iso date
// T.DateStr()) // "2016-03-17"
func DateStr() string {
	return time.Now().Format(YMD)
}

// ToHumanStr convert time to human date
//
//	T.ToHumanStr(time.Now()) // "17-Mar-2016 10:06"
func ToHumanStr(t time.Time) string {
	if t.Equal(EMPTY) {
		return ``
	}
	return t.Format(HUMAN)
}

// HumanStr current human date
//
//	T.HumanStr() // "17-Mar-2016 10:06"
func HumanStr() string {
	return time.Now().Format(HUMAN)
}

// ToDateHourStr convert time to iso date and hour:minute
//
//	T.ToDateHourStr(time.Now()) // "2016-03-17 10:07"
func ToDateHourStr(t time.Time) string {
	if t.Equal(EMPTY) {
		return ``
	}
	return t.Format(YMD_HM)
}

// ToHhmmssStr convert time to iso date and hourminutesecond
//
//	T.ToDateHourStr(time.Now()) // "230744"
func ToHhmmssStr(t time.Time) string {
	if t.Equal(EMPTY) {
		return ``
	}
	return t.Format(HMS)
}

// DateHhStr current iso date and hour
//
//	T.DateHhStr()// output "20160317.10"
func DateHhStr() string {
	return time.Now().Format(YMDH)
}

// DateHhMmStr current iso date and hour
//
//	T.DateHhMmStr()// output "20160317.1059"
func DateHhMmStr() string {
	return time.Now().Format(YMDHM)
}

// ToDateTimeStr convert time to iso date and time
//
//	T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"
func ToDateTimeStr(t time.Time) string {
	if t.Equal(EMPTY) {
		return ``
	}
	return t.Format(YMD_HMS)
}

// DateTimeStr current iso date and time
//
//	T.ToDateTimeStr(time.Now()) // "2016-03-17 10:07:50"
func DateTimeStr() string {
	return time.Now().Format(YMD_HMS)
}

// DayInt int64 day of current date
func DayInt() int64 {
	return int64(time.Now().Day())
}

// HourInt int64 current hour
func HourInt() int64 {
	return int64(time.Now().Hour())
}

// MonthInt int64 current month
func MonthInt() int64 {
	return int64(time.Now().Month())
}

// YearInt int64 current year
func YearInt() int64 {
	return int64(time.Now().Year())
}

// YearDayInt int64 current day of year
func YearDayInt() int64 {
	return int64(time.Now().YearDay())
}

// Filename get filename version of current date
//
//	T.Filename()) // "20160317_102543"
func Filename() string {
	return time.Now().Format(FILE)
}

// HhmmssStr get filename version of current time
func HhmmssStr() string {
	return ToHhmmssStr(time.Now())
}

// Sleep delay for nanosecond
func Sleep(ns time.Duration) {
	time.Sleep(ns)
}

// RandomSleep random 0.4-2 sec sleep
func RandomSleep() {
	dur := rand.Uint64()%(1600*1000*1000) + (400 * 1000 * 1000)
	time.Sleep(time.Duration(dur))
}

// Track measure elapsed time in nanosec
//
//	T.Track(func(){
//	  x:=0
//	  T.Sleep(1)
//	}) // "done in 1.00s"
func Track(fun func()) time.Duration {
	start := time.Now()
	fun()
	elapsed := time.Since(start)
	L.ParentDescribe(`done in ` + F.ToStr(elapsed.Seconds()) + `s`)
	return elapsed
}

// IsValidTimeRange check if time in are in the range
//
//	t1, _:=time.Parse(`1992-03-23`,T.DateFormat)
//	t2, _:=time.Parse(`2016-03-17`,T.DateFormat)
//	T.IsValidTimeRange(t1,t2,time.Now()) // bool(false)
func IsValidTimeRange(start, end, check time.Time) bool {
	res := check.After(start) && check.Before(end)
	return res
}

// Age returns age from current date
func Age(birthdate time.Time) float64 {
	return float64(time.Since(birthdate)/time.Hour) / 24 / 365.25
}

// AgeAt returns age from within 2 date
func AgeAt(birthdate, point time.Time) float64 {
	return float64(point.Sub(birthdate)/time.Hour) / 24 / 365.25
}

// UnixNano get current unix nano
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// UnixNanoAfter get current unix nano after added with certain duration
func UnixNanoAfter(d time.Duration) int64 {
	return time.Now().Add(d).UnixNano()
}

// Epoch get current unix (second) as integer
func Epoch() int64 {
	return time.Now().Unix()
}

// ToEpoch convert string date to epoch => '2019-01-01' -->1546300800
func ToEpoch(date string) int64 {
	d, err := time.Parse(YMD, date)
	if err != nil {
		return 0
	} else {
		return d.Unix()
	}
}

// EpochStr get current unix (second) as string
func EpochStr() string {
	return I.ToS(time.Now().Unix())
}

// EpochAfter get current unix time added with a duration
func EpochAfter(d time.Duration) int64 {
	return time.Now().Add(d).Unix()
}

// EpochAfterStr get current unix time added with a duration
func EpochAfterStr(d time.Duration) string {
	return I.ToS(time.Now().Add(d).Unix())
}

// UnixToFile convert unix time to file naming
func UnixToFile(i int64) string {
	return time.Unix(i, 0).Format(FILE)
}

// WeekdayStr get day's name
func WeekdayStr() string {
	return time.Now().Weekday().String()
}

// Weekday get what day is it today, Sunday => 0
func Weekday() int {
	return int(time.Now().Weekday())
}

// UnixToDateTimeStr convert unix seconds to YYYY-MM-DD_hh:mm:ss
func UnixToDateTimeStr(epoch float64) string {
	return time.Unix(int64(epoch), 0).Format(YMD_HMS)
}

// UnixToDateStr convert from unix sconds to YYYY-MM-DD
func UnixToDateStr(epoch float64) string {
	return time.Unix(int64(epoch), 0).Format(YMD)
}

// UnixToHumanDateStr convert from unix to human date format D MMM YYYY
func UnixToHumanDateStr(epoch float64) string {
	return time.Unix(int64(epoch), 0).Format(HUMAN_DATE)
}

// UnixToHumanStr convert from unix to human format D-MMM-YYYY hh:mm:ss
func UnixToHumanStr(epoch float64) string {
	return time.Unix(int64(epoch), 0).Format(HUMAN)
}

// LastTwoDigitYear return current last two digit year
func LastTwoDigitYear() string {
	return time.Now().Format(YY)
}
