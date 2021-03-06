package beego

import (
	"fmt"
	"net/url"
	"strings"
)

type FlashData struct {
	Data map[string]string
}

func NewFlash() *FlashData {
	return &FlashData{
		Data: make(map[string]string),
	}
}

func (fd *FlashData) Notice(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["notice"] = msg
	} else {
		fd.Data["notice"] = fmt.Sprintf(msg, args...)
	}
}

func (fd *FlashData) Warning(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["warning"] = msg
	} else {
		fd.Data["warning"] = fmt.Sprintf(msg, args...)
	}
}

func (fd *FlashData) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["error"] = msg
	} else {
		fd.Data["error"] = fmt.Sprintf(msg, args...)
	}
}

func (fd *FlashData) Store(c *Controller) {
	c.Data["flash"] = fd.Data
	var flashValue string
	for key, value := range fd.Data {
		flashValue += "\x00" + key + ":" + value + "\x00"
	}
	c.Ctx.SetCookie("BEEGO_FLASH", url.QueryEscape(flashValue), 0, "/")
}

func ReadFromRequest(c *Controller) *FlashData {
	flash := &FlashData{
		Data: make(map[string]string),
	}
	if cookie, err := c.Ctx.Request.Cookie("BEEGO_FLASH"); err == nil {
		vals := strings.Split(cookie.Value, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, ":")
				if len(kv) == 2 {
					flash.Data[kv[0]] = kv[1]
				}
			}
		}
		//read one time then delete it
		cookie.MaxAge = -1
		c.Ctx.Request.AddCookie(cookie)
	}
	c.Data["flash"] = flash.Data
	return flash
}
