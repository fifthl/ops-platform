package utilModel

type IdAndSecret struct {
	ID        string
	Secret    string
	AccountId string
	Name      string
}

var Key = map[string]IdAndSecret{
	"1": {
		ID:        "LBSnKHUvFYIwzHsce18BMqsbKyzLPZ5xiB18kFgSAnc=",
		Secret:    "GMJk9z5r881ZKPvtvzCR3rE2AXK2z/ug5B/asgpK7wA=",
		AccountId: "1",
		Name:      "zb", //总部
	},
	"2": {
		ID:        "MI58WUTio0ECCioErhQ/p3F757kbAZfc0QMURF/otKg=",
		Secret:    "3bPK5PrwhU6Z+fjhOkF75Hrc4yekPR2KNRA6lR6GIbo=",
		AccountId: "2",
		Name:      "gj", //管家
	},
	"3": {
		ID:        "taQW+3qRXs61Bno9k15y4DlUVrJoq4Axr2p5ngjDmk0=",
		Secret:    "lSZWEyOG5o1xUZk+SGXpm/SgQPiyZnComWPAxNZnLNc=",
		AccountId: "3",
		Name:      "st", //生泰
	},
	"4": {
		ID:        "MgSYFn5nC3qD027gTRCqEKV82yFkzdfRzd7aHZ1l7B4=",
		Secret:    "1DY0wKV1VYKt+tNYGCpLP40aT+Dpa7AY9mxgAV2f25k=",
		AccountId: "4",
		Name:      "xx", //信息
	},
	"5": {
		ID:        "jYvDT9Su0MX1LnGi6kLMTTrRGREUgDssSvJOu7AzW34=",
		Secret:    "PZ+G0ua+b6R0qXRo1xLengriUneiT6r6Dn1sAD4uO20=",
		AccountId: "5",
		Name:      "sc", // 市场
	},
	"6": {
		ID:        "iH2t6+Ep+QH91GLFQmKmZleDP9y5hBMzoZvyH8zCSTE=",
		Secret:    "CoYhbLgv8Q9ccQOQjAIxpfgeO0XwjX+q/3Mfk/mqxyA=",
		AccountId: "6",
		Name:      "stOld", // 旧机电
	},
	"7": {
		ID:        "SIGIX6VILU3F+9YutT2h19IpeddamRP80Y/fFYhhcG8=",
		Secret:    "bV9WFRyJ21dES3EPXZn9uNM7Ki84p9GW7j7B3KIveII=",
		AccountId: "7",
		Name:      "gjOld", // 旧管家
	},
}
