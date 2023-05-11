package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"yandex-team.ru/bstask/db/schemas"
)

type rowsScanner struct {
	rows *sql.Rows
}

type rowScanner struct {
	row *sql.Row
}

func (scanner *rowsScanner) ScanCourierData() (schemas.CourierDto, error) {
	var res schemas.CourierDto
	var regions []uint8
	var workingHours []uint8

	err := scanner.rows.Scan(&res.CourierId, &res.CourierType, &regions, &workingHours)
	if err != nil {
		log.Fatal(err)
	}

	res.Regions, err = stringToInt32Slice(bytesToString(regions))
	res.WorkingHours = stringToStringSlice(bytesToString(workingHours))

	return res, err
}

func (scanner *rowScanner) ScanCourierData() (schemas.CourierDto, error) {
	var res schemas.CourierDto
	var regions []uint8
	var workingHours []uint8

	err := scanner.row.Scan(&res.CourierId, &res.CourierType, &regions, &workingHours)
	if err != nil {
		log.Fatal(err)
	}

	res.Regions, err = stringToInt32Slice(bytesToString(regions))
	res.WorkingHours = stringToStringSlice(bytesToString(workingHours))

	return res, err
}

func (scanner *rowScanner) ScanOrderData() (schemas.OrderDto, error) {
	var res schemas.OrderDto
	var deliveryHours []uint8
	var completedTime []uint8

	err := scanner.row.Scan(&res.OrderId, &res.Weight, &res.Region, &deliveryHours, &res.Cost, &completedTime)
	if err != nil {
		log.Fatal(err)
	}

	res.DeliveryHours = stringToStringSlice(bytesToString(deliveryHours))
	res.CompletedTime = bytesToString(completedTime)

	return res, err
}

func (scanner *rowsScanner) ScanOrderData() (schemas.OrderDto, error) {
	var res schemas.OrderDto
	var deliveryHours []uint8
	var completedTime []uint8

	err := scanner.rows.Scan(&res.OrderId, &res.Weight, &res.Region, &deliveryHours, &res.Cost, &completedTime)
	if err != nil {
		log.Fatal(err)
	}

	res.DeliveryHours = stringToStringSlice(bytesToString(deliveryHours))
	res.CompletedTime = bytesToString(completedTime)

	return res, err
}

func bytesToString(data []uint8) string {
	return string(data)
}

func stringToStringSlice(str string) []string {

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	arr := strings.Split(str, ",")

	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
		arr[i] = strings.Trim(arr[i], `"`)
	}

	return arr
}

func stringToInt32Slice(str string) ([]int32, error) {

	arr := stringToStringSlice(str)
	intArr := make([]int32, len(arr))

	for i, elem := range arr {
		val, err := strconv.ParseInt(elem, 10, 32)
		if err != nil {
			return nil, err
		}
		intArr[i] = int32(val)
	}

	return intArr, nil
}

func int32SliceToUint8Array(intSlice []int32) []byte {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, val := range intSlice {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%d", val))
	}
	buf.WriteString("}")
	return buf.Bytes()
}

func stringSliceToUint8Array(stringSlice []string) []byte {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, val := range stringSlice {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf(`"%s"`, val))
	}
	buf.WriteString("}")
	return buf.Bytes()
}
