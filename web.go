package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	db  *pgxpool.Pool
	ctx context.Context
)

func StartServer(p *pgxpool.Pool, cont context.Context) {
	db = p
	ctx = cont
	router := mux.NewRouter()
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./ui/dist/js"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./ui/dist/css"))))
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/param", paramHandler)
	router.HandleFunc("/controller", controllerHandler)
	router.HandleFunc("/device", deviceHandler)
	router.HandleFunc("/project", projectHandler)
	router.HandleFunc("/address", addressHandler)
	router.HandleFunc("/event", eventHandler)
	router.HandleFunc("/statistic", statisticHandler)
	http.Handle("/", router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl, err := template.ParseFiles("ui/dist/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = tmpl.Execute(w, nil)
}

func statisticHandler(w http.ResponseWriter, r *http.Request) {

}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	q, err := db.Query(ctx, "select rn, code, event from event")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
	} else {
		var res []interface{}
		for q.Next() {
			t := struct {
				Rn    int    `json:"rn"`
				Code  string `json:"code"`
				Event string `json:"event"`
			}{}
			_ = q.Scan(&t.Rn, &t.Code, &t.Event)
			res = append(res, t)
		}
		out, _ := json.Marshal(res)
		_, _ = w.Write(out)
	}
}

func addressHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	switch r.Method {
	case http.MethodGet:
		var res []Address
		ids, ok := r.URL.Query()["id"]
		if id, _ := strconv.Atoi(ids[0]); ok {
			q, err := db.Query(ctx, "select rn, prn, code from address where prn = $1", id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
			} else {
				for q.Next() {
					tmp := Address{}
					q.Scan(&tmp.Rn, &tmp.Project, &tmp.Code)
					q2, err := db.Query(ctx, "select prn, val from address_param where arn = $1", tmp.Rn)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
					}
					for q2.Next() {
						t := struct {
							Prn int     `json:"rn"`
							Val float64 `json:"val"`
						}{}
						q2.Scan(&t.Prn, &t.Val)
						tmp.Params = append(tmp.Params, t)
					}
					res = append(res, tmp)
				}
				out, _ := json.Marshal(res)
				_, _ = w.Write(out)
			}
		}
	case http.MethodPut:
		tmp := Address{}
		err := decoder.Decode(&tmp)
		q, err := db.Query(ctx, "insert into address (prn, code) VALUES ($1, $2) returning rn", tmp.Project, tmp.Code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			q.Next()
			q.Scan(&tmp.Rn)
			for _, param := range tmp.Params {
				_, err = db.Exec(ctx, "insert into address_param values ($1, $2, $3)", param.Prn, tmp.Rn, param.Val)
			}
		}
	case http.MethodPost:
		tmp := Address{}
		err := decoder.Decode(&tmp)
		_, err = db.Query(ctx, "update address set prn = $1, code = $2 where rn = $3", tmp.Project, tmp.Code, tmp.Rn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			db.Exec(ctx, "delete from address_params where arn = $1", tmp.Rn)
			for _, param := range tmp.Params {
				_, _ = db.Exec(ctx, "insert into address_params values ($1, $2, $3)", param.Prn, tmp.Rn, param.Val)
			}
		}
	case http.MethodDelete:
		tmp := Address{}
		_ = decoder.Decode(&tmp)
		_, err := db.Query(ctx, "delete from address where rn = $1", tmp.Rn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			_, err = db.Query(ctx, "delete from address_param where arn = $1", tmp.Rn)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
			}
		}
	case http.MethodOptions:
		var res []interface{}
		q, err := db.Query(ctx, "select rn, code from address")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			for q.Next() {
				tmp := struct {
					Rn   int    `json:"value"`
					Code string `json:"text"`
				}{}
				_ = q.Scan(&tmp.Rn, &tmp.Code)
				res = append(res, tmp)
			}
			if out, err := json.Marshal(res); err == nil {
				_, _ = w.Write(out)
			}
		}
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	switch r.Method {
	case http.MethodGet:
		var res []Project
		q, err := db.Query(ctx, "select rn, code from project")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			for q.Next() {
				tmp := Project{}
				q.Scan(&tmp.Rn, &tmp.Code)
				q2, _ := db.Query(ctx, "select prn, val from project_param where pjrn = $1", tmp.Rn)
				for q2.Next() {
					t := struct {
						Prn int     `json:"rn"`
						Val float64 `json:"val"`
					}{}
					q2.Scan(&t.Prn, &t.Val)
					tmp.Params = append(tmp.Params, t)
				}
				res = append(res, tmp)
			}
			out, _ := json.Marshal(res)
			_, _ = w.Write(out)
		}
	case http.MethodPut:
		tmp := Project{}
		err := decoder.Decode(&tmp)
		q, err := db.Query(ctx, "insert into project(code) values ($1) returning rn", tmp.Code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			rn := 0
			q.Next()
			_ = q.Scan(&rn)
			for _, param := range tmp.Params {
				_, err = db.Exec(ctx, "insert into project_param values ($1, $2, $3)", param.Prn, rn, param.Val)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
				}
			}
		}
	case http.MethodPost:
		tmp := Project{}
		_ = decoder.Decode(&tmp)
		_, err := db.Query(ctx, "update project set code = $1 where rn = $2", tmp.Code, tmp.Rn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			_, _ = db.Query(ctx, "delete from project_param where pjrn = $1", tmp.Rn)
			for _, param := range tmp.Params {
				_, _ = db.Exec(ctx, "insert into project_param values ($1, $2, $3)", param.Prn, tmp.Rn, param.Val)
			}
		}
	case http.MethodDelete:
		tmp := Project{}
		_ = decoder.Decode(&tmp)
		_, _ = db.Query(ctx, "delete from project where rn = $1", tmp.Rn)
		_, _ = db.Query(ctx, "delete from project_param where pjrn = $1", tmp.Rn)
	case http.MethodOptions:
		ids, ok := r.URL.Query()["id"]
		if id, _ := strconv.Atoi(ids[0]); ok {
			var res []interface{}
			q, _ := db.Query(ctx, "select prn, val from project_param where pjrn = $1", id)
			for q.Next() {
				t := struct {
					Prn int     `json:"prn"`
					Val float64 `json:"val"`
				}{}
				q.Scan(&t.Prn, &t.Val)
				res = append(res, t)
			}
			if out, err := json.Marshal(res); err == nil {
				_, _ = w.Write(out)
			}
		}
	}
}

func deviceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	switch r.Method {
	case http.MethodGet:
		var res []interface{}
		ids, ok := r.URL.Query()["id"]
		if id, err := strconv.Atoi(ids[0]); ok && err == nil {
			q, err := db.Query(ctx, "select d.rn, d.comment, d.uid, a.rn, a.code\nfrom devices d inner join location l on d.rn = l.drn inner join address a on a.rn = l.arn\nwhere d.crn = $1", id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
			} else {
				for q.Next() {
					tmp := struct {
						Rn      int    `json:"rn"`
						Uid     string `json:"uid"`
						Comment string `json:"comment"`
						Address struct {
							Rn   int    `json:"value"`
							Code string `json:"text"`
						} `json:"address"`
					}{}
					_ = q.Scan(&tmp.Rn, &tmp.Comment, &tmp.Uid, &tmp.Address.Rn, &tmp.Address.Code)
					res = append(res, tmp)
				}
				if out, err := json.Marshal(res); err == nil {
					_, _ = w.Write(out)
				}
			}
		}
	case http.MethodPut:
		tmp := struct {
			Rn      int
			Crn     int    `json:"crn"`
			Uid     string `json:"uid"`
			Comment string `json:"comment"`
			Arn     int    `json:"arn"`
		}{}
		_ = decoder.Decode(&tmp)
		q, err := db.Query(ctx, "insert into devices (crn, comment, uid) VALUES ($1, $2, $3) returning rn",
			tmp.Crn, tmp.Comment, tmp.Uid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			q.Next()
			q.Scan(&tmp.Rn)
			_, err := db.Query(ctx, "insert into location(drn, arn) VALUES ($1, $2)", tmp.Rn, tmp.Arn)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
			}
		}
	case http.MethodPost:
		tmp := struct {
			Rn      int    `json:"rn"`
			Crn     int    `json:"crn"`
			Uid     string `json:"uid"`
			Comment string `json:"comment"`
			Arn     int    `json:"arn"`
		}{}
		_ = decoder.Decode(&tmp)
		_, err := db.Query(ctx, "update devices set crn = $1, comment = $2, uid = $3 where rn = $4",
			tmp.Crn, tmp.Comment, tmp.Uid, tmp.Rn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			_, err := db.Query(ctx, "insert into location(drn, arn) VALUES ($1, $2)", tmp.Rn, tmp.Arn)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
			}
		}
	case http.MethodDelete:
		tmp := struct {
			Rn int `json:"rn"`
		}{}
		_ = decoder.Decode(&tmp)
		_, _ = db.Query(ctx, "delete from devices where rn = $1", tmp.Rn)
		_, _ = db.Query(ctx, "delete from location where drn = $1", tmp.Rn)
	}
}

