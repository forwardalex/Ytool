package enum

import (
	"strings"
	"sync"
)

type CronTaskEnum struct {
	Week   Enum
	Day    Enum
	Hour   Enum
	Minute Enum
	Second Enum
}

var inss *CronTaskEnum
var onces sync.Once
var itemList []*Enum

func newCronTaskEnum() *CronTaskEnum {
	return &CronTaskEnum{
		Week:   Enum{Id: 1, Desc: "Week"},   // 周
		Day:    Enum{Id: 2, Desc: "Day"},    // 天
		Hour:   Enum{Id: 3, Desc: "Hour"},   // 小时
		Minute: Enum{Id: 4, Desc: "Minute"}, // 分钟
		Second: Enum{Id: 5, Desc: "Second"}, // 秒
	}
}

func GetCronTaskEnum() *CronTaskEnum {

	if inss == nil {
		onces.Do(func() {
			if inss == nil {
				inss = newCronTaskEnum()
			}
		})
	}

	return inss
}

func (*CronTaskEnum) Convert(desc string) *Enum {

	if len(itemList) == 0 {
		itemList = append(itemList, &GetCronTaskEnum().Week)
		itemList = append(itemList, &GetCronTaskEnum().Day)
		itemList = append(itemList, &GetCronTaskEnum().Hour)
		itemList = append(itemList, &GetCronTaskEnum().Minute)
		itemList = append(itemList, &GetCronTaskEnum().Second)
	}

	for _, item := range itemList {
		if strings.EqualFold(item.Desc, desc) {
			return item
		}
	}

	return nil
}
