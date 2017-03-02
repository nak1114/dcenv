package kii

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
  "encoding/base64"
  "math"

  "github.com/Songmu/prompter"
)

type KiiResLogin struct {
  Id     string `json:"id,omitempty"`           //" : "*************************************",
  Token  string `json:"access_token,omitempty"` //
  Expire int64  `json:"expires_in,omitempty"`   //" : 2147483639,
  TType  string `json:"token_type,omitempty"`   //" : "bearer"
}
type KiiResLoginError struct {
  ECode string `json:"errorCode,omitempty"`         //" : "invalid_grant",
  EDesc string `json:"error_description,omitempty"` //" : "The user was not found or a wrong password was provided",
  Error string `json:"error,omitempty"`             //" : "invalid_grant"
}
type KiiResError struct {
  ErrorCode string `json:"errorCode"` //Error code "QUERY_NOT_SUPPORTED".
  Message   string `json:"message"`   //The error message.
}
type KiiResCreateObj struct {
  Id         string `json:"objectID"`  //   string  The ID of the object.
  CcreatedAt int64  `json:"createdAt"` //   long  The creation date of the object in Unix epoch (milliseconds in UTC).
  DataType   string `json:"dataType"`  //   string  The data type of the object.
}
type KiiObjBase struct {
  Id      string `json:"_id,omitempty"`
  Owner   string `json:"_owner,omitempty"`
  Created int64  `json:"_created,omitempty"`
  Updated int64  `json:"_modified,omitempty"`
  Version string `json:"_version,omitempty"`
}

type KiiResQuery struct {
  QueryDescription  string
  Results           []interface{}
  NextPaginationKey string
}

type User struct {
  Id      string
  Name    string
  Token   string
  TType   string
  Expire  int64
  Ap      *App
  Buckets string
}
type App struct {
  User
  Site string
}
type Query struct {
  R   []interface{}
  L   int
  U   *User
  K   string
  Q   string
  EoQ bool
}

const radix = 26
const BestEffortLimit = "200"

var exit = os.Exit
var StrLoginError = "\nIf you don't sign up for the dcenv resistry yet. Please sign up on [ https://nak1114.github.io/dcenv/sign_up.html ]."

func NewApp(appID string, appKey string, site string, buckets string) (a App) {
  a.Id = appID
  a.Name = appKey
  a.Token = base64.StdEncoding.EncodeToString([]byte(appID + ":" + appKey))
  a.TType = "Basic"
  a.Expire = math.MaxInt64
  a.Ap = &a
  a.Site = site
  a.Buckets = buckets
  return
}

func (u User) AuthHeader() (s string) {
  return "Authorization: " + u.TType + " " + u.Token
}
func (u *User) CB(backets string) *User {
  u.Buckets = backets
  return u
}
func (u User) ExistObj(obj string, res interface{}) (e error) {
  h := []string{
    u.AuthHeader(),
  }
  st, str := curl("GET", "https://"+u.Ap.Site+"/api/apps/"+u.Ap.Id+"/buckets/"+u.Buckets+"/objects/"+obj, h, "")
  if st != 200 {
    v := KiiResError{}
    if e = json.Unmarshal([]byte(str), &v); e != nil {
      return
    }
    return fmt.Errorf("%s : %s", v.ErrorCode, v.Message)
  }
  if e = json.Unmarshal([]byte(str), &res); e != nil {
    return
  }
  return nil
}

func (u *User) DeleteObj(id string) (e error) {
  h := []string{
    u.AuthHeader(),
    "Content-Type: application/json",
  }
  st, str := curl("DELETE", "https://"+u.Ap.Site+"/api/apps/"+u.Ap.Id+"/buckets/"+u.Buckets+"/objects/"+id, h, "")
  if st != 204 {
    v := KiiResError{}
    if e = json.Unmarshal([]byte(str), &v); e != nil {
      return
    }
    e = fmt.Errorf("%s : %s", v.ErrorCode, v.Message)
    return
  }
  return
}

