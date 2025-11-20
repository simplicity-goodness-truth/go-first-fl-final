package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Calculation of next date
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	var nextDateString string
	var date time.Time

	// Getting rule type and value
	ruleType, ruleValue, err := getRepeatTypeAndValue(repeat)

	if err != nil {
		return "", err
	}

	//Getting start date
	startDate, err := time.Parse(apiDateFormat, dstart)

	if err != nil {
		return "", err
	}

	switch ruleType {

	case "d":

		// Getting repeat value of D type
		var repeatDValue, err = getDTypeRepeatValue(ruleValue)

		if err != nil {
			return "", err
		}

		// Searching for a next D date

		date = startDate
		
			for {
				date = date.AddDate(0, 0, repeatDValue)

				if afterNow(date, now) {
					break
				}
			}

	case "y":

		// Searching for a next Y date
		date = startDate

		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				break
			}
		}

	case "w":

		// Getting repeat value of W type
		var repeatWValue, err = getWTypeRepeatValue(ruleValue)

		if err != nil {
			return "", err
		}

		// Searching for a next W date
		date = now

		for {

			date = date.AddDate(0, 0, 1)

			// Validating if a calculated date from now is included into required weekdays array
			if isDayInWeekDaySet(date, repeatWValue) {
				break
			}

		}

	case "m":

		// Getting repeat value of M type
		var repeatMDays, repeatMMonths, err = getMTypeRepeatValue(ruleValue)

		if err != nil {
			return "", err
		}

		// Searching for a next M date
		date = startDate

		if len(repeatMMonths) > 0 {

			for {
				date = date.AddDate(0, 0, 1)

				if afterNow(date, now) && isDayInDaysAndMonthsSet(date, repeatMDays, repeatMMonths) {
					break
				}
			}

		} else {
			for {
				date = date.AddDate(0, 0, 1)

				if afterNow(date, now) && isDayInDaysSet(date, repeatMDays) {
					break
				}

			}

		}
	}

	nextDateString = date.Format(apiDateFormat)

	return nextDateString, nil

}

// Validation if specific date is in sets of provided days
func isDayInDaysSet(date time.Time, days []int) bool {

	var dateToCheck int

	for _, d := range days {

		if d > 0 {

			if date.Day() == d {
				return true
			}

		} else {

			// Calculation of a last day of a month
			if d == -1 {
				dateToCheck = date.AddDate(0, 1, -date.Day()).Day()
			}

			// Calculation of a day before a last day of a month
			if d == -2 {
				dateToCheck = date.AddDate(0, 1, -date.Day()-1).Day()
			}

			if date.Day() == dateToCheck {
				return true
			}

		}
	}

	return false
}

// Validation if specific date is in sets of provided days and months
func isDayInDaysAndMonthsSet(date time.Time, days []int, months []int) bool {

	for _, m := range months {

		if int(date.Month()) == m {

			if isDayInDaysSet(date, days) {
				return true
			}
		}
	}

	return false
}

// Validation if specific date is in set of provided week day numbers
func isDayInWeekDaySet(date time.Time, weekDaySet []int) bool {

	for _, weekDayNum := range weekDaySet {

		if int(date.Weekday()) == weekDayNum {

			return true

		}
	}

	return false
}

// Validation if time is later than now
func afterNow(date time.Time, now time.Time) bool {

	if date.After(now) {
		return true
	}

	return false

}

// Repeat type and value validator and parser
func getRepeatTypeAndValue(repeat string) (ruleType string, ruleValue string, err error) {

	var (
		errEmptyRepeat       = errors.New("repeat type is empty")
		errInvalidRepeatType = errors.New("invalid repeat type")
	)

	if repeat == "" {
		return "", "", errEmptyRepeat
	}

	// Parsing repeat string and returning rule type and rule value

	ruleType, ruleValue, _ = strings.Cut(repeat, " ")

	// Validation of set repeat type correctness
	if !isCorrectRepeatType(ruleType) {
		return "", "", errInvalidRepeatType
	}

	// Returning type and value, taking into an account, that
	// there could be no return value (for Y type, for instance)
	if len(ruleValue) > 0 {

		return ruleType, ruleValue, nil

	} else {

		return ruleType, "", nil

	}

}

