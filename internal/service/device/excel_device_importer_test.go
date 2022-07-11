package device

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestImport(t *testing.T) {
	importer := NewExcelDeviceImporter()
	f, err := os.Open("/home/shenweijie/新增爬虫设备.xlsx")
	require.Nil(t, err)

	devices, err := importer.ImportDevices(f)
	require.Nil(t, err)
	fmt.Println(devices)
}
