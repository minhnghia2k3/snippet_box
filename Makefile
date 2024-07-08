.PHONY: run

docker.image:
	docker pull mysql

docker.run:
	docker run --name mysql_db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:tag

docker.exec:
	docker run exec -it mysql_db mysql -u root -p

run:
	go run ./cmd/web