func controllerHandler(w http.ResponseWriter, r *http.Request) {
	var obj interface{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&obj)
	controller := struct {
		Rn     int    `json:"rn"`
		Code   string `json:"code"`
		Params []int  `json:"params"`
	}{}
	switch r.Method {
	case http.MethodGet:
		var res []interface{}
		q, err := db.Query(ctx, "select rn, code from controller")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		} else {
			for q.Next() {
				tc := controller
				_ = q.Scan(&tc.Rn, &tc.Code)
				qp, err := db.Query(ctx, "select rn from parameter p inner join controller_parameter cp on p.rn = cp.prn where cp.crn = $1", tc.Rn)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
				}
				for qp.Next() {
					i := 0
					qp.Scan(&i)
					tc.Params = append(tc.Params, i)
				}
				res = append(res, tc)
			}
			if out, err := json.Marshal(res); err == nil {
				w.Write(out)
			}
		}
	case http.MethodPut:
		qr, err := db.Query(ctx, "insert into controller(code) values ($1) returning rn", obj.(map[string]interface{})["code"].(string))
		if qr.Next() && err == nil {
			rn := 0
			qr.Scan(&rn)
			for i, p := range obj.(map[string]interface{})["params"].([]interface{}) {
				db.Exec(ctx, "insert into controller_parameter(crn, prn, pos) values ($1, $2, $3)", rn, p.(float64), i+1)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		}
	case http.MethodPost:
		js := obj.(map[string]interface{})
		_, err := db.Exec(ctx, "update controller set code = $1 where rn = $2", js["code"].(string), js["rn"].(float64))
		if err == nil {
			db.Exec(ctx, "delete from controller_parameter where crn = $1", js["rn"].(float64))
			for i, p := range obj.(map[string]interface{})["params"].([]interface{}) {
				db.Exec(ctx, "insert into controller_parameter(crn, prn, pos) values ($1, $2, $3)",
					js["rn"].(float64), p.(float64), i+1)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		}
	case http.MethodDelete:
		db.Exec(ctx, "delete from controller where rn = $1", obj.(map[string]interface{})["rn"].(float64))
		db.Exec(ctx, "delete from controller_parameter where crn = $1", obj.(map[string]interface{})["rn"].(float64))
	}
}

func paramHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	param := Param{}
	_ = decoder.Decode(&param)
	switch r.Method {
	case http.MethodGet:
		q, err := db.Query(ctx, "select rn, code, formula from parameter where formula is null")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500 - Something bad happened!"))
		} else {
			var res []Param
			for q.Next() {
				tmp := Param{}
				_ = q.Scan(&tmp.Rn, &tmp.Code, &tmp.Formula)
				res = append(res, tmp)
			}
			out, _ := json.Marshal(res)
			w.Write(out)
		}
	case http.MethodPut:
		_, err := db.Exec(ctx, "insert into parameter (code, formula) VALUES ($1, $2)", param.Code, param.Formula)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		}
	case http.MethodPost:
		_, err := db.Exec(ctx, "update parameter set code = $1, formula = $2 where rn = $3", param.Code, param.Formula, param.Rn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		}
	case http.MethodDelete:
		_, err := db.Exec(ctx, "delete from parameter where rn = $1", param.Rn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprint(err.Error())))
		}
	}
}
