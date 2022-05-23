package feedback

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
)

func Test_Excel(t *testing.T) {
	file := excelize.NewFile()
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		fmt.Println(err)
	}
	styleID, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "#777777"}})
	if err != nil {
		fmt.Println(err)
	}
	if err := streamWriter.SetRow("A1", []interface{}{
		excelize.Cell{StyleID: styleID, Value: "Data"}}); err != nil {
		fmt.Println(err)
	}
	//for rowID := 2; rowID <= 102400; rowID++ {
	//	row := make([]interface{}, 50)
	//	for colID := 0; colID < 50; colID++ {
	//		row[colID] = rand.Intn(640000)
	//	}
	//	cell, _ := excelize.CoordinatesToCellName(1, rowID)
	//	if err := streamWriter.SetRow(cell, row); err != nil {
	//		fmt.Println(err)
	//	}
	//}
	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
	}
	if err := file.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
