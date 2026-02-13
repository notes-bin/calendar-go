package api

type LaoHuangLiResponse struct {
	ID        string `json:"id"`        // 黄历ID
	Yangli    string `json:"yangli"`    // 阳历日期
	Yinli     string `json:"yinli"`     // 阴历日期
	Wuxing    string `json:"wuxing"`    // 五行
	Chongsha  string `json:"chongsha"`  // 冲煞
	Baiji     string `json:"baiji"`     // 百忌
	Jishen    string `json:"jishen"`    // 吉神宜趋
	Yi        string `json:"yi"`        // 宜
	Xiongshen string `json:"xiongshen"` // 凶神宜忌
	Ji        string `json:"ji"`        // 忌
}
