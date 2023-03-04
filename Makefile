agent-build:
	cd ./agent/ && go build -o ../build/agent ./cmd/main.go && cd ..

agent-run: agent-build
	./build/agent -c ./config_local/agent.yml

admin-build:
	cd ./admin/ && go build -o ../build/admin ./cmd/main.go && cd ..

admin-run: admin-build
	./build/admin -c ./config_local/admin.yml

broker-build:
	cd ./broker/ && go build -o ../build/broker ./cmd/main.go && cd ..

broker-run: broker-build
	./build/broker -c ./config_local/broker.yml


ingress-build:
	cd ./ingress/ && go build -o ../build/ingress ./cmd/main.go && cd ..

ingress-run: ingress-build
	./build/ingress -c ./config_local/ingress.yml

storagea-build:
	cd ./storage_a/ && go build -o ../build/storagea ./cmd/main.go && cd ..

storagea-run: storagea-build
	./build/storagea -c ./config_local/storage.yml

storageb-build:
	cd ./storage_b/ && go build -o ../build/storageb ./cmd/main.go && cd ..

storageb-run: storageb-build
	./build/storageb -c ./config_local/storage.yml

storagef-build:
	cd ./storage_f/ && go build -o ../build/storagef ./cmd/main.go && cd ..

storagef-run: storagef-build
	./build/storagef -c ./config_local/storagef.yml