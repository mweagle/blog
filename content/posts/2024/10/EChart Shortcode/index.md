---
title: "Hugo Apache ECharts shortcode"
subtitle: "Extending Hugo with custom javascript"
description:
draft: false
date: 2024-10-25T22:07:34-07:00
lastmod: 2024-10-25T22:07:34-07:00
image: "pankaj-patel-u2Ru4QBXA5Q-unsplash.jpg"
toc: false
tags: [""]
categories: []
mermaid: true
echarts: true
cover:
    image: "pankaj-patel-u2Ru4QBXA5Q-unsplash.jpg"
    alt: ""
    caption: ""
    relative: true 
    linkFullImages: false
---

<!--more-->

This is mostly documentation for myself summarizing what I did to add a `<echarts>` shortcode to Hugo that allows
me to support [Apache Echarts](https://echarts.apache.org/en/index.html) on this blog. I initially thought I would use this to
add a [Treemap](https://echarts.apache.org/examples/en/index.html#chart-type-treemap) visualization of
[Todd Conklin's Five Principles of Human Performance](http://localhost:1313/posts/2024/10/five-principles-sankey/). These steps are very similar to Navendu Pottekkat's [Adding Diagrams to Your Hugo Blog With Mermaid](https://navendu.me/posts/adding-diagrams-to-your-hugo-blog-with-mermaid/).

Some of these steps are depend on the [hugo-theme-stack](https://github.com/CaiJimmy/hugo-theme-stack) partials. They
will likely need to be changed to support other themes.

## Steps

### Include the ECharts source

Download the ECharts source from the [download page](https://echarts.apache.org/en/download.html). Expand the archive and ensure the contents exist in the _static/echarts_ directory.

### Configure the theme params

Configure the light and dark echart themes via the site's _config.yaml_. The partial renderers
will use these basenames when resolving theme paths in the _static/echarts/theme_ directory.

```yaml
params:
  echarts:
    theme:
      light: infographic
      dark: dark-digerati
```

### Header Partial

Upsert the _layouts/partials/head/custom.html_ partial so that it conditionally includes the
ECharts partial render. Posts will declare whether they need EChart support by frontmatter. This
eliminates downloading unnecessary files for the majority of the pages.

```html
<!-- Add echarts js file -->
{{ if (.Params.echarts) }}
{{ partial "echarts.html" }}
{{ end }}
```

The _layouts/partials/head/custom.html_ conditionally includes the EChart specific partial
_layouts/partials/echarts.html_. The JavaScript in the `echarts` shortcode will check the Hugo
theme's [local storage value](https://github.com/CaiJimmy/hugo-theme-stack/blob/839fbd0ecb5bba381f721f31f5195fb6517fc260) 
when rendering to determine which theme to use.

```html
<!-- Include the ECharts file you just downloaded -->
<script src="/echarts/dist/echarts.js"></script>

<!-- Include the theme files -->
<script src="/echarts/theme/{{ if site.Params.echarts.theme.light }}{{ site.Params.echarts.theme.light }}{{ else }}default{{ end }}.js"></script>
<script src="/echarts/theme/{{ if site.Params.echarts.theme.dark }}{{ site.Params.echarts.theme.dark }}{{ else }}dark{{ end }}.js"></script>
```

This source optionally dereferences the theme names defined in the site's primary _config.yaml_. There's
also the [ECharts Theme Builder](https://echarts.apache.org/en/theme-builder.html) to create
a completely custom design.

### Usage

The first requirement is to declare that the post requires EChart support in the post's frontmatter at the root:

```yaml
echarts: true
```

To build a graph, the shortcode requires user data to create the DOM elements and initialize the chart:

- Required
  - ElementID
  - Chart Definition (either inline or as an external JS file)
- Optional
  - Dimensions (`width` and `height` optional)
  - Interactivity

#### Required

The shortcode requires an element ID to use for the DOM. This enables multiple
charts per page:

```html
{{</* echarts id="contentInline" */>}}
```

Chart content can be provided either inline or externally by a JavaScript file
that defines an immediately executed function. The options returned by either source
are provided to the [setOption](https://echarts.apache.org/en/api.html#echartsInstance.setOption)
call.

For inline data:

```html
{{</* echarts id="contentInline" */>}}
{
  title: {
      text: "Inline Chart Definition",
  },
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: [150, 230, 224, 218, 135, 147, 260],
      type: 'line'
    }
  ]
}
{{</* /echarts */>}}
```

Renders as (ignoring the optional `width` and `height` params):

{{< echarts id="contentInline" >}}
{
  title: {
      text: "Inline Chart Definition",
  },
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: [150, 230, 224, 218, 135, 147, 260],
      type: 'line'
    }
  ]
}
{{< /echarts >}}

I started with this approach, but missed VSCode's syntax highlighting. An alternative
version is to define an immediate JS function in an external file.

```html
{{</* echarts id="externalContent" srcChart="content/posts/2024/10/EChart Shortcode/line.js" */>}}
{{</* /echarts */>}}
```

Renders as:

{{< echarts id="externalContent" srcChart="content/posts/2024/10/EChart Shortcode/line.js" >}}
{{< /echarts >}}

#### Optional - Dimensions

The shortcode supports optional `width` (default=_1024px_) and `height` (default=_800px_) params:

```html
{{</* echarts id="externalContent" width="800px" height="600px" */>}}
{{</* /echarts */>}}
```

http://localhost:1313/posts/2024/10/echart-shortcode/posts/2024/10/EChart%20Shortcode/gauge.js

#### Optional - Interactivity

Some ECharts like the [gauge](https://echarts.apache.org/examples/en/editor.html?c=gauge-clock) require
event handlers to be configured. For these types of charts the entire configuration can be externalized into a 
`jsSource` param:

```html
{{</* echarts id="jsInteractive" width="400px" height="400px" jsSource="./gauge.js" */>}}
{{</* /echarts */>}}
```

The content of the `jsSource` file is an immediate JavaScript function like:

```js
(function() {
    var chartOps = {
        ...
    };
    return chartOps;
})
```

Renders as:

{{< echarts id="jsInteractive" width="400px" height="400px"jsSource="./gauge.js" >}}
{{< /echarts >}}

### Shortcode

The shortcode itself is responsible for creating a unique DOM element, constructing
the chart input from either the inline or external source, and then optionally loading
the interactive data. 

{{% embed "/layouts/shortcodes/echarts.html" %}}

## Attributions

Photo by <a href="https://unsplash.com/@pankajpatel?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Pankaj Patel</a> on <a href="https://unsplash.com/photos/turned-on-monitor-displaying-programming-language-u2Ru4QBXA5Q?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Unsplash</a>
  