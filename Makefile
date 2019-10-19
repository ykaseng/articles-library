serve:
	docker-compose -f docker-compose.yml up --build

test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes

test_app:
	docker build -t test_image -f Dockerfile.test .
	docker run --name test_app --env-file ./config/config.development.articles-library.list --rm -it test_image

test_db:
	docker run --name test_db --rm -it -p 5432:5432 --env POSTGRES_USER=postgres --env POSTGRES_PASSWORD=r17BDyxd3rJmF9NIlGZP --env POSTGRES_DB=library postgres:latest

test_coverage:
	go test -timeout 30s github.com/ykaseng/articles-library/database -coverprofile=C:\Users\myhkaya\AppData\Local\Temp\vscode-gouP2GgF\go-code-cover