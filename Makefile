VERSION ?= latest
TAG ?= $(VERSION)


PORT = 8888

debug: tidy fmt
	PORT=$(PORT) COOKIE=$(shell cat ./cookie) go run ./cmd/genshin-level-up

open.local:
	open "http://localhost:8888/194435467"

open:
	open "https://genshin-level-up.vercel.app/194435467"

test: tidy fmt
	go test -v ./...

tidy:
	go mod tidy

fmt:
	goimports -w -l .

sync.%:
	pnpx ts-node ./scripts/sync-$*.ts

pnpm.install:
	pnpm i

ensure.genshin-data:
	if [[ -d ./GenshinData ]]; then \
  		cd ./GenshinData && git pull --rebase; \
  	else \
  	  	git clone --depth=1 https://github.com/Dimbreath/GenshinData.git ./GenshinData; \
  	fi

run.materials.scripts:
	pnpx ts-node ./scripts/materials.sync.ts

run.genshin-data.convert:
	pnpx ts-node ./scripts/genshin-data.convert.ts

sync.materials: pnpm.install run.materials.scripts

convert.genshin-data: ensure.genshin-data run.genshin-data.convert
