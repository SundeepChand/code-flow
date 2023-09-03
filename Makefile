build-wasm: clean
	GOARCH=wasm GOOS=js go build -o cmd/web/frontend/assets/wasm/main.wasm github.com/SundeepChand/code-flow/cmd/web

clean:
	rm -f cmd/web/frontend/assets/wasm/main.wasm
