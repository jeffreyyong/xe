.PHONY : helm test all bi cover

all:
	make lint
	make test
	make helm

fix:
	bash script/build.sh fix

lint:
	bash script/build.sh lint

test:
	bash script/build.sh test

local_run:
	bash script/build.sh local_run

cover:
	bash script/build.sh cover
