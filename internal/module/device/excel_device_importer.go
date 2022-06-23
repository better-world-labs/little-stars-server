package device

import (
	"aed-api-server/internal/interfaces/entities"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"time"
)

type ExcelDeviceImporter struct {
}

func NewExcelDeviceImporter() IDeviceImporter {
	return &ExcelDeviceImporter{}
}

func (s ExcelDeviceImporter) ImportDevices(reader io.Reader) ([]*entities.Device, error) {
	r, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}

	sheets := r.GetSheetList()

	return s.parseSheetsDevices(r, sheets)
}

func (s ExcelDeviceImporter) parseSheetsDevices(excel *excelize.File, sheets []string) ([]*entities.Device, error) {
	var devices []*entities.Device

	for _, sheet := range sheets {
		d, err := s.parseSheetDevices(excel, sheet)
		if err != nil {
			return nil, err
		}

		devices = append(devices, d...)
	}

	return devices, nil
}

func (s ExcelDeviceImporter) parseSheetDevices(excel *excelize.File, sheet string) ([]*entities.Device, error) {
	rows, err := excel.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	return s.parseDevices(rows)
}

func (s ExcelDeviceImporter) parseDevices(rows [][]string) ([]*entities.Device, error) {
	var devices []*entities.Device

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

func (s ExcelDeviceImporter) parseDevice(head []string, row []string) (*entities.Device, error) {
	var device entities.Device

	device.Id = uuid.NewString()
	device.Created = time.Now().UnixMilli()
	device.Source = entities.DeviceSourceImported
	for i, r := range row {
		switch head[i] {
		case "来源":
			if r == "小星星" {
				return nil, nil
			}

		case "标题":
			device.Title = r

		case "经度":
			lng, err := strconv.ParseFloat(r, 64)
			if err != nil {
				return nil, err
			}

			device.Longitude = lng

		case "纬度":
			lat, err := strconv.ParseFloat(r, 64)
			if err != nil {
				return nil, err
			}

			device.Latitude = lat

		case "设备图":
			device.DeviceImage = r

		case "环境图":
			device.EnvironmentImage = r

		case "联系电话":
			device.Tel = r

		case "地址":
			device.Address = r
		}
	}

	return &device, nil
}
