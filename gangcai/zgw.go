package justcy_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	//	"log"
	"regexp"
	"strconv"
	// "strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	Zgwcompany.Register()
}

var Zgwcompany = &Spider{
	Name:        "中国钢材网",
	Description: "中国钢材网钢厂信息 [Auto Page] [http://www.zgw.com]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 300000}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "https://www.zgw.com/Shop/" + strconv.Itoa(loop[0]),
							Rule: "获取结果",
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					ss := query.Find(".content_right")
					ss.Each(func(i int, goq *goquery.Selection) {
						ctx.SetTemp("company_info", goq)
						ctx.Parse("获取结果")

					})
				},
			},

			"获取结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"公司名称",
					"企业类型",
					"主营",
					"手机",
					"传真",
					"业务联系人",
					"业务联系电话",
					"仓库地址",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var selectObj = query.Find(".content_right")

					// 获取内容
					content := selectObj.Find(".leixing").Text()

					// 过滤标签
					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					content = re.ReplaceAllString(content, "")

					//					re, _ = regexp.Compile("[\\n\\t ]")
					//					content = re.ReplaceAllString(content, "")

					companyName := selectObj.Find(".gogs").Text()

					// companytype
					re, _ = regexp.Compile("企业类型：\\n(.*?)\\n")
					companyTypeTemp := re.FindStringSubmatch(content)
					companyType := ""
					if len(companyTypeTemp) > 0 {
						companyType = companyTypeTemp[1]
					}

					//					companyType := selectObj.Find(".leixing").Find(".lei1").Find("p").Eq(0).Text()
					//					companyMajor := selectObj.Find(".leixing").Find(".lei1").Find("p").Eq(1).Text()

					re, _ = regexp.Compile("主营：(.*?)\\n")
					companyMajorTemp := re.FindStringSubmatch(content)
					companyMajor := ""
					if len(companyMajorTemp) > 0 {
						companyMajor = companyMajorTemp[1]
					}

					re, _ = regexp.Compile("机：(.*?)\\n")
					phoneTemp := re.FindStringSubmatch(content)
					phone := ""
					if len(phoneTemp) > 0 {
						phone = phoneTemp[1]
					}

					re, _ = regexp.Compile("传真：(.*?)\\n")
					telTemp := re.FindStringSubmatch(content)
					tel := ""
					if len(telTemp) > 0 {
						tel = telTemp[1]
					}

					re, _ = regexp.Compile("业务联系人：(.*?)\\n")
					bosstemp := re.FindStringSubmatch(content)
					boss := ""
					if len(bosstemp) > 0 {
						boss = bosstemp[1]
					}
					re, _ = regexp.Compile("业务联系人电话：(.*?)\\n")
					boss_phoneTemp := re.FindStringSubmatch(content)
					boss_phone := ""
					if len(boss_phoneTemp) > 0 {
						boss_phone = boss_phoneTemp[1]
					}

					addressText := selectObj.Find(".lei3").Text()
					re, _ = regexp.Compile("仓库地址：(.*?)\\n")
					addressTemp := re.FindStringSubmatch(addressText)
					address := ""
					if len(addressTemp) > 0 {
						address = addressTemp[1]
					}
					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: companyName,
						1: companyType,
						2: companyMajor,
						3: phone,
						4: tel,
						5: boss,
						6: boss_phone,
						7: address,
					})
				},
			},
		},
	},
}
