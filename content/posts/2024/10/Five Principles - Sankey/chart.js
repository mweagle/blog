(function () {
  var sankeyData = {
    nodes: [
      { name: "FivePrinciples" },
      { name: "People Make Mistakes" },
      { name: "Errors are normal" },
      {
        name: "Since error is a\nnormal state of existence,\nerror is never causal",
      },
      {
        name: "Error is not the opposite of success.\nError is a part of being successful",
      },
      { name: "Error exists in success\nas well as failure" },
      {
        name: "Errors are not choices.\nError only becomes\na choice in retrospect.",
      },
      {
        name: "You can’t remove error\nso you must defend against\nthe inevitability of error",
      },
      { name: "Good systems\nbuild in error tolerance" },
      { name: "Knowing errors will happen\nis a good thing" },
      {
        name: "An error without significance\nconsequence is the closest\nthing to a leading indicator data\nyou will find",
      },
      { name: "Learning And Improving" },
      {
        name: "Organizations have two choices\nwhen responding to failure:\nto learn and improve\nor to blame and punish",
      },
      {
        name: "Learning is a strategic and operational\nchoice towards improvement",
      },
      { name: "Learning is a deliberate\nimprovement strategy" },
      { name: "Knowing how work is done is difficult" },
      {
        name: "Workers are the experts,\nthe profound users of the work process",
      },
      { name: "Workers always complete\nthe process design" },
      {
        name: "Defenses are placed in systems,\ntested in systems,\nand strengthened in systems\nby learning how successful\nwork is done",
      },
      { name: "Context Drives Behavior" },
      {
        name: "Workers do what they do for a reason.\nAnd the reason makes sense\nto the worker given the context",
      },
      {
        name: "Complex systems don’t lend themselves\nto traditional metrics",
      },
      {
        name: "Local rationale is information\nto be discovered, not\nto be weaponized",
      },
      {
        name: "The environment in which work occurs\nmainly determines workers\nbehavior and actions",
      },
      {
        name: "Individual behavior is influenced\nby the organizational processes and values",
      },
      { name: "How You Respond" },
      { name: "You have two choices:\ngetting better or getting even" },
      {
        name: "You can blame and punish,\nor you can learn and improve,\nbut you can’t do both",
      },
      { name: "You create the feedback\nsystems you have" },
      {
        name: "Managers shape how\norganizations learn by\ntheir reaction to failure",
      },
      {
        name: "Every aspect of improvement\nis contingent on leadership’s\ndeliberate decision to get better",
      },
      { name: "People are watching you" },
      { name: "Blame Fixes Nothing" },
      {
        name: "Blame is emotionally important,\nbut not operationally important",
      },
      { name: "Blame makes error a choice in retrospect" },
      {
        name: "Blame takes up emotional and intellectual\nspace with little added value",
      },
      { name: "Blame misdirects resources and strategies" },
      { name: "Blame is the opposite of encouragement" },
    ],
    links: [
      {
        source: "FivePrinciples",
        target: "People Make Mistakes",
        value: 9,
      },
      { source: "People Make Mistakes", target: "Errors are normal", value: 1 },
      {
        source: "People Make Mistakes",
        target:
          "Since error is a\nnormal state of existence,\nerror is never causal",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target:
          "Error is not the opposite of success.\nError is a part of being successful",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target: "Error exists in success\nas well as failure",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target:
          "Errors are not choices.\nError only becomes\na choice in retrospect.",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target:
          "You can’t remove error\nso you must defend against\nthe inevitability of error",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target: "Good systems\nbuild in error tolerance",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target: "Knowing errors will happen\nis a good thing",
        value: 1,
      },
      {
        source: "People Make Mistakes",
        target:
          "An error without significance\nconsequence is the closest\nthing to a leading indicator data\nyou will find",
        value: 1,
      },
      {
        source: "FivePrinciples",
        target: "Learning And Improving",
        value: 7,
      },
      {
        source: "Learning And Improving",
        target:
          "Organizations have two choices\nwhen responding to failure:\nto learn and improve\nor to blame and punish",
        value: 1,
      },
      {
        source: "Learning And Improving",
        target:
          "Learning is a strategic and operational\nchoice towards improvement",
        value: 1,
      },
      {
        source: "Learning And Improving",
        target: "Learning is a deliberate\nimprovement strategy",
        value: 1,
      },
      {
        source: "Learning And Improving",
        target: "Knowing how work is done is difficult",
        value: 1,
      },
      {
        source: "Learning And Improving",
        target:
          "Workers are the experts,\nthe profound users of the work process",
        value: 1,
      },
      {
        source: "Learning And Improving",
        target: "Workers always complete\nthe process design",
        value: 1,
      },
      {
        source: "Learning And Improving",
        target:
          "Defenses are placed in systems,\ntested in systems,\nand strengthened in systems\nby learning how successful\nwork is done",
        value: 1,
      },
      {
        source: "FivePrinciples",
        target: "Context Drives Behavior",
        value: 4,
      },
      {
        source: "Context Drives Behavior",
        target:
          "Workers do what they do for a reason.\nAnd the reason makes sense\nto the worker given the context",
        value: 1,
      },
      {
        source: "Context Drives Behavior",
        target: "Complex systems don’t lend themselves\nto traditional metrics",
        value: 1,
      },
      {
        source: "Context Drives Behavior",
        target:
          "Local rationale is information to\nbe discovered,\nno to \ne weaponized",
        value: 1,
      },
      {
        source: "Context Drives Behavior",
        target:
          "The environment in which work occurs\nmainly determines workers\nbehavior and actions",
        value: 1,
      },
      {
        source: "Context Drives Behavior",
        target:
          "Individual behavior is influenced\nby the organizational processes and values",
        value: 1,
      },
      {
        source: "FivePrinciples",
        target: "How You Respond",
        value: 6,
      },
      {
        source: "How You Respond",
        target: "You have two choices:\ngetting better or getting even",
        value: 1,
      },
      {
        source: "How You Respond",
        target:
          "You can blame and punish,\nor you can learn and improve,\nbut you can’t do both",
        value: 1,
      },
      {
        source: "How You Respond",
        target: "You create the feedback\nsystems you have",
        value: 1,
      },
      {
        source: "How You Respond",
        target:
          "Managers shape how\norganizations learn by\ntheir reaction to failure",
        value: 1,
      },
      {
        source: "How You Respond",
        target:
          "Every aspect of improvement\nis contingent on leadership’s\ndeliberate decision to get better",
        value: 1,
      },
      {
        source: "How You Respond",
        target: "People are watching you",
        value: 1,
      },
      {
        source: "FivePrinciples",
        target: "Blame Fixes Nothing",
        value: 4,
      },
      {
        source: "Blame Fixes Nothing",
        target:
          "Blame is emotionally important,\nbut not operationally important",
        value: 1,
      },
      {
        source: "Blame Fixes Nothing",
        target: "Blame makes error a choice in retrospect",
        value: 1,
      },
      {
        source: "Blame Fixes Nothing",
        target:
          "Blame takes up emotional and intellectual\nspace with little added value",
        value: 1,
      },
      {
        source: "Blame Fixes Nothing",
        target: "Blame misdirects resources and strategies",
        value: 1,
      },
      {
        source: "Blame Fixes Nothing",
        target: "Blame is the opposite of encouragement, value: 1",
      },
    ],
  };

  var foo = {
    title: {
      text: "Sankey Test",
    },
    nodeGap: 128,
    nodeWidth: 64,
    layoutIterations: 32,
    tooltip: {
      trigger: "item",
      triggerOn: "mousemove",
      formatter: "{b}",
    },
    series: {
      type: "sankey",
      layout: "none",
      emphasis: {
        focus: "adjacency",
      },
      lineStyle: {
        color: "gradient",
        curveness: 0.8,
      },
      data: sankeyData.nodes,
      links: sankeyData.links,
    },
  };
  return foo;
})();
