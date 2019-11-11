package handler

import (
	"2019_2_IBAT/internal/pkg/auth"
	csrf "2019_2_IBAT/internal/pkg/csrf"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"log"
	"time"

	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateEmployer(w http.ResponseWriter, r *http.Request) { //+
	// defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("CreateEmployer Start")

	uuid, err := h.UserService.CreateEmployer(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: err.Error()})
		w.Write([]byte(errJSON))
		return
	}

	log.Println("CreateEmployer Employer was created")

	authInfo, cookieValue, err := h.AuthService.CreateSession(uuid, EmployerStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	token, err := csrf.Tokens.Create(authInfo.ID.String(), cookieValue, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Create token failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	expiresAt, err := time.Parse(TimeFormat, authInfo.Expires)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Handle CreateSession:  Time parsing failed")
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	cookie := http.Cookie{
		Name:    auth.CookieName,
		Value:   cookieValue,
		Expires: expiresAt,
	}

	w.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
	w.Header().Set("X-Csrf-Token", token)
	http.SetCookie(w, &cookie)
	RoleJSON, _ := json.Marshal(Role{Role: authInfo.Role})

	w.Write([]byte(RoleJSON))
}

func (h *Handler) GetEmployerById(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id_string := mux.Vars(r)["id"]
	log.Printf("GetEmployerById id_string: %s\n", id_string)
	emplId, err := uuid.Parse(id_string)
	log.Println("GetEmployerById Handler Start")

	if err != nil {
		log.Printf("GetEmployerById Parse id error: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	employer, err := h.UserService.GetEmployer(emplId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errJSON, _ := json.Marshal(Error{Message: InvalidIdMsg})
		w.Write([]byte(errJSON))
		return
	}

	employer.Password = "" //danger
	employerJSON, _ := json.Marshal(employer)

	w.Write([]byte(employerJSON))
}

func (h *Handler) GetEmployers(w http.ResponseWriter, r *http.Request) { //+
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	params := h.ParseEmplQuery(r.URL.Query())
	employers, err := h.UserService.GetEmployers(params)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errJSON, _ := json.Marshal(Error{Message: InternalErrorMsg})
		w.Write([]byte(errJSON))
		return
	}

	employerJSON, _ := json.Marshal(employers)

	w.Write([]byte(employerJSON))
}

func (h *Handler) ParseEmplQuery(query url.Values) map[string]interface{} {
	params := make(map[string]interface{})

	if query.Get("company_name") != "" {
		params["company_name"] = query.Get("company_name")
	} else {
		if query.Get("empl_num") != "" {
			params["empl_num"] = query.Get("empl_num")
		}
		if query.Get("region") != "" {
			params["region"] = query.Get("region")
		}
	}

	return params
}

// write tcp 127.0.0.1:37786->127.0.0.1:6379: use of closed network connection

// Can not get auth info: redigo: unexpected response line
// (possible server error or unsupported concurrent read by application)

// fail
// &{GET /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f HTTP/2.0 2 0
// 	map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Accept-Language:
// 	[ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7] Cookie:[session-id=uoa7If53qWDgQdFR3DrGSeZXL9nLnBZT]
// 	Origin:[http://localhost:8080] Referer:[http://localhost:8080/]
// 	User-Agent:[Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36
// 	(KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36]]
// 	0xc0000751a0 <nil> 0 [] false 82.146.43.113:8080 map[]
// 	map[] <nil> map[] 93.171.198.4:37498 /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f
// 	0xc000268c60 <nil> <nil> 0xc0000260d0}

//correct
// &{GET /employer/b7c44bf4-cd2c-491b-985a-3704965b567e HTTP/2.0 2 0
// 	map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Accept-Language:
// 	[ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7] Cookie:[session-id=uoa7If53qWDgQdFR3DrGSeZXL9nLnBZT]
// 	Origin:[http://localhost:8080] Referer:[http://localhost:8080/]
// 	User-Agent:[Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 z
// 	(KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36]]
// 	0xc0002ce7b0 <nil> 0 [] false 82.146.43.113:8080 map[]
// 	 map[] <nil> map[] 93.171.198.4:37498 /employer/b7c44bf4-cd2c-491b-985a-3704965b567e
// 	 0xc000268c60 <nil> <nil> 0xc0002ce930}

// 2019/11/10 07:52:08 AuthMiddleware: passing to serve

//no error
// 2019/11/10 07:52:22 LOG START [GET] 93.171.198.4:37498, /employer/b7c44bf4-cd2c-491b-985a-3704965b567e
// 2019/11/10 07:52:22 AuthMiddleware: started
// 2019/11/10 07:52:22 &{0xc000086dc0}
// 2019/11/10 07:52:22 AuthStorage: Get started
// 2019/11/10 07:52:22 LOG START [GET] 93.171.198.4:37498, /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f
// 2019/11/10 07:52:22 AuthMiddleware: started
// 2019/11/10 07:52:22 &{0xc000086dc0}
// 2019/11/10 07:52:22 AuthStorage: Get started
// 2019/11/10 07:52:22 AuthMiddleware: auth_record was setted
// 2019/11/10 07:52:22 CTX
// 2019/11/10 07:52:22 &{GET /employer/b7c44bf4-cd2c-491b-985a-3704965b567e HTTP/2.0 2 0 map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Accept-Language:[ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7] Cookie:[session-id=uoa7If53qWDgQdFR3DrGSeZXL9nLnBZT] Origin:[http://localhost:8080] Referer:[http://localhost:8080/] User-Agent:[Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36]] 0xc0002cfa40 <nil> 0 [] false 82.146.43.113:8080 map[] map[] <nil> map[] 93.171.198.4:37498 /employer/b7c44bf4-cd2c-491b-985a-3704965b567e 0xc000268c60 <nil> <nil> 0xc0002cfda0}
// 2019/11/10 07:52:22 AuthMiddleware: passing to serve
// 2019/11/10 07:52:22 req.RequestURI = /employer/b7c44bf4-cd2c-491b-985a-3704965b567e
// 2019/11/10 07:52:22 req.Method = GET
// 2019/11/10 07:52:22 GetEmployerById id_string: %s b7c44bf4-cd2c-491b-985a-3704965b567e
// 2019/11/10 07:52:22 GetEmployerById Handler Start
// 2019/11/10 07:52:22 GetEmployer Service Start
// 2019/11/10 07:52:22 GetEmployer Repository Start
// 2019/11/10 07:52:22 LOG END [GET] 93.171.198.4:37498, /employer/b7c44bf4-cd2c-491b-985a-3704965b567e 2.840019ms

//error
// 2019/11/10 07:52:22 LOG START [GET] 93.171.198.4:37498, /employer/b7c44bf4-cd2c-491b-985a-3704965b567e
// 2019/11/10 07:52:22 AuthMiddleware: started
// 2019/11/10 07:52:22 &{0xc000086dc0}
// 2019/11/10 07:52:22 AuthStorage: Get started
// 2019/11/10 07:52:22 AuthStorage: Can not get auth info: redigo: unexpected response line (possible server error or unsupported concurrent read by application)
// 2019/11/10 07:52:22 AuthMiddleware: failed to set auth_record
// 2019/11/10 07:52:22 CTX
// 2019/11/10 07:52:22 &{GET /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f HTTP/2.0 2 0 map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Accept-Language:[ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7] Cookie:[session-id=uoa7If53qWDgQdFR3DrGSeZXL9nLnBZT] Origin:[http://localhost:8080] Referer:[http://localhost:8080/] User-Agent:[Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36]] 0xc0002cfc20 <nil> 0 [] false 82.146.43.113:8080 map[] map[] <nil> map[] 93.171.198.4:37498 /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f 0xc000268c60 <nil> <nil> 0xc0000260d0}
// 2019/11/10 07:52:22 AuthMiddleware: passing to serve
// 2019/11/10 07:52:22 req.RequestURI = /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f
// 2019/11/10 07:52:22 req.Method = GET
// 2019/11/10 07:52:22 GetEmployerById id_string: %s
// 2019/11/10 07:52:22 GetEmployerById Handler Start
// 2019/11/10 07:52:22 GetEmployerById Parse id error: invalid UUID length: 0
// 2019/11/10 07:52:22 LOG END [GET] 93.171.198.4:37498, /employer/1668c0b9-653d-4a93-83b9-e8b32187c18f 10.800273ms

// 2019/11/10 07:52:22 http2: panic serving 93.171.198.4:37498: runtime error: slice bounds out of range [217:212]
// goroutine 112 [running]:
// net/http.(*http2serverConn).runHandler.func1(0xc0000105e8, 0xc0000b7f67, 0xc000248480)
// 	/usr/local/go/src/net/http/h2_bundle.go:5706 +0x16b
// panic(0x890e00, 0xc00030e140)
// 	/usr/local/go/src/runtime/panic.go:679 +0x1b2
// bufio.(*Reader).ReadSlice(0xc000060c00, 0x4b2b0a, 0xc0000cc800, 0x0, 0x0, 0xc0000cc800, 0x0)
// 	/usr/local/go/src/bufio/bufio.go:334 +0x22d
// github.com/gomodule/redigo/redis.(*conn).readLine(0xc000086dc0, 0xc0000b7630, 0xc0000b7630, 0x56be03, 0xc0000cc800, 0x0)
// 	/home/vl/hh/pkg/mod/github.com/gomodule/redigo@v2.0.0+incompatible/redis/conn.go:431 +0x38
// github.com/gomodule/redigo/redis.(*conn).readReply(0xc000086dc0, 0x0, 0x0, 0x0, 0x0)
// 	/home/vl/hh/pkg/mod/github.com/gomodule/redigo@v2.0.0+incompatible/redis/conn.go:504 +0x2f
// github.com/gomodule/redigo/redis.(*conn).DoWithTimeout(0xc000086dc0, 0x0, 0x8c4218, 0x3, 0xc0003063f0, 0x1, 0x1, 0x10, 0x823e40, 0x1, ...)
// 	/home/vl/hh/pkg/mod/github.com/gomodule/redigo@v2.0.0+incompatible/redis/conn.go:665 +0x154
// github.com/gomodule/redigo/redis.(*conn).Do(0xc000086dc0, 0x8c4218, 0x3, 0xc0003063f0, 0x1, 0x1, 0xc0000ccb00, 0x0, 0x400, 0x24)
// 	/home/vl/hh/pkg/mod/github.com/gomodule/redigo@v2.0.0+incompatible/redis/conn.go:616 +0x73
// 2019_2_IBAT/internal/pkg/auth/repository.(*SessionManager).Get(0xc000021180, 0xc000255f5b, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc00009dc00)
// 	/home/vl/hh/src/2019_2_IBAT/internal/pkg/auth/repository/auth_redstorage.go:37 +0x13e
// 2019_2_IBAT/internal/pkg/auth/service.(*AuthService).AuthMiddleware.func1(0x965280, 0xc0000105e8, 0xc000302900)
// 	/home/vl/hh/src/2019_2_IBAT/internal/pkg/auth/service/auth.go:61 +0x192
// net/http.HandlerFunc.ServeHTTP(0xc0002f4b40, 0x965280, 0xc0000105e8, 0xc000302900)
// 	/usr/local/go/src/net/http/server.go:2007 +0x44
// 2019_2_IBAT/internal/pkg/middleware.Logger.AccessLogMiddleware.func1(0x965280, 0xc0000105e8, 0xc000302900)
// 	/home/vl/hh/src/2019_2_IBAT/internal/pkg/middleware/logger.go:22 +0x1a2
// net/http.HandlerFunc.ServeHTTP(0xc0002f4b60, 0x965280, 0xc0000105e8, 0xc000302900)
// 	/usr/local/go/src/net/http/server.go:2007 +0x44
// 2019_2_IBAT/internal/pkg/middleware.CorsMiddleware.func1(0x965280, 0xc0000105e8, 0xc000302900)
// 	/home/vl/hh/src/2019_2_IBAT/internal/pkg/middleware/cors.go:103 +0x1d1
// net/http.HandlerFunc.ServeHTTP(0xc0002f4b80, 0x965280, 0xc0000105e8, 0xc000302900)
// 	/usr/local/go/src/net/http/server.go:2007 +0x44
// github.com/gorilla/mux.(*Router).ServeHTTP(0xc0000ce0c0, 0x965280, 0xc0000105e8, 0xc000302700)
// 	/home/vl/hh/pkg/mod/github.com/gorilla/mux@v1.7.3/mux.go:212 +0xe2
// net/http.serverHandler.ServeHTTP(0xc00013a2a0, 0x965280, 0xc0000105e8, 0xc000302700)
// 	/usr/local/go/src/net/http/server.go:2802 +0xa4
// net/http.initNPNRequest.ServeHTTP(0x9680c0, 0xc00026ce40, 0xc000274700, 0xc00013a2a0, 0x965280, 0xc0000105e8, 0xc000302700)
// 	/usr/local/go/src/net/http/server.go:3374 +0x8d
// net/http.(*http2serverConn).runHandler(0xc000248480, 0xc0000105e8, 0xc000302700, 0xc0002f4ae0)
// 	/usr/local/go/src/net/http/h2_bundle.go:5713 +0x9f
// created by net/http.(*http2serverConn).processHeaders
// 	/usr/local/go/src/net/http/h2_bundle.go:5447 +0x4eb
