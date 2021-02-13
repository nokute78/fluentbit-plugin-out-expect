all:
	go build -buildmode=c-shared -o out_gexpect.so .

fast:
	go build flb_api.go

clean:
	rm -rf *.so *.h *~
