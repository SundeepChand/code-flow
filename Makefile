build: build-wasm build-frontend
	git checkout gh-page
	git add .
	git commit -m 'Update build'
	git push origin gh-page

build-wasm: clean
	GOARCH=wasm GOOS=js go build -o cmd/web/frontend/public/assets/wasm/main.wasm github.com/SundeepChand/code-flow/cmd/web

build-frontend:
	cd cmd/web/frontend && npm run build:prod

run-frontend-dev-server:
	cd cmd/web/frontend && npm run serve

clean:
	rm -f cmd/web/frontend/assets/wasm/main.wasm
