package i18n

import (
	"strings"
	"fmt"
)

// 获取自定义语言内容, langContent: en-US:SynServer|zh-CN:同步服  lang:zh-CN, 返回: 同步服
func GetCustomLangContent(langContent string, lang string) string {
	if lang == "" {
		lang = "zh-CN"
	}
	arr := strings.Split(langContent, "|")
	if len(arr) <= 1 {
		return langContent
	}
	for _, v := range arr {
		idx := strings.Index(v, ":")
		if idx != -1 {
			lanKey := v[0:idx]
			lanVal := v[idx+1:]
			if lanKey == lang {
				return lanVal
			}
		}
	}
	return fmt.Sprintf("undefined:%v", lang)
}
