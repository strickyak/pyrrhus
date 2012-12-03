# Simple 'Hello World' web server.

from go import io  # hack for io.Writer
#from go import os
from go import fmt
from go.net import http

def home(w, r):
  fmt.Fprintf(w, 'Hello World!')

print http.HandleFunc('/', home)
print http.ListenAndServe('127.0.0.1:1111', None)

# Based on this golang code:
##  package main
##  
##  import (
##    "os"
##    "fmt"
##    "net/http"
##  )
##  
##  func home(w http.ResponseWriter, r *http.Request) {
##    fmt.Fprintf(w, "Hello World!")
##  }
##  
##  func main() {
##    port := os.Args[1]
##    http.HandleFunc("/", home)
##    http.ListenAndServe("127.0.0.1:" + port, nil)
##  }
