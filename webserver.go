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

var flag503 bool
var flagnoresp bool

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



func generate503(w http.ResponseWriter, r *http.Request) {

fmt.Println("Header Data Start")
        for name, values := range r.Header {
    // Loop over all values for the name.
    for _, value := range values {
        fmt.Println(name, value)
    }
}
fmt.Println("Header Data End")

        switch r.Method {
        case "GET":
                for k, v := range r.URL.Query() {
                        if k == "alt503" && v[0] == "true" && flag503 == true {
                         				
			 flag503 = false
			     fmt.Println("200 sent")
			     w.WriteHeader(http.StatusOK)
                             w.Write([]byte("200 - Something good happened!"))
                             return

                        } else if k == "alt503" && v[0] == "true" {
			 flag503 = true
			     fmt.Println("503 sent")
			     w.WriteHeader(http.StatusServiceUnavailable)
                             w.Write([]byte("503 - Something bad happened!"))
                             return
			}
                }
			     w.WriteHeader(http.StatusServiceUnavailable)
			     fmt.Println("503 sent")
                             w.Write([]byte("503 - Something bad happened!"))
        default:
                w.WriteHeader(http.StatusNotImplemented)
                w.Write([]byte(http.StatusText(http.StatusNotImplemented)))

}


}



func noresp(w http.ResponseWriter, r *http.Request) {

	
fmt.Println("Header Data Start")
        for name, values := range r.Header {
    // Loop over all values for the name.
    for _, value := range values {
        fmt.Println(name, value)
    }
}
fmt.Println("Header Data End")

        switch r.Method {
        case "GET":
                for k, v := range r.URL.Query() {
                        if k == "altempty" && v[0] == "true" && flagnoresp == true {
                             flagnoresp = false
                             fmt.Println("200 sent")
                             w.WriteHeader(http.StatusOK)
                             w.Write([]byte("200 - Something good happened!"))
                             return
			} else if  k == "altempty" && v[0] == "true" {
                             flagnoresp = true

		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Close = true
		// Don't forget to close the connection:
		conn.Close()
		} else {

		r.Close = true
		}

	     } 

        default:
                w.WriteHeader(http.StatusNotImplemented)
                w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	//mux.Lock()
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

fmt.Println("Header Data Start")	
	for name, values := range r.Header {
    // Loop over all values for the name.
    for _, value := range values {
        fmt.Println(name, value)
    }
}
fmt.Println("Header Data End")	

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
	go http.HandleFunc("/503", generate503)
	go http.HandleFunc("/noresp", noresp)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
