all: t1 t2

t1 : t1.go
	go build t1.go
t2 : t2.go
	go build t2.go

t1.go : t1.py
	python c.py t1.py > t1.go
t2.go : t2.py
	python c.py t2.py > t2.go

clean:
	rm t*[0-9].go t*[0-9]
