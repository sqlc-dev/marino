.PHONY: all parser clean

all: fmt parser generate

test: fmt parser
	sh test.sh

parser: parser/parser.go parser/hintparser.go

genkeyword: generate_keyword/genkeyword.go
	go build -C generate_keyword -o ../parser/genkeyword

generate: genkeyword parser/parser.y
	go generate ./parser/...

parser/parser.go: parser/parser.y bin/goyacc
	@echo "bin/goyacc -o $@ -p yy -t Parser $<"
	@bin/goyacc -o $@ -p yy -t Parser $< || ( rm -f $@ && echo 'Please check y.output for more information' && exit 1 )
	@rm -f y.output

parser/hintparser.go: parser/hintparser.y bin/goyacc
	@echo "bin/goyacc -o $@ -p yyhint -t hintParser $<"
	@bin/goyacc -o $@ -p yyhint -t hintParser $< || ( rm -f $@ && echo 'Please check y.output for more information' && exit 1 )
	@rm -f y.output

parser/%arser_golden.y: parser/%arser.y
	@bin/goyacc -fmt -fmtout $@ $<
	@(git diff --no-index --exit-code $< $@ && rm $@) || (mv $@ $< && >&2 echo "formatted $<" && exit 1)

bin/goyacc: goyacc/main.go goyacc/format_yacc.go
	GO111MODULE=on go build -o bin/goyacc goyacc/main.go goyacc/format_yacc.go

fmt: bin/goyacc parser/parser_golden.y parser/hintparser_golden.y
	@echo "gofmt (simplify)"
	@gofmt -s -l -w . 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'

clean:
	go clean -i ./...
	rm -rf *.out
	rm -f parser/parser.go parser/hintparser.go parser/genkeyword