/*
curl -v -X PUT \
  -H "Authorization: Bearer *******************************************" \
  -H "Content-Type: application/json" \
  "https://api-jp.kii.com/api/apps/efo16zkedmd7/buckets/dcenv/objects/defcf5d0-fbd3-11e6-865e-22000b07265b" \
  -d '{"score": 100, "name": "game3"}'
*/
func (u *User) UpdateObj(id string, obj string) (res KiiResCreateObj, e error) {
  h := []string{
    u.AuthHeader(),
    "Content-Type: application/json",
  }
  st, str := curl("PUT", "https://"+u.Ap.Site+"/api/apps/"+u.Ap.Id+"/buckets/"+u.Buckets+"/objects/"+id, h, obj)
  if st != 200 {
    v := KiiResError{}
    if e = json.Unmarshal([]byte(str), &v); e != nil {
      return
    }
    e = fmt.Errorf("%s : %s", v.ErrorCode, v.Message)
    return
  }
  if e = json.Unmarshal([]byte(str), &res); e != nil {
    return
  }
  return
}
func (u *User) CreateObj(obj string) (res KiiResCreateObj, e error) {
  /*
     curl -v -X POST \
     -H "Authorization: Bearer *******************************************" \
     -H "Content-Type: application/json" \
     "https://api-jp.kii.com/api/apps/efo16zkedmd7/buckets/dcenv/objects" \
     -d '{"score": 1800, "name": "game1"}'
  */
  h := []string{
    u.AuthHeader(),
    "Content-Type: application/json",
  }
  st, str := curl("POST", "https://"+u.Ap.Site+"/api/apps/"+u.Ap.Id+"/buckets/"+u.Buckets+"/objects", h, obj)
  if st != 201 {
    v := KiiResError{}
    if e = json.Unmarshal([]byte(str), &v); e != nil {
      return
    }
    e = fmt.Errorf("%s : %s", v.ErrorCode, v.Message)
    return
  }

  if e = json.Unmarshal([]byte(str), &res); e != nil {
    return
  }
  return
}

/*
curl -v -X DELETE \
  -H "Authorization: Bearer *******************************************" \
  "https://api-jp.kii.com/api/apps/efo16zkedmd7/buckets/dcenv/objects/e130c620-fc1d-11e6-865e-22000b07265b/acl/WRITE_EXISTING_OBJECT/UserID:ANY_AUTHENTICATED_USER"
*/
func (u *User) ACLObj(objId string) (e error) {
  h := []string{
    u.AuthHeader(),
  }
  st, str := curl("DELETE", "https://"+u.Ap.Site+"/api/apps/"+u.Ap.Id+"/buckets/"+u.Buckets+"/objects/"+objId+"/acl/WRITE_EXISTING_OBJECT/UserID:ANY_AUTHENTICATED_USER", h, "")
  if st == 204 {
    return
  }
  v := KiiResError{}
  if e = json.Unmarshal([]byte(str), &v); e != nil {
    return
  }
  e = fmt.Errorf("%s : %s", v.ErrorCode, v.Message)
  return
}

func (u *User) WriteToFile(fname string) error {
  buf := fmt.Sprintf("%s\t%s\t%s\t%s", u.Name, u.Id, u.Token, strconv.FormatInt(u.Expire, radix))
  return ioutil.WriteFile(fname, []byte(buf), 0600)
}

