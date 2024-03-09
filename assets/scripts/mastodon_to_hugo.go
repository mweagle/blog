package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"slices"
	"strings"
	"text/template"
	"time"
)

// Sample usage:
// go run mastodon_to_hugo.go --input "~/Downloads/mastodon-archive" --output "./blog/content/mastodon"

// Where the input folder is the root of the expanded Mastodon archive. It MUST
// include an outbox.json file
//
// The contents in the --output folder will be purged before rendering the
// new toots
//
// /////////////////////////////////////////////////////////////////////////////
// _                  _      _
// | |_ ___ _ __  _ __| |__ _| |_ ___ ___
// |  _/ -_) '  \| '_ \ / _` |  _/ -_|_-<
//  \__\___|_|_|_| .__/_\__,_|\__\___/__/
// 			  |_|
// /////////////////////////////////////////////////////////////////////////////

var TEMPLATE_TOOT_ROOT = `---
title: "Mastodon - {{ .Toot.Published }}"
subtitle: ""
canonical: {{ .Toot.Object.ID }}
description:
image: "/images/mastodon.png"

date: {{ .Toot.Published }}
lastmod: {{ .ExecutionTime }}
image: ""
tags: [{{ range $index, $eachTag := .Toot.Object.Tags}}{{if $index}},{{end}}"{{$eachTag.Name}}"{{end}}]

categories: ["mastodon"]
---
![Mastodon](/images/mastodon.png)
`

var TEMPLATE_TOOT = `
{{ .Toot.Object.Content }}
{{ range $index, $eachAttachment := .Toot.Object.Attachments}}
{{ if eq $eachAttachment.MediaType "video/mp4"}}<video controls autoplay muted loop width="512"><source src="{{$eachAttachment.BaseFilename}}" type="{{ $eachAttachment.MediaType}}" /></video>{{else}}![{{$eachAttachment.Name}}]({{$eachAttachment.BaseFilename}}){{end}}{{end}}

###### [Mastodon Source 🐘]({{ .Toot.Object.URL }})

___
`

// /////////////////////////////////////////////////////////////////////////////
// _            _
// __ ___ _ _  __| |_ __ _ _ _| |_ ___
// / _/ _ \ ' \(_-<  _/ _` | ' \  _(_-<
// \__\___/_||_/__/\__\__,_|_||_\__/__/
//
// /////////////////////////////////////////////////////////////////////////////

var HOST = "hachyderm.io"
var USER = "mweagle"
var MY_FOLLOWERS_URL = fmt.Sprintf("https://%s/users/%s/followers", HOST, USER)

// /////////////////////////////////////////////////////////////////////////////
// _
// | |_ _  _ _ __  ___ ___
// |  _| || | '_ \/ -_|_-<
//  \__|\_, | .__/\___/__/
// 	 |__/|_|
//
// /////////////////////////////////////////////////////////////////////////////

type FilterTootFunc func(*ActivityEntry) bool

// //////////////////////////////////////////////////////////////////////////////
// commandLineArgs
type commandLineArgs struct {
	inputRootPathExpandedArchive string
	outputRootPathHugoAssets     string
}

func (cla *commandLineArgs) parseCommandLine(log *slog.Logger) bool {
	flag.StringVar(&cla.inputRootPathExpandedArchive, "input", "", "Path to unzipped archive")
	flag.StringVar(&cla.outputRootPathHugoAssets, "output", "", "Path to root directory for output. Existing contents will be deleted.")
	flag.Parse()

	if (len(cla.inputRootPathExpandedArchive) <= 0) || len(cla.outputRootPathHugoAssets) <= 0 {
		log.Error("Invalid command line arguments")
		return false
	}
	return true
}

// /////////////////////////////////////////////////////////////////////////////
// publishingStats
type PublishingStats struct {
	totalTootCount    uint
	renderedTootCount uint
	filteredTootCount uint
	mediaFilesCount   uint
	replyThreadsCount uint
}

