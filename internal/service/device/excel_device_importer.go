package device

import (
	"aed-api-server/internal/interfaces/entities"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"strings"
	"time"
)

type ExcelDeviceImporter struct {
}

func NewExcelDeviceImporter() IDeviceImporter {
	return &ExcelDeviceImporter{}
}

func (s ExcelDeviceImporter) ImportDevices(reader io.Reader) ([]*entities.BaseDevice, error) {
	r, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	sheets := r.GetSheetList()

	return s.parseSheetsDevices(r, sheets)
}

func (s ExcelDeviceImporter) parseSheetsDevices(excel *excelize.File, sheets []string) ([]*entities.BaseDevice, error) {
	var devices []*entities.BaseDevice

	for _, sheet := range sheets {
		d, err := s.parseSheetDevices(excel, sheet)
		if err != nil {
			return nil, err
		}

		devices = append(devices, d...)
	}

	return devices, nil
}

func (s ExcelDeviceImporter) parseSheetDevices(excel *excelize.File, sheet string) ([]*entities.BaseDevice, error) {
	rows, err := excel.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	return s.parseDevices(rows)
}

func (s ExcelDeviceImporter) parseDevices(rows [][]string) ([]*entities.BaseDevice, error) {
	var devices []*entities.BaseDevice

	for i, row := range rows {
		if i == 0 {
			continue
		}

		d, err := s.parseDevice(rows[0], row)
		if err != nil {
			return nil, err
		}

		if d == nil {
			continue
		}

		devices = append(devices, d)
	}

	return devices, nil
}

func (s ExcelDeviceImporter) parseDevice(head []string, row []string) (*entities.BaseDevice, error) {
	var device entities.BaseDevice

	device.Id = uuid.NewString()
	device.Created = time.Now().UnixMilli()
	device.Source = entities.DeviceSourceImported

	for i, r := range row {
		switch head[i] {
		case "来源":
			if r == "小星星" {
				return nil, nil
			}
			device.SourceName = r

		case "来源设备ID":
			device.SourceDeviceId = r

		case "经度":
			err := s.parseLng(r, &device)
			if err != nil {
				return nil, err
			}

		case "纬度":
			err := s.parseLat(r, &device)
			if err != nil {
				return nil, err
			}

		case "标题":
			s.parseTitle(r, &device)

		case "设备图":
			s.parseDeviceImage(r, &device)

		case "环境图":
			s.parseEnvironmentImage(r, &device)

		case "联系电话":
			s.parseContact(r, &device)

		case "地址":
			s.parseAddress(r, &device)
		}
	}

	return &device, nil
}

func (s ExcelDeviceImporter) parseTitle(row string, device *entities.BaseDevice) {
	if device.Title == "" || len(row) > len(device.Title) {
		device.Title = row
	}
}

func (s ExcelDeviceImporter) parseLng(row string, device *entities.BaseDevice) error {
	lng, err := strconv.ParseFloat(row, 64)
	if err != nil {
		return err
	}

	device.Longitude = lng

	return nil
}

func (s ExcelDeviceImporter) parseLat(row string, device *entities.BaseDevice) error {
	lat, err := strconv.ParseFloat(row, 64)
	if err != nil {
		return err
	}

	device.Latitude = lat

	return nil
}

func (s ExcelDeviceImporter) parseDeviceImage(row string, device *entities.BaseDevice) {
	var images []string
	err := json.Unmarshal([]byte(row), &images)
	if err != nil {
		split := strings.Split(row, ",")
		for _, s := range split {
			if s != "" {
				device.DeviceImage = s
				break
			}
		}
	}

	if len(images) > 0 {
		device.DeviceImage = images[0]
	}
}

func (s ExcelDeviceImporter) parseEnvironmentImage(row string, device *entities.BaseDevice) {
	var images []string
	err := json.Unmarshal([]byte(row), &images)
	if err != nil {
		device.EnvironmentImage = row
		return
	}

	if len(images) > 0 {
		device.EnvironmentImage = images[0]
		return
	}
}

func (s ExcelDeviceImporter) parseContact(row string, device *entities.BaseDevice) {
	device.Contact = row
}

func (s ExcelDeviceImporter) parseAddress(row string, device *entities.BaseDevice) {
	if device.Address == "" || len(row) > len(device.Title) {
		device.Address = row
	}
}
