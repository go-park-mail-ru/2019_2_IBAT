package main

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func mainPage(w http.ResponseWriter, r *http.Request) {
// 	session, err := r.Cookie("session_id")
// 	loggedIn := (err != http.ErrNoCookie)

// 	if loggedIn {
// 		fmt.Fprintln(w, `<a href="/auth">logout</a>`)
// 		fmt.Fprintln(w, "Welcome, "+session.Value)

// 	} else {
// 		// fmt.Fprintln(w, `<form action="/auth" method="post">`)
// 		// fmt.Fprintln(w, "You need to login")

// 		// // Username:<input type="text" name="username">
// 		// // Password:<input type="password" name="password">
// 		// fmt.Fprintln(w, `<input type="submit" value="Login">`)
// 		// fmt.Fprintln(w, "</form>")
// 		fmt.Fprintln(w, `<a href="/auth">login</a>`)
// 		fmt.Fprintln(w, "You need to login")
// 	}
// }

// func loginPage(w http.ResponseWriter, r *http.Request) {
// 	expiration := time.Now().Add(10 * time.Hour)
// 	cookie := http.Cookie{
// 		Name:    "session_id",
// 		Value:   "rvasily",
// 		Expires: expiration,
// 	}
// 	http.SetCookie(w, &cookie)
// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// func logoutPage(w http.ResponseWriter, r *http.Request) {
// 	session, err := r.Cookie("session_id")
// 	if err == http.ErrNoCookie {
// 		http.Redirect(w, r, "/", http.StatusFound)
// 		return
// 	}

// 	session.Expires = time.Now().AddDate(0, 0, -1)
// 	http.SetCookie(w, session)

// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// // func main() {
// // 	r := mux.NewRouter()
// // 	// r.HandleFunc("/login", loginPage)
// // 	// r.HandleFunc("/logout", logoutPage)
// // 	r.HandleFunc("/auth", loginPage).Methods(http.MethodGet, http.MethodOptions)
// // 	r.HandleFunc("/auth", logoutPage).Methods(http.MethodDelete, http.MethodOptions)
// // 	http.HandleFunc("/", mainPage)

// // 	fmt.Println("starting server at :8080")
// // 	http.ListenAndServe(":8080", nil)
// // }