// /////////////////////////////////////////////////////////////////////////////
// ActivityObjectAttachment
type ActivityObjectAttachment struct {
	Type         string `json:"type"`
	MediaType    string `json:"mediaType"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	BaseFilename string
	AtomURI      string `json:"atomUri"`
	Width        uint   `json:"width"`
	Height       uint   `json:"height"`
}

// /////////////////////////////////////////////////////////////////////////////
// ActivityObjectTag
type ActivityObjectTag struct {
	Type string `json:"type"`
	Name string `json:"name"`
	HREF string `json:"href"`
}

// /////////////////////////////////////////////////////////////////////////////
// ActivityObject
type ActivityObject struct {
	Announcement string
	ID           string                      `json:"id"`
	Type         string                      `json:"type"`
	InReplyTo    string                      `json:"inReplyTo"`
	Published    string                      `json:"published"`
	URL          string                      `json:"url"`
	CC           []string                    `json:"cc"`
	AtomURI      string                      `json:"atomUri"`
	Content      string                      `json:"content"`
	Attachments  []*ActivityObjectAttachment `json:"attachment"`
	Tags         []*ActivityObjectTag        `json:"tag"`
}

func (ao *ActivityObject) UnmarshalJSON(data []byte) error {
	var s string
	stringUnmarshalErr := json.Unmarshal(data, &s)
	// If this succeeded, we need to ignore the rest of the data
	if stringUnmarshalErr == nil {
		ao.Announcement = s
	} else {
		dictMap := map[string]interface{}{}
		objUnmarshalErr := json.Unmarshal(data, &dictMap)
		if objUnmarshalErr != nil {
			return objUnmarshalErr
		}
		ao.ID = jsonScalar[string]("id", dictMap)
		ao.Type = jsonScalar[string]("type", dictMap)
		ao.InReplyTo = jsonScalar[string]("inReplyTo", dictMap)
		ao.Published = jsonScalar[string]("published", dictMap)
		ao.URL = jsonScalar[string]("url", dictMap)
		ao.AtomURI = jsonScalar[string]("atomUri", dictMap)
		ao.Content = jsonScalar[string]("content", dictMap)

		fieldValue, fieldValueExists := dictMap["cc"]
		if fieldValueExists {
			jsonBytes, _ := json.Marshal(fieldValue)
			fieldUnmarshalErr := json.Unmarshal(jsonBytes, &ao.CC)
			if fieldUnmarshalErr != nil {
				return fieldUnmarshalErr
			}
		}

		fieldValue, fieldValueExists = dictMap["attachment"]
		if fieldValueExists {
			jsonBytes, _ := json.Marshal(fieldValue)
			fieldUnmarshalErr := json.Unmarshal(jsonBytes, &ao.Attachments)
			if fieldUnmarshalErr != nil {
				return fieldUnmarshalErr
			}
			// For each one, update the BaseFilename to make the template
			// easier
			for _, eachAttachment := range ao.Attachments {
				urlPathParts := strings.Split(eachAttachment.URL, "/")
				eachAttachment.BaseFilename = urlPathParts[len(urlPathParts)-1]
			}
		}
		fieldValue, fieldValueExists = dictMap["tag"]
		if fieldValueExists {
			jsonBytes, _ := json.Marshal(fieldValue)
			fieldUnmarshalErr := json.Unmarshal(jsonBytes, &ao.Tags)
			if fieldUnmarshalErr != nil {
				return fieldUnmarshalErr
			}
		}
		// Always add a "Social Media" tag
		if len(ao.Tags) <= 0 {
			ao.Tags = make([]*ActivityObjectTag, 0)
		}
		ao.Tags = append(ao.Tags, &ActivityObjectTag{
			Type: "Hashtag",
			HREF: "https://hachyderm.io/tags/social media",
			Name: "Social Media",
		})
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////////////
// ActivityEntry
type ActivityEntry struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Published string          `json:"published"`
	CC        []string        `json:"cc"`
	Object    *ActivityObject `json:"object"`
}

// /////////////////////////////////////////////////////////////////////////////
// Outbox
type Outbox struct {
	TotalItems           uint             `json:"totalItems"`
	OrderedItems         []*ActivityEntry `json:"orderedItems"`
	ArchiveDirectoryRoot string
	// TODO - update this to be ID to ActivityEntry so that we can
	// partition the output files to include the date
	ThreadIDChain map[string]*ActivityEntry
}

func (ob *Outbox) filterToots(filterFunc FilterTootFunc) {
	filteredToots := []*ActivityEntry{}
	for _, eachEntry := range ob.OrderedItems {
		if filterFunc(eachEntry) {
			filteredToots = append(filteredToots, eachEntry)
		}
	}
	ob.OrderedItems = filteredToots
}

func jsonScalar[V any](key string, dict map[string]interface{}) V {
	curVal, curValOk := dict[key]
	if !curValOk {
		curVal = new(V)
	}

	typedVal, typedValOk := curVal.(V)
	if !typedValOk {
		return *new(V)
	}
	return typedVal
}

func selfPublishFilter(entry *ActivityEntry) bool {
	selfReplyToURL := fmt.Sprintf("https://%s/users/%s", HOST, USER)
	// Include only Create toots
	if entry.Type != "Create" {
		return false
	}
	// Include self-replies only
	if len(entry.Object.InReplyTo) != 0 &&
		!strings.HasPrefix(entry.Object.InReplyTo, selfReplyToURL) {
		return false
	}
	// ok, what about CCs
	if len(entry.Object.CC) > 1 || !slices.Contains(entry.Object.CC, MY_FOLLOWERS_URL) {
		return false
	}
	return true
}

func newOutbox(inputFile string) (*Outbox, error) {
	inputData, inputDataErr := os.ReadFile(inputFile)
	if inputDataErr != nil {
		return nil, inputDataErr
	}
	outbox := Outbox{}
	err := json.Unmarshal(inputData, &outbox)
	if err != nil {
		return nil, err
	}
	// Get the input file source. That's the root directory
	// for all media references
	outbox.ArchiveDirectoryRoot = path.Dir(inputFile)

	// For each activity, find the root thread element, which may be empty...
	outbox.ThreadIDChain = map[string]*ActivityEntry{}
	for _, eachActivity := range outbox.OrderedItems {
		outbox.ThreadIDChain[eachActivity.Object.ID] = eachActivity
	}
	return &outbox, nil
}

type cleanupFunc func(log *slog.Logger)

// /////////////////////////////////////////////////////////////////////////////
//  __              _   _
// / _|_  _ _ _  __| |_(_)___ _ _  ___
// |  _| || | ' \/ _|  _| / _ \ ' \(_-<
// |_|  \_,_|_||_\__|\__|_\___/_||_/__/
//
// /////////////////////////////////////////////////////////////////////////////

func ensureDirectory(root string, deleteExisting bool, log *slog.Logger) error {
	_, emptyDirectoryStatErr := os.Stat(root)
	log.Debug("Ensuring directory", "path", root, "deleteExisting", deleteExisting)
	if emptyDirectoryStatErr == nil && deleteExisting {
		removeAllErr := os.RemoveAll(root)
		log.Info("Deleting existing directory contents", "path", root)
		if removeAllErr != nil {
			return removeAllErr
		}
	}
	return os.MkdirAll(root, os.ModePerm)
}

func renderTootsToDisk(outputRoot string, filteredOutbox *Outbox, log *slog.Logger) error {
	// When rendering out, use the current time as the lastModTime
	nowTime := time.Now().Format(time.RFC3339)

	publishingStats := PublishingStats{
		totalTootCount:    filteredOutbox.TotalItems,
		renderedTootCount: uint(len(filteredOutbox.OrderedItems)),
		filteredTootCount: filteredOutbox.TotalItems - uint(len(filteredOutbox.OrderedItems)),
	}
	tootRootTemplate, tootRootTemplateErr := template.New("tootRoot").Parse(TEMPLATE_TOOT_ROOT)
	if tootRootTemplateErr != nil {
		return tootRootTemplateErr
	}
	tootTemplate, tootTemplateErr := template.New("toot").Parse(TEMPLATE_TOOT)
	if tootTemplateErr != nil {
		return tootTemplateErr
	}

	for _, eachItem := range filteredOutbox.OrderedItems {
		threadRootActivityItem := eachItem

		// By default, each toot is it's own root. If there is a replyTo chain,
		// recurse that to the root which becomes the active root
		//defaultID := eachItem.Object.ID
		// We need to keep recursing this to get to the root of the thread...
		for {
			replyToID := threadRootActivityItem.Object.InReplyTo
			if len(replyToID) <= 0 {
				break
			}
			parentActivityItem, parentActivityItemExists := filteredOutbox.ThreadIDChain[replyToID]
			if !parentActivityItemExists {
				// Skip it, it's probably a reply to someone...
				break
			}
			if parentActivityItem == threadRootActivityItem {
				return fmt.Errorf("Loop detected for item: %s", threadRootActivityItem.Object.ID)
			}
			threadRootActivityItem = parentActivityItem
			publishingStats.replyThreadsCount += 1
		}
		// Let's use the published information about the parent
		// to put some structure on the output.
		// Sample date: 2024-02-02T17:40:31Z
		parsedDate, parsedDateErr := time.Parse(time.RFC3339, threadRootActivityItem.Published)
		if parsedDateErr != nil {
			return fmt.Errorf("Failed to parse date: %s. Error: %s", threadRootActivityItem.Published, parsedDateErr)
		}
		idParts := strings.Split(threadRootActivityItem.Object.ID, "/")
		fileID := idParts[len(idParts)-1]
		tootRootBundleDirectory := path.Join(outputRoot,
			fmt.Sprintf("%d", parsedDate.Year()),
			fmt.Sprintf("%.2d", parsedDate.Month()),
			fileID,
		)
		// Is there a reply to entry?
		errDirectory := ensureDirectory(tootRootBundleDirectory, false, log)
		if errDirectory != nil {
			return errDirectory
		}
		tootOutputPath := path.Join(tootRootBundleDirectory, "index.md")
		log.Debug("Rendering toot", "id", eachItem.ID, "path", tootOutputPath)

		// Setup the template param map
		templateParamMap := map[string]interface{}{
			"ExecutionTime": nowTime,
			"Toot":          eachItem,
		}
		// Does the file already exist? If not, open in create and render the toot header. If it
		// does exist, open in append mode and render the next in the thread.
		var tootFS *os.File = nil
		_, fileExistsErr := os.Stat(tootOutputPath)
		if os.IsNotExist(fileExistsErr) {
			createFS, createFSErr := os.OpenFile(tootOutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if createFSErr != nil {
				return createFSErr
			}
			tootFS = createFS
			// If the file doesn't exist, render the toot header to the file...
			if err := tootRootTemplate.Execute(tootFS, templateParamMap); err != nil {
				return err
			}
		} else if fileExistsErr != nil {
			return fileExistsErr
		} else {
			appendFS, appendFSErr := os.OpenFile(tootOutputPath, os.O_APPEND|os.O_WRONLY, 0600)
			if appendFSErr != nil {
				return appendFSErr
			}
			log.Debug("Appending to thread",
				"replyTo", eachItem.Object.InReplyTo,
				"tootPath", tootOutputPath,
				"id", eachItem.Object.ID)
			tootFS = appendFS
		}

		// Then render the toot to the open file
		if err := tootTemplate.Execute(tootFS, templateParamMap); err != nil {
			return err
		}
		// Flush it
		tootFS.Close()

		// Any media objects we need to move? We're just going to use the basename for the
		// attachment and put it in the page bundle directory
		for _, eachAttachment := range eachItem.Object.Attachments {
			sourceFilePath := path.Join(filteredOutbox.ArchiveDirectoryRoot, eachAttachment.URL)
			destFilePath := path.Join(tootRootBundleDirectory, eachAttachment.BaseFilename)
			srcFile, srcFileErr := os.Open(sourceFilePath)
			if srcFileErr != nil {
				return srcFileErr
			}
			defer srcFile.Close()

			destFile, destFileErr := os.Create(destFilePath)
			if destFileErr != nil {
				return destFileErr
			}
			defer destFile.Close()
			bytesCopied, copyErr := io.Copy(destFile, srcFile) //copy the contents of source to destination file
			if copyErr != nil {
				return copyErr
			}
			log.Debug("Copied media file to source",
				"type", eachAttachment.MediaType,
				"name", eachAttachment.BaseFilename,
				"bytes", bytesCopied,
				"id", eachItem.Object.ID)
			publishingStats.mediaFilesCount += 1
		}
	}
	log.Info("Publishing statistics",
		"totalTootCount", publishingStats.totalTootCount,
		"renderedTootCount", publishingStats.renderedTootCount,
		"filteredTootCount", publishingStats.filteredTootCount,
		"replyThreadCount", publishingStats.replyThreadsCount,
		"mediaFilesCount", publishingStats.mediaFilesCount)
	return nil
}

//
////////////////////////////////////////////////////////////////////////////////

// //////////////////////////////////////////////////////////////////////////////
//
// _ __  __ _(_)_ _
// | '  \/ _` | | ' \
// |_|_|_\__,_|_|_||_|
//
// //////////////////////////////////////////////////////////////////////////////
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Welcome to Hugodon!")
	cleanupFuncs := []cleanupFunc{}

	cla := commandLineArgs{}
	if !cla.parseCommandLine(logger) {
		logger.Error("Failed to parse command line arguments")
		os.Exit(-1)
	}

	// Unmarshal the data into
	outboxFilePath := path.Join(cla.inputRootPathExpandedArchive, "outbox.json")
	outboxFeed, outboxFeedErr := newOutbox(outboxFilePath)
	if outboxFeedErr != nil {
		logger.Error("Failed to read output JSON", "path", outboxFilePath, "error", outboxFeedErr)
		os.Exit(-1)
	}
	totalToots := outboxFeed.TotalItems
	outboxFeed.filterToots(selfPublishFilter)
	logger.Info("Toots filtered", "totalCount", totalToots, "filteredCount", len(outboxFeed.OrderedItems))

	// Transform the toots into documents.
	ensureDirectory(cla.outputRootPathHugoAssets, true, logger)
	renderErr := renderTootsToDisk(cla.outputRootPathHugoAssets,
		outboxFeed,
		logger)
	if renderErr != nil {
		logger.Error("Failed to render toots", "error", renderErr)
		os.Exit(-1)
	}
	// Copy the tweets to the appropriate location
	for _, eachFunc := range cleanupFuncs {
		eachFunc(logger)
	}
	logger.Info("Toot replication complete")
}
