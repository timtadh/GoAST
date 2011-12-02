build: clean
	gobuild -a

test: build
	-gobuild -t -run
	rm _testmain _testmain.6 _testmain.go

fmt: clean
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w .

.PHONY : clean
clean :
	-rm -r goast/main
	-rm _testmain _testmain.6 _testmain.go
	-find . -name "*.6" | xargs -I"%s" rm %s
	-find . -name "*.a" | xargs -I"%s" rm %s

