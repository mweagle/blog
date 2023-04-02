default: build

HUGO_VERSION := "0.110.0"
HUGO_SOURCE_URL := "https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VERSION)/hugo_extended_$(HUGO_VERSION)_darwin-universal.tar.gz"
.PHONY: run edit

hugo_install:
	curl -L -vs --create-dirs -O --output-dir ./.vscode $(HUGO_SOURCE_URL)
	tar -xvf ./.vscode/hugo_extended_$(HUGO_VERSION)_darwin-universal.tar.gz --directory ./.vscode
	mv -fv ./.vscode/hugo ./

clean:
	rm -rf ./public

build: clean
	./hugo

commit:
	git add --all . && git commit

commit-nomessage:
	git add --all . && git commit -m "Updated documentation"

pull:
	git pull origin main

push:
	git push -f origin main

publish: commit-nomessage push
	echo "Docs published"

commit-auto: commit-nomessage
	git push

edit: clean
	# Used for localhost editing
	LOCAL_DEVELOPMENT=TRUE ./hugo server --watch