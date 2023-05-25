---
title: "Resilience is a Verb"
date: 2023-05-25T05:51:53-08:00
lastmod: 2023-05-25T05:51:53-08:00
draft: false
subtitle:
description:
tags: ["resilience", "adaptive capacity", "paper"]
image:

ShowToc: true
---

## Source

Woods, David. (2018). Resilience is a Verb.
https://www.researchgate.net/publication/329035477

## TL;DR

This is a very readable paper on resilience engineering. At a high level, Woods' proposal is that resilience is better thought of in terms of a system's potential capacity, rather than a state to be achieved. Adaptive capacity is "the potential for adjusting patterns of activities to handle future changes in the kinds of  events, opportunities and disruptions experienced."

## Summary

As laid out in [“Resilience Engineering : Concepts and Precepts](https://www.researchgate.net/publication/50232053), resilience is about what a system can do, including the capacity to: 1. anticipate to mitigate; 2. synchronize to reduce coordination friction/misalignment; 3. promptly respond; 4. proactively learn.

Adaptive capacity is a potential held in reserve rather than an essentialist system property. This implies some level of system self-awareness so that there are explicit junctures to evaluate whether the current course is appropriate (pivot, proceed, pause). Adaptive capacity is an option, not a mandate for continual change:

> Adaptive capacity does not mean a system is constantly changing what it has planned or does so all the time, but rather that the system has some ability to recognize when it's adequate to continue the plan, to continue to work in the usual way, and when it is not adequate to continue on, given the demands, changes and context ongoing or upcoming.

Adaptive capacity is a latent system property and all systems are messy and running degraded. SNAFUs (_Situation Normal, All Fouled Up_) provide opportunities to surface where the current system's capacity is limited and could or should be re-evaluated. Similar to the "Learning from near-misses" perspective.

There's a cycle of 1. systems growing in scale, tempo, and complexity (interconnectedness, Law of Stretched Systems); 2. actors in the system exploiting those changes for new purposes; 3. creating new SNAFUs which move the system to a new normal and move things back to 1.

Rather than viewing SNAFUs as intrinsic and opportunities to improve resilience, organizations often "rationalize this core finding away on grounds of rarity, prevention, compliance.". It's more appealing to view SNAFUs as uncommon, obsolete due to prior improvements, or caused by a lack of human compliance. The emphasis on compliance and work-to-{plan, role, rule} pushes the adaptive capacity underground.

Complex systems improve over time and also catastrophically fail when encountering new circumstances. What to do? Continually adapt!

> Organizational systems succeed despite the basic limits of automata and plans in a complex, interdependent and changing environment because responsible people adapt to make the system work despite its design—SNAFU Catching.

There are four capabilities needed for continual adaptation: 1. initiative, favoring local cognition; 2. reciprocity across people and teams to load balance; 3. exploiting surprises to drive learnings; and 4. effective learning reviews.

> To be proactive in learning about resilience shifts the focus: study how systems work well usually despite difficulties, limited resources, trade-offs, and surprises—SNAFU Catching.

## Notable Quotes

> As organizations focus on making systems work faster, better, and cheaper, they develop new plans embodied in procedures, automation, policies, and forcing functions. These plans are seen as effective since they represent improvements relative to how the system worked previously. When surprising results occur, the organization interprets the surprises as deviations—erratic people were unable to work to plan, to work to their role within the plan, and to work to the rules prescribed for their role. The countermeasures become more stringent pressures to work-to-plan, work-to-role and work-to-rule.
