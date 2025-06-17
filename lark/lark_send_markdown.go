/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package lark
 *@file    lark_test
 *@date    2025/4/10 15:39
 */

package lark

import (
    "StarRocksQuerys/util"
    "fmt"
    "github.com/go-resty/resty/v2"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

func Send2Markdown(title, message string, roboots []string) {
    body := strings.NewReplacer("\n", "<br>").Replace(message)
    msg := fmt.Sprintf(`
{
  "msg_type": "interactive",
  "card": {
    "schema": "2.0",
    "config": {
      "update_multi": true,
      "style": {
        "text_size": {
          "normal_v2": {
            "default": "normal",
            "pc": "normal",
            "mobile": "heading"
          }
        }
      }
    },
    "body": {
      "direction": "vertical",
      "padding": "12px 12px 12px 12px",
      "elements": [
        {
          "tag": "markdown",
          "content": "%v",
          "text_align": "left",
          "text_size": "normal_v2",
          "margin": "0px 0px 0px 0px"
        }
      ]
    },
    "header": {
      "title": {
        "tag": "plain_text",
        "content": "%s"
      },
      "subtitle": {
        "tag": "plain_text",
        "content": "%s"
      },
      "template": "blue",
      "padding": "12px 12px 12px 12px"
    }
  }
}`, body, filepath.Base(os.Args[0]), title)

    var wg sync.WaitGroup
    for _, roboot := range roboots {
        wg.Add(1)
        go func(roboot string) {
            defer wg.Done()
            r := restys(fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", roboot), msg)
            util.Loggrs.Info(string(r))
        }(roboot)
    }
    wg.Wait()
}

func restys(uri, body string) []byte {
    //发送POST请求并处理响应
    respones, err := resty.New().SetProxy(larkproxy).R().
        SetHeader("Content-Type", "application/json;charset=utf-8").
        SetBody(body).
        Post(uri)
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(string(respones.Body()))
    return respones.Body()
}
