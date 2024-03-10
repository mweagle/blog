---
title: "Archiving Toots to Hugo Sites"
subtitle: ""
description: "Own a copy of your own content" 
draft: false

date: 2024-03-09T21:11:33-08:00
lastmod: 2024-03-09T21:11:33-08:00
image: "/posts/2024/03/mastodon-to-hugo/mastodon.svg"
tags: ["Social Media"]
categories: [""]

---

## TL;DR;

Run a go script to archive your Mastodon archive to your Hugo site
([source](https://github.com/mweagle/blog/blob/main/assets/scripts/mastodon_to_hugo.go)).
Visit this site's [Mastodon category](https://mweagle.net/categories/mastodon/) for an illustration.

## PESOS

The rise of social media has dramatically reduced the cost of authoring and sharing information. One downside
of this movement is that I've relinquished content ownership to each site's admins or owners. No doubt there are
benefits to this arrangement, among them being I don't need to manage the technical bits, but it would be
nice for me to maintain a copy in case of the worst. I learned this is abbreviated [PESOS](https://indieweb.org/PESOS)
(publish to third party, archive to self) and while it's less desireable than [POSSE](https://indieweb.org/POSSE)
(publish to self, syndicate to third parties), it's something I'd like to keep doing.

Taking inspiration from [Ross Baker's](https://social.rossabaker.com/@ross) Mastodon to Hugo
[archive script](https://rossabaker.com/projects/toot-archive/), I decided to create a `go` version of his
[Scala script](https://git.rossabaker.com/ross/cromulent/src/branch/main/gen/scala/toot-archive.scala). It's been a
while since I wrote Scala on a regular basis and while I could read Ross' code, I'm not confident I'd be able to
make changes that preserved the pure functional aspects. Plus, I wasn't looking forward to keeping JVM and
Scala dependencies up to date for this script.

I'm sharing this for anyone interested in the _PESOS_ philosophy. You'll probably need to tweak the source a bit
for your own use case and Hugo theme, but the code will get you most of the way there. The script will only
archive public toots or self-replies.

## Steps

### Request Archive

The first step is to [request a copy of your Mastodon archive](https://docs.joinmastodon.org/user/moving/#export). This
will take a bit of time to create. I received an email from [Hachyderm](https://hachyderm.io/@mweagle) when
the archive was ready for download.

### Expand Archive

After the archive is created, download it locally and expand it somewhere. The script expects the ZIP archive to be
previously expanded. The _outbox.json_ file in the archive is an [Activity Streams](https://www.w3.org/TR/activitystreams-core/)
serialization and its contents drive the archive process.

### Modify Templates

Toots Markdown frontmatter and content are rendered using go's [text/template](https://pkg.go.dev/text/template)
package. The _Execute_ parameter dictionary is:

{{< highlight go >}}

templateParamMap := map[string]interface{}{
    "ExecutionTime": nowTime,
    "Toot":          eachItem,
}

{{< /highlight >}}

The Toot object is an `ActivityEntry` struct with fields extracted from deserializing the _outbox.json_ 
`OrderedItems` contents.

The templates make assumptions about my site and the [theme](https://github.com/CaiJimmy/hugo-theme-stack)
I'm using. For example, the frontmatter template includes a hardcoded reference to a site asset:

{{< highlight go >}}

---
title: "Mastodon - {{ .Toot.Published }}"
subtitle: ""
canonical: {{ .Toot.Object.ID }}
description:
image: "/images/mastodon.png"
...

{{< /highlight >}}

Change both the `TEMPLATE_TOOT_FRONTMATTER` and `TEMPLATE_TOOT` to conform with your
sites' active theme.

### Modify Identity Constants

There are three constants that are hardcoded as well. These values are used to determine
toot visibility and reply threads.

{{< highlight go >}}

var HOST = "hachyderm.io"
var USER = "mweagle"

{{< /highlight >}}

If there's interest I will update these to command line arguments.

### Run

The script accepts two parameters:

- __input__: Path to the root of the expanded archive. This is the directory housing _outbox.json_
- __output__: Path to which the toots should be rendered
  - ‼️ Note that __all__ contents in this directory __will be deleted__ during rendering

Example:

{{< highlight shell >}}

go run mastodon_to_hugo.go \
    --input "~/Downloads/mastodon-archive" \
    --output "./blog/content/mastodon"

{{< /highlight >}}

### Output

Assuming all goes well, the script will emit some summary stats:

{{< highlight shell >}}

time=2024-03-09T23:14:11.187-08:00 level=INFO msg="Welcome to Hugodon!"
time=2024-03-09T23:14:11.267-08:00 level=INFO msg="Toots filtered" totalCount=1498 filteredCount=963
time=2024-03-09T23:14:11.383-08:00 level=INFO msg="Deleting existing directory contents" path=./blog/content/mastodon
time=2024-03-09T23:14:11.608-08:00 level=INFO msg="Publishing statistics" totalTootCount=1498 renderedTootCount=963 filteredTootCount=535 replyThreadCount=224 mediaFilesCount=43
time=2024-03-09T23:14:11.616-08:00 level=INFO msg="Toot replication complete"

{{< /highlight >}}

and a a series of
[Page Bundles](https://gohugo.io/content-management/page-bundles/)
with your toots. Markdown files will be created in the **output**
directory using the following template:

`{year}/{monthNumber}/{lastPathComponentOfTootObjectID}`

For example, using the hardcoded templates, the source entry:

{{< highlight json >}}

{
    "id": "https://hachyderm.io/users/mweagle/statuses/111860128005683093/activity",
    "type": "Create",
    "actor": "https://hachyderm.io/users/mweagle",
    "published": "2024-02-02T05:01:37Z",
    "to": [
        "https://www.w3.org/ns/activitystreams#Public"
    ],
    "cc": [
        "https://hachyderm.io/users/mweagle/followers"
    ],
    "object": {
        "id": "https://hachyderm.io/users/mweagle/statuses/111860128005683093",
        "type": "Note",
        "summary": null,
        "inReplyTo": null,
        "published": "2024-02-02T05:01:37Z",

{{</highlight>}}

is rendered to [mastodon/2024/02/111860128005683093](https://mweagle.net/mastodon/2024/02/111860128005683093/).

### Notes

- All self-reply threads are appended to the primary Toot's markdown file
- Toot attachments are copied to the page bundle directory. Attachments can be referenced in the
toot template via the `ActivityObjectAttachment.BaseFilename` field value
- ActivityFeed tags include a leading `#` character. This is stripped from the `ActivityObjectTag.Name` field
- Only `Hashtag` tag types are deserialized

## Source

Get the source [here](https://github.com/mweagle/blog/blob/main/assets/scripts/mastodon_to_hugo.go).
