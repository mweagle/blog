default: build
.PHONY: run edit

# Execute time
YEAR=`date +'%Y'`
MONTH=`date +'%m'`
DAY=`date +'%m'`
TIME=`date +'%H_%M_%S%z'`
POST_DIR=posts/$(YEAR)/$(MONTH)/$(TIME)
CONTENT_POSTS_DIR=content/$(POST_DIR)
NEW_POST_FILENAME=index.md
NOW_TIME=$(YEAR)-$(MONTH)-$(DAY)T$(TIME)

# Hugo requirements.
# Ensure when this is updated we also regenerate the netlify toml
# file://./netlify.toml
HUGO_VERSION := "0.147.5"
HUGO_SOURCE_URL := "https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VERSION)/hugo_extended_$(HUGO_VERSION)_darwin-universal.tar.gz"

# External resources that need to be rebuilt
STAGING_DIR := "./.staging"
GITHUB_STAGING_DIR := "$(STAGING_DIR)/github"

hugo_install: netlify_toml
	curl -L -vs --create-dirs -O --output-dir ./.vscode $(HUGO_SOURCE_URL)
	tar -xvf ./.vscode/hugo_extended_$(HUGO_VERSION)_darwin-universal.tar.gz --directory ./.vscode
	mv -fv ./.vscode/hugo ./
	./hugo version

clean:
	rm -rfv ./public
	rm -rfv "$(GITHUB_STAGING_DIR)/C4-PlantUML-Themes"

post:
	mkdir -pv $(CONTENT_POSTS_DIR)
	@echo "POST FILENAME: $(CONTENT_POSTS_DIR)/$(NEW_POST_FILENAME)"
	./hugo new content $(POST_DIR)/$(NEW_POST_FILENAME)
	code "$(CONTENT_POSTS_DIR)/$(NEW_POST_FILENAME)"

# ARCHIVE_PATH=/Users/mattweagle/Downloads/mastodon-archive.zip make mastodon_replication
mastodon_replication:
	rm -rfv ${TMPDIR}/mastodon
	unzip -d ${TMPDIR}/mastodon $(ARCHIVE_PATH)
	rm -rfv ./content/mastodon
	mkdir -pv ./content/mastodon
	go run "../mastodon-to-hugo/mastodon-to-hugo.go" --input ${TMPDIR}/mastodon --output "./content/mastodon"
#
# TODO: https://christianspecht.de/2020/08/10/creating-an-image-gallery-with-hugo-and-lightbox2/
#
# Ensure that the Netlify config uses the same Hugo
# version we've configured in the Makefile. Output:
#
# # Created: 2024-01-03T11:56:31-0800
# [build]
# command = "hugo"
# publish = "public"
# [build.environment]
# HUGO_VERSION = "0.121.1"
#
netlify_toml:
	 @echo "# Created: $(NOW_TIME)\n[build]\ncommand = \"hugo\"\npublish = \"public\"\n\n[build.environment]\nHUGO_VERSION = \"$(HUGO_VERSION)\"" > netlify.toml

ext_assets: clean netlify_toml
	git clone --depth=1 https://github.com/mweagle/C4-PlantUML-Themes "$(GITHUB_STAGING_DIR)/C4-PlantUML-Themes"
	rm -rfv ./content/posts/c4pumlthemes/puml/resources/palettes
	mkdir -pv ./content/posts/c4pumlthemes/puml/resources/palettes
	cp -Rv "$(GITHUB_STAGING_DIR)/C4-PlantUML-Themes/palettes" ./content/posts/c4pumlthemes/puml/resources

template_overrides:
# Copy the search override file into the theme directory so that the
# mastodon and twitter contents are searchable, but don't show up in the
# chronological feed. See config.yaml for details.
	cp -fv ./layouts/page/search.json ./themes/hugo-theme-stack/layouts/page/search.json

build: clean ext_assets template_overrides
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

edit: clean netlify_toml
	# Used for localhost editing
	open http://localhost:1313
	LOCAL_DEVELOPMENT=TRUE ./hugo server --watch --disableFastRender --forceSyncStatic --noHTTPCache --buildDrafts

netlify:
	open "https://app.netlify.com/teams/mweagle/overview"
