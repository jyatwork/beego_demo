package controllers

import (
	"apibeego/models"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"sort"
	"strings"

	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

// Operations about User
type UserController struct {
	beego.Controller
}

func GetMd5String(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	a := hex.EncodeToString(md5Ctx.Sum(nil))
	return a
}

type RegisterReq struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// @Title AddUser
// @Description add user
// @Param   body     body    models.User  true        "The email for login"
// @Success 200 {object} models.User
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /addUser [post]
func (c *UserController) AddUser() {
	var v RegisterReq

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		md5pwd := GetMd5String(v.Password)
		var u models.User
		u.Uname = v.Username
		u.Psw = md5pwd
		u.CreateTime = time.Now().UTC()
		u.LastTime = time.Now().UTC()
		if _, err := models.AddUser(&u); err == nil {
			c.Data["json"] = u
		} else {
			fmt.Println(err.Error())
			c.Data["json"] = err.Error()
		}
	} else {
		fmt.Println(err.Error())
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

type LoginReq struct {
	LoginId  int64  `json:"LoginId"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// @Title CreatUser
// @Description create user
// @Param   body     body    models.User  true        "The email for login"
// @Success 200 {object} models.User
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /creatUser [post]
func (c *UserController) CreatUser() {
	var v LoginReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		md5pwd := GetMd5String(v.Password)
		usrOld := &models.User{Id: v.LoginId, Uname: v.Username, Psw: md5pwd}
		user, _ := models.VerifyUser(usrOld)
		if user != nil {
			// get uid

			id := user[0]["uid"]
			c.Data["json"] = map[string]interface{}{"success": 0, "msg": "登录成功", "uname": v.Username, "id": id}
		} else {
			c.Data["json"] = map[string]interface{}{"success": -1, "msg": "账号密码错误"}
		}

	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title Get Health
// @Description Get Health info
// @Success 200 {object} models.User
// @Param   time     query   string false       "category id"
// @Param   sign    query   string false       "brand id"
// @Param   source   query   string  false       "query of search"
// @Failure 400 no enough input
// @Failure 500 get products common error
// @router /health [get]
func (c *UserController) Health() {
	m := make(map[string]interface{})
	m["time"] = c.GetString("time")
	//	fmt.Println(c.GetString("time"))
	m["sign"] = c.GetString("sign")
	//	fmt.Println(c.GetString("sign"))
	m["source"] = c.GetString("source")
	//	fmt.Println(c.GetString("source"))
	if m["time"] == "" || m["sign"] == "" || m["source"] == "" {
		c.Data["json"] = map[string]interface{}{"code": 2, "error": "数据格式错误"}
	} else {

		if SignVerify(m) {
			if LoginVerify() {
				c.Data["json"] = map[string]interface{}{"code": 0, "sign": "UP", "score": 7}
			} else {
				c.Data["json"] = map[string]interface{}{"code": 0, "sign": "DOWN"}
			}

		} else {
			c.Data["json"] = map[string]interface{}{"code": 1, "error": "鉴权错误"}
		}
	}

	c.ServeJSON() //ServeJSON sends a json response with encoding charset
}

func SignVerify(m map[string]interface{}) bool {
	var strlist []string
	var vc string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		value, _ := ToStr(v)

		strlist = append(strlist, k+"="+value)

	}
	sort.Strings(strlist)
	for k, v := range strlist {
		if k == 0 {
			vc = v
		} else {
			vc = vc + "&" + v
		}
	}
	fmt.Println("vc:", vc)
	sign, ok := m["sign"].(string)
	if ok {
		signvc_md5 := GetMd5String(vc)
		fmt.Println(signvc_md5)
		if strings.Compare(strings.ToUpper(signvc_md5), strings.ToUpper(sign)) == 0 {
			return true
		}
	}

	return false
}

func ToStr(i interface{}) (string, error) {
	var (
		str string
	)
	switch v := i.(type) {
	case []interface{}:

		for k, u := range v {
			a, _ := ToStr(u)
			if k == 0 {
				str = "[" + a
			} else if k != len(v)-1 {
				str = str + " " + a
			} else if k == len(v)-1 {
				str = str + " " + a + "]"
			}
		}

		return str, nil
	case []string:
		for k, u := range v {
			a, _ := ToStr(u)
			if k == 0 {
				str = "[" + a
			} else if k != len(v)-1 {
				str = str + " " + a
			} else if k == len(v)-1 {
				str = str + " " + a + "]"
			}
		}
		return str, nil
	case string:
		return v, nil
		//	case interface{}:
		//		str, err = ToStr(i)
		//		if err != nil {
		//			return "", fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		//		}
		//		return str, nil
	case bool:
		return strconv.FormatBool(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case uint:
		return strconv.FormatInt(int64(v), 10), nil
	case uint64:
		return strconv.FormatInt(int64(v), 10), nil
	case uint32:
		return strconv.FormatInt(int64(v), 10), nil
	case uint16:
		return strconv.FormatInt(int64(v), 10), nil
	case uint8:
		return strconv.FormatInt(int64(v), 10), nil
	case []byte:
		return string(v), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return v.String(), nil
	case error:
		return v.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}

func LoginVerify() bool {
	var m map[string]interface{}
	req, err := http.NewRequest("POST", "http://localhost:8088/login", strings.NewReader(`{"Username":"test","Password":"15523"}`))
	if err != nil {
		fmt.Println(err.Error())

	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	res, err := client.Do(req) //Do方法发送请求，返回HTTP回复
	if err != nil {
		//noinspection GoPlaceholderCount
		fmt.Println("PushNotice http client sending an HTTP request error: %s", err.Error())

	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		//noinspection GoPlaceholderCount
		fmt.Println("PushNotice the response's status code from an HTTP request is: %d", res.StatusCode)
	} else {
		d, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("PushNotice ReadAll error:", err.Error())
		}

		if json.Unmarshal(d, &m) == nil {

			if value, ok := m["success"].(float64); ok {

				if value == 0 {
					return true
				}
			}
		}
	}
	return false
}
