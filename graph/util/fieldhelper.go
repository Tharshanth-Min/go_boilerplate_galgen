package util

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"example/graph/logengine"
	suuid "github.com/gofrs/uuid"
)

func UUIDV4ToString(id [16]byte) string {
	u, err := suuid.FromBytes(id[:])
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	return u.String()
}

func IsValidDateWithUSDateString(date string) (time.Time, string, error) {
	const (
		layoutISO = "01/02/2006"
		layoutUS  = "2006/01/02"
	)
	timeObj, err := time.Parse(layoutISO, date)
	if err != nil {
		return timeObj, "", err
	}
	return timeObj, timeObj.Format(layoutUS), err
}

func DateDifference(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func BytesToUUIDV4(id [16]byte) suuid.UUID {
	u, err := suuid.FromBytes(id[:])
	if err != nil {
		logengine.GetTelemetryClient().TrackException(errors.New("Failed to parse UUID from byte array"))
	}
	return u
}

func StringToUUID4(uuidstr string) suuid.UUID {
	uid, err := suuid.FromString(uuidstr)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	return uid
}

func CheckStringExceedLimit(input string, limit int) bool {
	result := false
	if utf8.RuneCountInString(input) > limit {
		result = true
	}
	return result
}

func getUniqueStringArray(stringSlice []string) []string {
	keys := make(map[string]int)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = 1
			list = append(list, entry)
		}
	}
	return list
}

func GetUSLocaleTimeString(date time.Time) string {
	return date.Format("01/02/2006")
}

func GetDateFromUTCTime(date string) string {
	var formatedDate string
	if date != "" {
		parse_time, _ := time.Parse(time.RFC3339, date)
		formatedDate = parse_time.Format("01/02/2006")
	}
	return formatedDate
}

func GetSGLocaleTimeString(date time.Time) string {
	loc, err := time.LoadLocation("Singapore")
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	date = date.In(loc)
	return date.Format("02/01/2006T15:04")
}
func GetTimeUnixTimeStamp(date time.Time) string {
	timestamp := int32(date.Unix())
	return String(timestamp)
}

func AppendStringUnique(list *[]interface{}, idmap map[string]int, idstr string) {
	if _, value := idmap[idstr]; !value {
		idmap[idstr] = 1
		*list = append(*list, idstr)
	}
}

// Convert int32 to string
func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
func GetCurrentTime() time.Time {
	timeNow := time.Now().UTC().Add(0 * time.Hour)
	return timeNow
}

func GetInteger64Length(value int64) (count int) {
	for value != 0 {
		value /= 10
		count = count + 1
	}
	return count
}

func GetIntegerLength(value int) (count int) {
	for value != 0 {
		value /= 10
		count = count + 1
	}
	return count
}

func IsValidEmail(email string) bool {
	Re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return Re.MatchString(email)
}

func RemoveDuplicatesFromSlice(s []string) []string {
	mapOfStrings := make(map[string]bool)
	for _, item := range s {
		_, ok := mapOfStrings[item]
		if !ok {
			mapOfStrings[item] = true
		}
	}
	var result []string
	for item, _ := range mapOfStrings {
		result = append(result, item)
	}
	return result
}

func IsValidDateWithDateObect(date string) (time.Time, error) {
	const (
		layoutISO = "01/02/2006"
		layoutUS  = "2 January, 2006"
	)
	timeObj, err := time.Parse(layoutISO, date)
	if err != nil {
		return timeObj, err
	}
	return timeObj, err
}

func IsValidDate(date string) error {
	const (
		layoutISO = "01/02/2006"
		layoutUS  = "2 January, 2006"
	)
	_, err := time.Parse(layoutISO, date)
	if err != nil {
		return err
	}
	return err
}

func CompareDate(startDate string, endDate string) bool {
	startDateArr := strings.Split(startDate, "/")
	endDateArr := strings.Split(endDate, "/")
	if startDateArr[2] > endDateArr[2] {
		return false
	} else if startDateArr[2] == endDateArr[2] && startDateArr[0] > endDateArr[0] {
		return false
	} else if startDateArr[2] == endDateArr[2] && startDateArr[0] == endDateArr[0] && startDateArr[1] > endDateArr[1] {
		return false
	} else {
		return true
	}
}

func CurrentDate() string {
	currentTime := time.Now()
	date := strings.ReplaceAll(currentTime.Format("01-02-2006"), "-", "/")
	return date
}

func GetStartDateEndDate(validityDate string) (string, string) {
	if validityDate != "" {
		replacer := strings.NewReplacer("[", "", ")", "")

		validityDate := replacer.Replace(validityDate)
		dateRange := strings.Split(validityDate, ",")
		if len(validityDate) < 21 { //When start date and end date same
			return dateRange[0], dateRange[0]
		}
		return dateRange[0], dateRange[1]
	}
	return "", ""
}

func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func ValidateMonthYear() (int, int) {
	year, month, _ := time.Now().Date()

	return int(month), int(year)
}

func ValidateURL(URL string) (bool, string) {
	var valid bool = false
	var fileName string
	if strings.Contains(strings.ToLower(URL), "http") || strings.Contains(strings.ToLower(URL), "https") {
		valid = true
		splitedURL := strings.Split(URL, "/")
		length := len(splitedURL)
		fileName = splitedURL[length-1]
	}
	return valid, fileName
}

func IsCurrentPastMonthYear(month string, year string) bool {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	currentDateTime := firstOfMonth.AddDate(0, 1, 0).Format("01/2006")
	currentTime, _ := time.Parse("01/2006", currentDateTime)

	var m, _ = strconv.Atoi(month)
	var y, _ = strconv.Atoi(year)
	comparisonDateTime := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC).Format("01/2006")
	comparisonDate, _ := time.Parse("01/2006", comparisonDateTime)

	result := comparisonDate.Before(currentTime)
	return result
}

func GetFirstDayOfCurrentMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	return firstOfMonth
}
