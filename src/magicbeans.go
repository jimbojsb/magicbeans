package main

import (
	"github.com/maxid/beanstalkd"
	"os"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

func main() {
	r := mux.NewRouter()

	if (len(os.Args) < 2) {
		fmt.Println("no host specified")
		os.Exit(-1)
	}

	r.HandleFunc("/host", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		response, _ := json.Marshal(os.Args[1])
		res.Write(response)
	})

	r.HandleFunc("/stats", func(res http.ResponseWriter, req *http.Request) {
		tubeNames, err := getClient().ListTubes()
		if (err != nil) {
			fmt.Println(err)
		}
		stats := make(map[string]map[string]string)
		for _, tubeName := range tubeNames {
			tubeStats, _ := getClient().StatsTube(tubeName)
			stats[tubeName] = tubeStats
		}
		res.Header().Add("Content-Type", "application/json")
		response, _ := json.Marshal(stats)
		res.Write(response)
	})

	r.HandleFunc("/peek-ready/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		client := getClient()
		client.Use(tube)
		readyJob, err := client.PeekReady()
		if (err != nil) {
			fmt.Println(err)
		}
		res.Header().Add("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]string {"job_id": strconv.Itoa(int(readyJob.Id)), "job_data": string(readyJob.Data)});
		res.Write(response)
	})

	r.HandleFunc("/peek-buried/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		client := getClient()
		client.Use(tube)
		readyJob, err := client.PeekBuried()
		if (err != nil) {
			fmt.Println(err)
		}
		res.Header().Add("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]string {"job_id": strconv.Itoa(int(readyJob.Id)), "job_data": string(readyJob.Data)});
		res.Write(response)
	})

	r.HandleFunc("/kick-all/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		client := getClient()
		tubeStats, err := client.StatsTube(tube)
		if (err != nil) {
			fmt.Println(err)
		}
		numToKick, _ := strconv.Atoi(tubeStats["current-jobs-buried"])
		client.Use(tube)
		client.Kick(numToKick)
	}).Methods("POST")

	r.HandleFunc("/kick/{tube}/{num}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		client := getClient()
		numToKick, _ := strconv.Atoi(vars["num"])
		client.Use(tube)
		client.Kick(numToKick)
	}).Methods("POST")


	r.HandleFunc("/bury-all/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		go func(tube string) {
			client := getClient()
			tubeStats, err := client.StatsTube(tube)
			if (err != nil) {
				fmt.Println(err)
			}
			numToBury, _ := strconv.Atoi(tubeStats["current-jobs-ready"])
			client.Watch(tube)
			for i := 0; i < numToBury; i++ {
				job, err := client.Reserve(30)
				if (err != nil) {
					fmt.Println(err)
				}
				client.Bury(job.Id, 1024)
			}
		}(tube)
	}).Methods("POST")

	r.HandleFunc("/bury/{tube}/{num}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		numToBury, _ := strconv.Atoi(vars["num"])
		go func(tube string, num int) {
			client := getClient()
			client.Watch(tube)
			for i := 0; i < num; i++ {
				job, err := client.Reserve(30)
				if (err != nil) {
					fmt.Println(err)
				}
				client.Bury(job.Id, 1024)
			}
		}(tube, numToBury)
	}).Methods("POST")

	r.HandleFunc("/drain-all-buried/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		go func(tube string) {
			client := getClient()
			tubeStats, err := client.StatsTube(tube)
			if (err != nil) {
				fmt.Println(err)
			}
			numToBury, _ := strconv.Atoi(tubeStats["current-jobs-buried"])
			client.Use(tube)
			for i := 0; i < numToBury; i++ {
				job, err := client.PeekBuried()
				if (err != nil) {
					fmt.Println(err)
				}
				client.Delete(job.Id)
			}
		}(tube)
	}).Methods("POST")

	r.HandleFunc("/drain-ready/{tube}/{num}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		numToBury, _ := strconv.Atoi(vars["num"])
		go func(tube string, num int) {
			client := getClient()
			client.Use(tube)
			for i := 0; i < num; i++ {
				job, err := client.PeekReady()
				if (err != nil) {
					fmt.Println(err)
				}
				client.Delete(job.Id)
			}
		}(tube, numToBury)
	}).Methods("POST")

	r.HandleFunc("/drain-all-ready/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		go func(tube string) {
			client := getClient()
			tubeStats, err := client.StatsTube(tube)
			if (err != nil) {
				fmt.Println(err)
			}
			numToBury, _ := strconv.Atoi(tubeStats["current-jobs-ready"])
			client.Use(tube)
			for i := 0; i < numToBury; i++ {
				job, err := client.PeekReady()
				if (err != nil) {
					fmt.Println(err)
				}
				client.Delete(job.Id)
			}
		}(tube)
	}).Methods("POST")

	r.HandleFunc("/drain-buried/{tube}/{num}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		numToBury, _ := strconv.Atoi(vars["num"])
		go func(tube string, num int) {
			client := getClient()
			client.Use(tube)
			for i := 0; i < num; i++ {
				job, err := client.PeekBuried()
				if (err != nil) {
					fmt.Println(err)
				}
				client.Delete(job.Id)
			}
		}(tube, numToBury)
	}).Methods("POST")

	r.HandleFunc("/insert/{tube}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tube := vars["tube"]
		req.ParseForm()
		jobData := req.FormValue("data");
		client := getClient()
		client.Use(tube)
		client.Put(1024, 0, 120, []byte(jobData))
	}).Methods("POST")


	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/magicbeans/www")))
	http.Handle("/", r)

	fmt.Println("Starting http service on :80")
	http.ListenAndServe(":80", nil)
}

func getClient() (*beanstalkd.BeanstalkdClient) {

	host := os.Args[1];

	var port string
	if (len(os.Args) == 3) {
		port = os.Args[2]
	} else {
		port = "11300"
	}

	connectionString := host + ":" + port

	beanstalk, err := beanstalkd.Dial(connectionString)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}
	return beanstalk
}