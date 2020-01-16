APP_VER := "1.0.0"
PREVIOUS_VER := "0.0.0"

GOOS=linux
GOARCH=amd64
BUILDPATH=$(CURDIR)
EXENAME=api-template
BINPATH=$(BUILDPATH)/cmd
PUBLISHFOLDER=$(BUILDPATH)/releases
PUBLISHPATH=$(PUBLISHFOLDER)/$(APP_VER)
#APP_BRANCH := $(shell git branch | grep '*')
#DEP_BRANCH := $(shell cd ../stapi/models ; git branch | grep '*')


# check and app branch and depended modules are on same branch
#.PHONY: check
#check:
#	@if [ "$(APP_BRANCH)" != "$(DEP_BRANCH)" ] ; then echo "App is on branch: $(APP_BRANCH), while dependend module stapimodels is on $(depbranch) branch"; exit 1; fi


.PHONY: default
default: help ;



.PHONY: clean
## clean: clean the binaries
clean:
	@rm -rf $(BINPATH)


	
.PHONY: build
## build: build application for current OS i.e. Linux
build: clean
	@if [ ! -d $(BINPATH) ] ; then mkdir -p $(BINPATH) ; fi
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.appVer=$(APP_VER) -X main.buildVer=`date -u +b-%Y%m%d.%H%M%S`" -o $(BINPATH)/$(EXENAME) . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'


.PHONY: run
## run: build and run application
run: build
	@cp conf.dev.json $(BINPATH)/
	@mkdir $(BINPATH)/sqls || true
	@cp sqls/sqlbuilder.json $(BINPATH)/sqls
	@cd ${BINPATH} && export TESTENV=1 && ./${EXENAME}


.PHONY: test
## test: run unit test for given controller, pass name of comtroller as c=controllername e.g. make c=user test
test:
	@echo Controller is $(c)
	@cd tests/$(c)tests && go test


--release: build
	@echo 'Releasing new version ${APP_VER} ...'
	@if [ -d $(PUBLISHPATH) ] ; then echo 'Build folder already exist' ; exit 1 ; fi
	@mkdir -p $(PUBLISHPATH)


--release-replace: build
	@echo 'Replacing version ${APP_VER} ...'
	@if [ ! -d $(PUBLISHPATH) ] ; then echo 'Build folder not exist' ; exit 1 ; fi


--prepfiles-for-package:
	@cp $(BINPATH)/$(EXENAME) $(PUBLISHPATH)/
	@cp conf.test.json $(PUBLISHPATH)/conf.dev.json
	@cp $(PUBLISHFOLDER)/conf.prod.json $(PUBLISHPATH)/
	@# delete old zip file
	@rm $(PUBLISHPATH)/${EXENAME}.zip || true
	@echo ...files synced



.PHONY: deploy-test
## deploy-test: build and deploy binaries to staging server
deploy-test: build
	@cp conf.test.json $(BINPATH)/conf.dev.json
	@7z a $(BINPATH)/${EXENAME}.zip $(BINPATH)/$(EXENAME) $(BINPATH)/conf.dev.json
	echo "Package ready to deploy to staging"



.PHONY: deploy-prod
## deploy-prod: build, package and upload binaries to s3 bucket
deploy-prod: --release --prepfiles-for-package
	# create zip file
	@7z a $(PUBLISHPATH)/${EXENAME}.zip $(PUBLISHPATH)/$(EXENAME) $(PUBLISHPATH)/conf.prod.json
	@echo "Zip package [${APP_VER}] ready to deploy to production"



.PHONY: deploy-prod-replace
## deploy-prod-replace: build, package and replace binaries at s3 bucket
deploy-prod-replace: --release-replace --prepfiles-for-package
	# create zip file
	@7z a $(PUBLISHPATH)/${EXENAME}.zip $(PUBLISHPATH)/$(EXENAME) $(PUBLISHPATH)/conf.prod.json $(PUBLISHPATH)/*.list
	@echo "Replacement package ready to deploy to production"
	

.PHONY: gen-sql
## gen-sql: generate SQLs by invoking sql generater
gen-sql:
	@cd sqls/tools && go run .

	
.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
