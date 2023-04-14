---
title: "C4 PlantUML Themes"
author: "Matt Weagle"
date: 2023-04-10T08:57:31-08:00
lastmod: 2023-04-10T08:57:31-08:00
tags: ["architecture", "visualization"]
ShowToc: true
description: "Plant UML Themes"
draft: false
subtitle: "Using ColorBrewer palettes for C4 architecture diagrams"

image: ""

---

## Hugo Shortcode with C4 Plant UML themes

References:
* https://paul.dugas.cc/post/plantuml-shortcode/
* https://colorbrewer2.org
* https://seaborn.pydata.org/generated/seaborn.color_palette.html


{{< plantuml id="eg" source="./puml/sample.puml" >}}

<!--[![my image](./puml/resources/colorbrewer/palettes/div-BrBG-9.png)](./puml/resources/colorbrewer/palettes/div-BrBG-9.puml)-->

### All C4 Plant UML Themes

{{< pumlthemes id="allthemes"
  path="content/posts/c4pumlthemes/puml/resources/colorbrewer/palettes"
  sitepath="/posts/c4pumlthemes/puml/resources/colorbrewer/palettes/">}}
