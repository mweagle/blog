default: build

HUGO_VERSION := "0.112.3"
HUGO_SOURCE_URL := "https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VERSION)/hugo_extended_$(HUGO_VERSION)_darwin-universal.tar.gz"
.PHONY: run edit

hugo_install:
	curl -L -vs --create-dirs -O --output-dir ./.vscode $(HUGO_SOURCE_URL)
	tar -xvf ./.vscode/hugo_extended_$(HUGO_VERSION)_darwin-universal.tar.gz --directory ./.vscode
	mv -fv ./.vscode/hugo ./


clean:
	rm -rfv ./public

ext_assets: clean
	rm -rfv ./content/posts/c4pumlthemes/puml/resources/palettes
	mkdir -pv ./content/posts/c4pumlthemes/puml/resources/palettes
	cp -Rv ~/Documents/GitHub/C4-PlantUML-Themes/palettes ./content/posts/c4pumlthemes/puml/resources

build: ext_assets clean
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
	open http://localhost:1313
	LOCAL_DEVELOPMENT=TRUE ./hugo server --watch --disableFastRender --forceSyncStatic --noHTTPCache --buildDrafts