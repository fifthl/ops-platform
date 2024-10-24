package utilModel

type IdAndSecret struct {
	ID        string
	Secret    string
	AccountId string
	Name      string
}

var Key = map[string]IdAndSecret{
	"1": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "GMJk9z5r881ZKPvtvzCR3rE2AXK2z/ug5B/asgpK7wA=",
		AccountId: "1",
		Name:      "zb", //总部
	},
	"2": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "3bPK5PrwhU6Z+fjhOkF75Hrc4yekPR2KNRA6lR6GIbo=",
		AccountId: "2",
		Name:      "gj", //管家
	},
	"3": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "lSZWEyOG5o1xUZk+SGXpm/SgQPiyZnComWPAxNZnLNc=",
		AccountId: "3",
		Name:      "st", //生泰
	},
	"4": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "1DY0wKV1VYKt+tNYGCpLP40aT+Dpa7AY9mxgAV2f25k=",
		AccountId: "4",
		Name:      "xx", //信息
	},
	"5": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "PZ+G0ua+b6R0qXRo1xLengriUneiT6r6Dn1sAD4uO20=",
		AccountId: "5",
		Name:      "sc", // 市场
	},
	"6": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "CoYhbLgv8Q9ccQOQjAIxpfgeO0XwjX+q/3Mfk/mqxyA=",
		AccountId: "6",
		Name:      "stOld", // 旧机电
	},
	"7": {
		ID:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		Secret:    "bV9WFRyJ21dES3EPXZn9uNM7Ki84p9GW7j7B3KIveII=",
		AccountId: "7",
		Name:      "gjOld", // 旧管家
	},
}