func (a App) Login(user string, pass string) (u User, e error) {
  if user == "" {
    user = prompter.Prompt("Enter your Login name", "")
  }
  if pass == "" {
    pass = prompter.Password(fmt.Sprintf("Enter your[%s] login password", user))
    if len(pass) < 4 {
      e = fmt.Errorf("Password Error")
      return
    }
  }

  h := []string{
    a.AuthHeader(),
    "Content-Type: application/json",
  }
  body := `{"grant_type": "password","username": "%s","password": "%s"}`
  st, str := curl("POST", "https://"+a.Site+"/api/apps/"+a.Id+"/oauth2/token", h, fmt.Sprintf(body, user, pass))
  if st != 200 {
    res := KiiResLoginError{}
    if e = json.Unmarshal([]byte(str), &res); e != nil {
      return
    }
    e = fmt.Errorf(res.EDesc + StrLoginError)
    return
  }
  res := KiiResLogin{}
  if e = json.Unmarshal([]byte(str), &res); e != nil {
    return
  }

  exp := time.Now().Add(time.Duration(res.Expire-30) * time.Second).Unix()
  u.Id = res.Id //strings.Replace(res.Id, "-", "", -1)
  u.Name = user
  u.Token = res.Token
  u.TType = "Bearer"
  u.Expire = exp
  u.Ap = &a
  u.Buckets = a.Buckets
  return
}

func (a App) Relogin(fname string) (User, error) {
  if _, err := os.Stat(fname); err != nil {
    return a.Login("", "")
  }
  b, err := ioutil.ReadFile(fname)
  if err != nil {
    return User{}, err
  }

  str := strings.SplitN(string(b), "\t", 4)

  if len(str) < 4 {
    if len(str) < 1 {
      return a.Login("", "")
    } else {
      return a.Login(str[0], "")
    }
  }
  exp, _ := strconv.ParseInt(str[3], radix, 64)
  if time.Now().Unix() > exp {
    fmt.Println("login timeout")
    return a.Login(str[0], "")
  }
  us := User{
    Id:      str[1],
    Name:    str[0],
    Token:   str[2],
    TType:   "Bearer",
    Expire:  exp,
    Ap:      &a,
    Buckets: a.Buckets,
  }

  return us, nil
}

func Logout(fname string) {
  os.Remove(fname)
}

func curl(xtype string, url string, header []string, ibody string) (stat int, body string) {

  bd := strings.NewReader(ibody)
  req, err := http.NewRequest(xtype, url, bd)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  for _, h := range header {
    t := strings.SplitN(h, `: `, 2)
    req.Header.Set(t[0], t[1])
  }
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    fmt.Println(err)
    exit(1)
    return
  }
  defer resp.Body.Close()

  ba, _ := ioutil.ReadAll(resp.Body)
  body = string(ba)
  stat = resp.StatusCode
  return
}

func (q *Query) Next(v interface{}) error {

  if bb, err := json.Marshal(q.R[q.L]); err == nil {
    if err := json.Unmarshal([]byte(bb), v); err != nil {
      q.EoQ = true
      return err
    }
  } else {
    q.EoQ = true
    return err
  }
  q.L += 1
  if q.L >= len(q.R) {
    if q.K == "" {
      q.EoQ = true
      q.L = len(q.R) - 1
      return nil
    }
    qq, err := q.U.Query(q.Q + `,"paginationKey": "` + q.K + `"`)
    *q = qq
    return err
  }
  return nil
}

func (u User) Query(query string) (Query, error) {
  h := []string{
    u.AuthHeader(),
    "Content-Type: application/vnd.kii.QueryRequest+json",
  }
  body := `{"bucketQuery": ` + query + `,"bestEffortLimit": ` + BestEffortLimit + `}`

  st, str := curl("POST", "https://"+u.Ap.Site+"/api/apps/"+u.Ap.Id+"/buckets/"+u.Buckets+"/query", h, body)
  if st != 200 {
    vv := KiiResError{}
    if err := json.Unmarshal([]byte(str), &vv); err != nil {
      return Query{}, err
    }
    return Query{}, fmt.Errorf("%d : %s", st, vv.Message)
  }

  vv := KiiResQuery{}
  if err := json.Unmarshal([]byte(str), &vv); err != nil {
    return Query{}, err
  }
  return Query{
    R:   vv.Results,
    L:   0,
    U:   &u,
    K:   vv.NextPaginationKey,
    Q:   query,
    EoQ: len(vv.Results) == 0,
  }, nil
}
