#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = models

GIT_GROUP = caltechlibrary

RELEASE_DATE=$(shell date +'%Y-%m-%d')

RELEASE_HASH=$(shell git log --pretty=format:'%h' -n 1)

PROGRAMS = modelgen

MAN_PAGES = modelgen.1

MAN_PAGES_MISC = model.5

VERSION = $(shell jq -r ".version" codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROGRAMS = $(shell ls -1 cmd)

PACKAGE = $(shell ls -1 *.go | grep -v 'version.go')

SUBPACKAGES = $(shell ls -1 */*.go)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
        PREFIX = $(prefix)
endif

OS = $(shell uname)

EXT =
ifeq ($(OS), Windows)
	EXT = .exe
endif

EXT_WEB = .wasm

DIST_FOLDERS = bin/* man/*

build: version.go $(PROGRAMS) man CITATION.cff about.md installer.sh installer.ps1

version.go: .FORCE
	@echo '' | pandoc --from t2t --to plain \
                --metadata-file codemeta.json \
                --metadata package=$(PROJECT) \
                --metadata version=$(VERSION) \
                --metadata release_date=$(RELEASE_DATE) \
                --metadata release_hash=$(RELEASE_HASH) \
                --template codemeta-version-go.tmpl \
                LICENSE >version.go
	-git add version.go


$(PROGRAMS): cmd/*/*.go $(PACKAGE)
	@mkdir -p bin
	go build -o bin/$@$(EXT) cmd/$@/$@.go
	@./bin/$@ -help >$@.1.md

man: $(MAN_PAGES) $(MAN_PAGES_LIB) $(MAN_PAGES_MISC)

$(MAN_PAGES): .FORCE
	mkdir -p man/man1
	pandoc $@.md --from markdown --to man -s >man/man1/$@

$(MAN_PAGES_MISC): .FORCE
	mkdir -p man/man5
	pandoc $@.md --from markdown --to man -s >man/man5/$@


CITATION.cff: codemeta.json .FORCE
	@cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	echo "" | pandoc --metadata title="Cite $(PROGRAM)" --metadata-file=_codemeta.json --template=codemeta-cff.tmpl >CITATION.cff

about.md: codemeta.json .FORCE
	@cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	echo "" | pandoc --metadata title="About $(PROGRAM)" --metadata-file=_codemeta.json --template=codemeta-about.tmpl >about.md

installer.sh: .FORCE
	@echo '' | pandoc --metadata title="Installer" --metadata git_org_or_person="$(GIT_GROUP)" --metadata-file codemeta.json --template codemeta-bash-installer.tmpl >installer.sh
	chmod 775 installer.sh
	git add -f installer.sh

installer.ps1: .FORCE
	@echo '' | pandoc --metadata title="Installer" --metadata git_org_or_person="$(GIT_GROUP)" --metadata-file codemeta.json --template codemeta-ps1-installer.tmpl >installer.ps1
	chmod 775 installer.ps1
	git add -f installer.ps1

# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv -v ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@for FNAME in $(MAN_PAGES); do if [ -f "./man/man1/$${FNAME}" ]; then cp -v "./man/man1/$${FNAME}" "$(PREFIX)/man/man1/$${FNAME}"; fi; done
	@for FNAME in $(MANPAGES_LIB); do if [ -f "./man/man3/$${FNAME}" ]; then cp -v "./man/man3/$${FNAME}" "$(PREFIX)/man/man3/$${FNAME}"; fi; done
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	@echo "Removing manpages in $(PREFIX)/man"
	@for FNAME in $(MAN_PAGES); do if [ -f "$(PREFIX)/man/man1/$${FNAME}" ]; then rm -v "$(PREFIX)/man/man1/$${FNAME}"; fi; done
	@for FNAME in $(MAN_PAGES_LIB); do if [ -f "$(PREFIX)/man/man3/$${FNAME}" ]; then rm -v "$(PREFIX)/man/man3/$${FNAME}"; fi; done


website: .FORCE
	make -f website.mak

check: .FORCE
	go vet *.go

test: clean build
	go test

cleanweb:
	@if [ -f index.html ]; then rm *.html; fi

clean:
	go clean
	@if [ -d bin ]; then rm -fR bin; fi
	@if [ -d dist ]; then rm -fR dist; fi
	@if [ -d testout ]; then rm -fR testout; fi
	-go clean -r

dist/Linux-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-x86_64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/Linux-aarch64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-aarch64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macOS-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-x86_64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macOS-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/Windows-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-x86_64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/Windows-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=arm64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

# Raspberry Pi OS (32 bit), as reported by Raspberry Pi 3B+
dist/Linux-armv7l:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-armv7l.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

	
distribute_docs:
	if [ -d dist ]; then rm -fR dist; fi
	mkdir -p dist
	cp -v codemeta.json dist/
	cp -v CITATION.cff dist/
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp installer.sh dist/
	cp installer.ps1 dist/
	cp -vR man dist/

update_version:
	$(EDITOR) codemeta.json
	codemeta2cff

release: clean build distribute_docs dist/Linux-x86_64 dist/Windows-x86_64 dist/Windows-arm64 dist/macOS-x86_64 dist/macOS-arm64 dist/Linux-aarch64 dist/Linux-armv7l

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish: website
	bash publish.bash

.FORCE:
