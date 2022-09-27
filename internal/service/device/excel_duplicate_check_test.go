package device

import (
	"aed-api-server/internal/interfaces/entities"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"os"
	"testing"
)

var importer = NewExcelDeviceImporter()

func TestDuplicateCheck(t *testing.T) {
	toBeImported, err := loadDevices("/home/shenweijie/下载/成都aed总和(去重)1.xlsx")
	require.Nil(t, err)

	localDevices, err := loadDevices("/home/shenweijie/db_star_device.xlsx")
	require.Nil(t, err)

	var duplicate int
	var invalidData int
	var newDevices []*entities.BaseDevice

	fmt.Println("----------------------------print duplicate device----------------------------------")
	for _, toBeImportedDevice := range toBeImported {
		if toBeImportedDevice.Address == "" && toBeImportedDevice.Title == "" {
			invalidData++
			continue
		}

		if duplicatedDevice := findDuplicateFromLocal(toBeImportedDevice, localDevices); duplicatedDevice != nil {
			duplicate++
			fmt.Printf("DB %s\n", devicePrintString(duplicatedDevice))
			fmt.Printf("File: %s\n", devicePrintString(toBeImportedDevice))
			continue
		}

		newDevices = append(newDevices, toBeImportedDevice)
	}

	fmt.Println("------------------------------print new device------------------------------")
	for _, d := range newDevices {
		fmt.Println(devicePrintString(d))
	}

	fmt.Printf("finished. %d device duplicate, %d new device, %d invalid data", duplicate, len(newDevices), invalidData)
	f := excelize.NewFile()
	require.Nil(t, err)
	index := f.NewSheet("新设备")
	f.SetActiveSheet(index)
	err = f.SetCellValue("新设备", "A1", "ID")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "B1", "地址")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "C1", "标题")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "D1", "经度")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "E1", "纬度")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "F1", "联系电话")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "G1", "设备图")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "H1", "来源")
	require.Nil(t, err)
	err = f.SetCellValue("新设备", "I1", "来源设备ID")
	require.Nil(t, err)

	for i, d := range newDevices {
		err := f.SetCellValue("新设备", fmt.Sprintf("B%d", i+2), d.Address)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", "B1", "地址")
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("C%d", i+2), d.Title)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("D%d", i+2), d.Longitude)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("E%d", i+2), d.Latitude)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("F%d", i+2), d.Latitude)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("G%d", i+2), d.DeviceImage)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("H%d", i+2), d.SourceName)
		require.Nil(t, err)
		err = f.SetCellValue("新设备", fmt.Sprintf("I%d", i+2), d.SourceDeviceId)
		require.Nil(t, err)
	}
	err = f.SaveAs("/home/shenweijie/out.xlsx")
	require.Nil(t, err)
}

func devicePrintString(device *entities.BaseDevice) string {
	return fmt.Sprintf("id=%s address=%s title=%s lng=%f lat=%f deviceImage=%s, sourceName=%s, sourceDeviceId=%s", device.Id, device.Address, device.Title, device.Longitude, device.Latitude, device.DeviceImage, device.SourceName, device.SourceDeviceId)
}

func findDuplicateFromLocal(toBeImported *entities.BaseDevice, localDevices []*entities.BaseDevice) *entities.BaseDevice {
	for _, localDevice := range localDevices {
		if isDuplicate(toBeImported, localDevice) {
			return localDevice
		}
	}

	return nil
}

func isDuplicate(d1, d2 *entities.BaseDevice) bool {
	return d1.DistanceOf(d2.Coordinate) <= 10
}

func loadDevices(path string) ([]*entities.BaseDevice, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	devices, err := importer.ImportDevices(f)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
