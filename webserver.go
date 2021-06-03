package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)
//var mux = sync.Mutex{}

func GenerateData(value []string,timedata string) {
 var databytes int
	datasize := value[0]
	sizeUnit := datasize[len(datasize)-2:]
	size := datasize[:len(datasize)-2]
	sizeint,_ := strconv.Atoi(size)
	switch sizeUnit {
	case "KB":
		databytes = sizeint*1024
	case "MB":
		databytes = sizeint*1024*1024
	case "GB":
		databytes = sizeint*1024*1024*1024
	}

	fd, err := os.Create(timedata)
	if err != nil {
		log.Fatal("Failed to create output")
	}
	_, err = fd.Seek(int64(databytes)-1, 0)
	if err != nil {
		log.Fatal("Failed to seek")
	}
	_, err = fd.Write([]byte{8})
	if err != nil {
		log.Fatal("Write failed")
	}
	err = fd.Close()
	if err != nil {
		log.Fatal("Failed to close file")
	}
}



func helloWorld(w http.ResponseWriter, r *http.Request) {
	//mux.Lock()
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		tNow := time.Now()

		//time.Time to Unix Timestamp
		tUnix := tNow.Unix()
		filename := strconv.FormatInt(tUnix,10)
		for k, v := range r.URL.Query() {
			if k == "size"{

				GenerateData(v,filename)
			}
		}

		b,_ := ioutil.ReadFile(filename)
		sendByte := len(string(b))
		fmt.Println("Started transferring data of %d bytes",sendByte)
		t1 := time.Now()
		size,err := w.Write(b)
		if err != nil {
			fmt.Println(err)
		}
                t2 := time.Now()
		diff := t2.Sub(t1)
		fmt.Println("Finished Sending data of %d bytes sending time  %s ",sendByte,diff)
		os.Remove(filename)
		w.Write([]byte("Received a GET request\n"))
		w.Write([]byte(fmt.Sprintf("Remote Address: %s\n", r.RemoteAddr)))
		w.Write([]byte(fmt.Sprintf("Content Length: %d\n", size)))
		w.Write([]byte(fmt.Sprintf("Sending TIme duration : %s\n", diff)))
		w.Write([]byte(fmt.Sprintf("URL: %s\n", r.URL)))
	case "POST":
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Remote Address: %s\n", r.RemoteAddr)
		fmt.Printf("Content Length: %d\n", r.ContentLength)
		fmt.Printf("URL: %s\n", r.URL)
		w.Write([]byte("Received a POST request\n"))
		w.Write([]byte(fmt.Sprintf("Remote Address: %s\n", r.RemoteAddr)))
		w.Write([]byte(fmt.Sprintf("Content Length: %d\n", r.ContentLength)))
		w.Write([]byte(fmt.Sprintf("URL: %s\n", r.URL)))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
//	mux.Unlock()
}
func main() {
	go http.HandleFunc("/", helloWorld)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
