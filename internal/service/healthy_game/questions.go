package healthy_game

import "aed-api-server/internal/interfaces/entities"

var questions = []*entities.Question{
	{
		Id:       1,
		OriginNo: "1",
		Desc:     "是否确诊冠心病？",
		Type:     "疾病背景",
		SubType:  "冠心病",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       2,
		OriginNo: "2",
		Desc:     "是否确诊以下任一疾病：肥厚型心肌病、限制型心肌病、致心律失常性右心室心肌病？",
		Type:     "疾病背景",
		SubType:  "冠心病",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       3,
		OriginNo: "3",
		Desc:     "直系亲属（含父母、直系兄弟姐妹、子女）是否确诊以下任一疾病：冠心病、肥厚型心肌病、限制型心肌病、致心律失常性右心室心肌病？",
		Type:     "疾病背景",
		SubType:  "冠心病",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "不清楚",
				Score: 0,
			},
		},
	},

	//4
	{
		Id:       4,
		OriginNo: "4-1",
		Desc:     "是否确诊高胆固醇？",
		Type:     "疾病背景",
		SubType:  "三高筛查",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 300,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "未体检,不知情",
				Score: 0,
			},
		},
	},
	{
		Id:       5,
		OriginNo: "4-2",
		Desc:     "是否喜欢大量食用高油脂食物，比如：油炸食物、肥肉（猪五花、羔羊肉等）、动物内脏、脑花、黄油食品、蛋糕、巧克力等？",
		Type:     "疾病背景",
		SubType:  "三高筛查",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       6,
		OriginNo: "4-3",
		Desc:     "直系亲属（含父母、直系兄弟姐妹、子女）是否确诊高胆固醇？",
		Type:     "疾病背景",
		SubType:  "三高筛查",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "不清楚",
				Score: 0,
			},
		},
	},

	//5
	{
		Id:       7,
		OriginNo: "5-1",
		Desc:     "是否确诊高血压？",
		Type:     "疾病背景",
		SubType:  "三高筛查",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 300,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "未体检,不知情",
				Score: 0,
			},
		},
	},
	{
		Id:       8,
		OriginNo: "5-2",
		Desc:     "是否喜欢大量食用高盐食物，比如：酸菜咸菜、熏肉腊肉、咸口酱料、重口味家常菜（水煮牛肉、水煮鱼等）、火锅、薯片等？",
		Type:     "疾病背景",
		SubType:  "三高筛查",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       9,
		OriginNo: "5-3",
		Desc:     "直系亲属（含父母、直系兄弟姐妹、子女）是否确诊高血压？",
		Type:     "疾病背景",
		SubType:  "三高筛查",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "不清楚",
				Score: 0,
			},
		},
	},

	//6
	{
		Id:       10,
		OriginNo: "6-1",
		Desc:     "以下哪个选项说明现在的吸烟（含电子烟）情况？",
		Type:     "生活习惯",
		SubType:  "吸烟",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "不吸烟",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "平均＜8支/天",
				Score: 50,
			},
			{
				Index: "C",
				Desc:  "平均9-14支/天",
				Score: 75,
			},
			{
				Index: "D",
				Desc:  "平均≥15支/天",
				Score: 100,
			},
		},
	},
	{
		Id:       11,
		OriginNo: "6-2",
		Desc:     "以下哪个选项说明戒烟情况？",
		Type:     "生活习惯",
		SubType:  "吸烟",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "从未吸烟或戒烟＞5年",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "戒烟＞1个月",
				Score: 75,
			},
			{
				Index: "C",
				Desc:  "戒烟＞1年",
				Score: 50,
			},
			{
				Index: "D",
				Desc:  "未戒烟",
				Score: 0,
			},
		},
	},

	//7
	{
		Id:       12,
		OriginNo: "7-1",
		Desc:     "是否确诊糖尿病？",
		Type:     "生活习惯",
		SubType:  "饮食",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "未体检,不知情",
				Score: 0,
			},
		},
	},
	{
		Id:       13,
		OriginNo: "7-2",
		Desc:     "是否每天高能量饮食（主观感受，每顿饭都吃得比较撑）？",
		Type:     "生活习惯",
		SubType:  "饮食",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 50,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       14,
		OriginNo: "7-3",
		Desc:     "是否喜欢大量食用高糖食物，比如白米饭、面条、包子馒头、含糖饮料、奶茶、蜜饯零食、冰淇淋、蛋糕点心？",
		Type:     "生活习惯",
		SubType:  "饮食",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 50,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       15,
		OriginNo: "7-4",
		Desc:     "直系亲属（含父母、直系兄弟姐妹、子女）是否确诊糖尿病？",
		Type:     "生活习惯",
		SubType:  "饮食",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 50,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "不清楚",
				Score: 0,
			},
		},
	},

	//8
	{
		Id:       16,
		OriginNo: "8-1",
		Desc:     "身体质量指数BMI（BMI=体重（kg）/身高（m）的平方）为？",
		Type:     "身体情况",
		SubType:  "肥胖",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "偏瘦，BMI＜18.5",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "18.5≤BMI＜24",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "超重，24≤BMI＜28",
				Score: 75,
			},
			{
				Index: "D",
				Desc:  "肥胖，BMI≥28",
				Score: 100,
			},
		},
	},

	//9
	{
		Id:       17,
		OriginNo: "9-1",
		Desc:     "以下哪个选项更符合你的每周运动情况?（包含跑步、打球、游泳、撸铁等正式运动，也包括慢走、骑车、登山等休闲运动）",
		Type:     "生活习惯",
		SubType:  "运动",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "几乎不运动",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "运动1-2次（累计超过30分钟）",
				Score: 75,
			},
			{
				Index: "C",
				Desc:  "运动2-3次（累计超过90分钟）",
				Score: 50,
			},
			{
				Index: "D",
				Desc:  "运动3次以上（累计超过150分钟）",
				Score: 0,
			},
		},
	},

	//10
	{
		Id:       18,
		OriginNo: "10-1",
		Desc:     "以下哪个选项更符合你的平均睡眠时间？",
		Type:     "生活习惯",
		SubType:  "睡眠",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "平均睡眠＜5小时/天",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "平均睡眠5-7小时/天",
				Score: 50,
			},
			{
				Index: "C",
				Desc:  "平均睡眠7-9小时/天",
				Score: 0,
			},
			{
				Index: "D",
				Desc:  "平均睡眠≥9小时/天",
				Score: 0,
			},
		},
	},
	{
		Id:       19,
		OriginNo: "10-2",
		Desc:     "以下哪个选项更符合你的睡眠质量？",
		Type:     "生活习惯",
		SubType:  "睡眠",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "睡得很好，起床后神清气爽",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "睡得一般，起床略微艰难",
				Score: 25,
			},
			{
				Index: "C",
				Desc:  "睡得很浅，起床后仍觉得疲惫",
				Score: 50,
			},
			{
				Index: "D",
				Desc:  "经常失眠",
				Score: 100,
			},
		},
	},

	//11
	{
		Id:       20,
		OriginNo: "11-1",
		Desc:     "以下哪个情况更能反映你的性格？",
		Type:     "性格",
		SubType:  "",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "争强好胜，做事时间紧迫感很强。",
				Score: 100,
			},
			{
				Index: "B",
				Desc:  "慢条斯理，做事平和不争不抢。",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "以上都不是。",
				Score: 0,
			},
		},
	},

	//12
	{
		Id:       21,
		OriginNo: "12-1",
		Desc:     "请选择性别？",
		Type:     "性别与年龄",
		SubType:  "",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "女性",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "男性",
				Score: 50,
			},
		},
	},
	{
		Id:       22,
		OriginNo: "12-2",
		Desc:     "请选择年龄?",
		Type:     "性别与年龄",
		SubType:  "",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "＜35岁",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "35-50岁",
				Score: 50,
			},
			{
				Index: "C",
				Desc:  "＞50岁",
				Score: 100,
			},
		},
	},

	//13
	{
		Id:       23,
		OriginNo: "13-1",
		Desc:     "是否了解猝死的大概发病原因及机制？",
		Type:     "健康知识及急救技能",
		SubType:  "健康素养",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       24,
		OriginNo: "13-2",
		Desc:     "是否了解猝死急救方法CPR？",
		Type:     "健康知识及急救技能",
		SubType:  "健康素养",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       25,
		OriginNo: "13-3",
		Desc:     "是否了解猝死黄金急救仪器AED？",
		Type:     "健康知识及急救技能",
		SubType:  "健康素养",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       26,
		OriginNo: "13-4",
		Desc:     "是否了解自己生活、工作周边的AED设备位置？",
		Type:     "周边AED配置率",
		SubType:  "",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "是",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "不是",
				Score: 0,
			},
		},
	},
	{
		Id:       27,
		OriginNo: "13-5",
		Desc:     "日常生活半径的AED风险情况？",
		Type:     "周边AED配置率",
		SubType:  "",
		Options: []*entities.AnswerOption{
			{
				Index: "A",
				Desc:  "高风险",
				Score: 0,
			},
			{
				Index: "B",
				Desc:  "中风险",
				Score: 0,
			},
			{
				Index: "C",
				Desc:  "低风险",
				Score: 0,
			},
		},
	},
}
