.PHONY: build dev install clean

build: clean install dist/s3backup

dist/s3backup: src/main.ts
	pnpm esbuild src/main.ts --bundle --minify --platform=node --outfile=dist/s3backup

dev: clean install
	pnpm esbuild src/main.ts --bundle --watch  --sourcemap --platform=node --outfile=dist/s3backup

install:
	pnpm i

clean:
	rm -rf dist