// Get repeat value for D type
func getDTypeRepeatValue(repeatValue string) (daysCount int, err error) {

	var (
		errEmptyRepeatValue   = errors.New("repeat value is empty")
		errInvalidRepeatValue = errors.New("invalid repeat value")
	)

	// Empty return value is not permitted
	if repeatValue == "" {
		return 0, errEmptyRepeatValue
	}

	// Validations for D type

	// Repeat value must be a number
	repeatDValue, err := strconv.Atoi(repeatValue)

	if err != nil {
		return 0, errInvalidRepeatValue
	}

	// Repeat value must not exceed value of 400 or be less or equal to zero
	if repeatDValue <= 0 || repeatDValue > 400 {
		return 0, errInvalidRepeatValue
	}

	return repeatDValue, nil

}

// Get repeat value for W type
func getWTypeRepeatValue(repeatValue string) (weekDays []int, err error) {

	var (
		errEmptyRepeatValue   = errors.New("repeat value is empty")
		errInvalidRepeatValue = errors.New("invalid repeat value")
	)

	// Empty return value is not permitted
	if repeatValue == "" {
		return nil, errEmptyRepeatValue
	}

	repeatWValues := strings.Split(repeatValue, ",")

	// Maximum 7 values split by comma are permitted for W type
	if len(repeatWValues) > 7 {
		return nil, errInvalidRepeatValue
	}

	for i := 0; i < len(repeatWValues); i++ {

		// Single element cannot be greater than 7 for a W type
		var value, err = strconv.Atoi(repeatWValues[i])

		if err != nil {
			return nil, errInvalidRepeatValue
		}

		if value == 0 || value > 7 {
			return nil, errInvalidRepeatValue
		}

		// Specific case of Sunday - Sunday starts from zero

		if value == 7 {
			value = 0
		}

		weekDays = append(weekDays, value)
	}

	return weekDays, nil

}

// Get repeat value for M type
func getMTypeRepeatValue(repeatValue string) (days []int, months []int, err error) {

	var (
		errEmptyRepeatValue   = errors.New("repeat value is empty")
		errInvalidRepeatValue = errors.New("invalid repeat value")
	)

	// Empty return value is not permitted
	if repeatValue == "" {
		return nil, nil, errEmptyRepeatValue
	}

	repeatMDaysMonths := strings.Split(repeatValue, " ")

	// Maximum 2 values split by space are permitted for M type
	if len(repeatMDaysMonths) > 2 {
		return nil, nil, errInvalidRepeatValue
	}

	// Validating days
	repeatMDays := strings.Split(repeatMDaysMonths[0], ",")

	for i := 0; i < len(repeatMDays); i++ {

		// Single element cannot be greater than 31 for a W type
		var value, err = strconv.Atoi(repeatMDays[i])

		if err != nil {
			return nil, nil, errInvalidRepeatValue
		}

		if value == 0 || value > 31 || value < -2 {
			return nil, nil, errInvalidRepeatValue
		}

		days = append(days, value)
	}

	// Validating months if there are any

	if len(repeatMDaysMonths) > 1 {

		repeatMMonth := strings.Split(repeatMDaysMonths[1], ",")

		for i := 0; i < len(repeatMMonth); i++ {

			// Single element cannot be greater than 31 for a M type
			var value, err = strconv.Atoi(repeatMMonth[i])

			if err != nil {
				return nil, nil, errInvalidRepeatValue
			}

			if value == 0 || value > 12 {
				return nil, nil, errInvalidRepeatValue
			}

			months = append(months, value)

		}

		return days, months, nil

	} else {

		return days, nil, nil

	}

}

// Validation of repeat type correctness
func isCorrectRepeatType(repeatType string) bool {

	var correctRepeatTypes []string

	correctRepeatTypes = getCorrectRepeatTypes()

	for _, value := range correctRepeatTypes {
		if value == repeatType {
			return true
		}
	}
	return false
}

// Get function for permitted repeat types
func getCorrectRepeatTypes() []string {

	return []string{"d", "y", "w", "m"}

}

// Handler for next date calculation
func nextDayHandler(res http.ResponseWriter, req *http.Request) {

	// Only GET method is permitted
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Getting request parameters
	nowString := req.FormValue("now")
	dateString := req.FormValue("date")
	repeatString := req.FormValue("repeat")

	now, err := time.Parse(apiDateFormat, nowString)

	if err != nil {
		http.Error(res, "Failed to parse now value:"+err.Error(), http.StatusBadRequest)
		return
	}

	// Calculating next date
	nextDate, err := NextDate(now, dateString, repeatString)

	if err != nil {
		http.Error(res, "Processing error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Setting a successfull HTTP status
	res.WriteHeader(http.StatusOK)

	// Sending a response
	res.Write([]byte(nextDate))
}
