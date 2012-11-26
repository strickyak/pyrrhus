all: t1 t2 t3 t4 t5

run_grok:
	go run grok.go /opt/go/src/pkg/*/

t1 : t1.go
	go build t1.go runtime.go
	./t1
t1.go : t1.py
	python c.py t1.py > t1.go || { mv t1.go t1.bad ; false ; }

t2 : t2.go
	go build t2.go runtime.go
	./t2
t2.go : t2.py
	python c.py t2.py > t2.go || { mv t2.go t2.bad ; false ; }

t3 : t3.go
	go build t3.go runtime.go
	./t3
t3.go : t3.py
	python c.py t3.py > t3.go || { mv t3.go t3.bad ; false ; }

t4 : t4.go
	go build t4.go runtime.go
	./t4
t4.go : t4.py
	python c.py t4.py > t4.go || { mv t4.go t4.bad ; false ; }

t5 : t5.go
	go build t5.go runtime.go
	./t5
t5.go : t5.py
	python c.py t5.py > t5.go || { mv t5.go t5.bad ; false ; }

clean:
	rm -f t*[0-9].go t*[0-9] t*[0-9].bad
