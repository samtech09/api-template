APP_VER := "1.0.0"
PREVIOUS_VER := "0.0.0"

GOOS=linux
GOARCH=amd64
BUILDPATH=$(CURDIR)
EXENAME=api-template
BINPATH=$(BUILDPATH)/cmd
PUBLISHFOLDER=$(BUILDPATH)/releases
PUBLISHPATH=$(PUBLISHFOLDER)/$(APP_VER)
SCRIPTPATH=`echo ${GOPATH}/src/mahendras/scripts`
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
	#scp $(BINPATH)/${EXENAME}.zip linode-st:~/goapps/temp/
	#@sh $(SCRIPTPATH)/deploy-app.sh "${EXENAME}"
	echo "Package deployed to staging"



.PHONY: deploy-prod
## deploy-prod: build, package and upload binaries to s3 bucket
deploy-prod: --release --prepfiles-for-package
	# create zip file
	@7z a $(PUBLISHPATH)/${EXENAME}.zip $(PUBLISHPATH)/$(EXENAME) $(PUBLISHPATH)/conf.prod.json
	@echo Uploading...
	# #rename old app.zip to app_v1.2.3.zip at s3
	# @aws s3 mv s3://mepl-goapps/${EXENAME}.zip s3://mepl-goapps/${EXENAME}_v${PREVIOUS_VER}.zip 2>/dev/null; true
	# @aws s3 cp $(PUBLISHPATH)/${EXENAME}.zip s3://mepl-goapps/${EXENAME}.zip		#upload latest file to s3
	@echo "Zip package [${APP_VER}] deployed to s3 bucket"



.PHONY: deploy-prod-replace
## deploy-prod-replace: build, package and replace binaries at s3 bucket
deploy-prod-replace: --release-replace --prepfiles-for-package
	# create zip file
	@7z a $(PUBLISHPATH)/${EXENAME}.zip $(PUBLISHPATH)/$(EXENAME) $(PUBLISHPATH)/conf.prod.json $(PUBLISHPATH)/*.list
	@echo Uploading...
	@aws s3 rm s3://mepl-goapps/${EXENAME}.zip
	#upload latest file to s3s
	@aws s3 cp $(PUBLISHPATH)/${EXENAME}.zip s3://mepl-goapps/${EXENAME}.zip
	@echo "Zip package [${APP_VER}] replaced at s3 bucket"



# .PHONY: push-on-prod-servers
# ## push-on-prod-servers: ssh to each server and pull latest binaries from s3
# push-on-prod-servers:
# 	@sh $(SCRIPTPATH)/aws-push.sh "${EXENAME}"



.PHONY: deploy-git
## deploy-git: create new tag in git with current version number
deploy-git:
	@bash ~/scripts/new-release-tag.sh ${APP_VER}

	

.PHONY: gen-sql
## gen-sql: generate SQLs by invoking sql generater
gen-sql:
	@cd sqls/tools && go run .

	
.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
