build:
	cd backend && go build .
	cd frontend && pnpm i && pnpm build

clean:
	cd backend && go clean
	cd frontend && rm -rf .parcel-cache/ dist/ node_modules/