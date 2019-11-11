import * as monaco from "monaco-editor";

export default function() {
  self.MonacoEnvironment = {
    getWorkerUrl: function(moduleId, label) {
      if (label === "json") {
        return "./json.worker.js";
      }
      if (label === "css") {
        return "./css.worker.js";
      }
      if (label === "html") {
        return "./html.worker.js";
      }
      if (label === "typescript" || label === "javascript") {
        return "./ts.worker.js";
      }
      return "./editor.worker.js";
    }
  };
  monaco.languages.register({
    id: "satysfi"
  });
  monaco.languages.setMonarchTokensProvider("satysfi", {
    keywords: [
      "not",
      "mod",
      "if",
      "then",
      "else",
      "let",
      "let-rec",
      "and",
      "in",
      "fun",
      "true",
      "false",
      "before",
      "while",
      "do",
      "let-mutable",
      "match",
      "with",
      "when",
      "as",
      "type",
      "of",
      "module",
      "struct",
      "sig",
      "val",
      "end",
      "direct",
      "constraint",
      "let-inline",
      "let-block",
      "let-math",
      "controls",
      "cycle",
      "inline-cmd",
      "block-cmd",
      "math-cmd",
      "command",
      "open"
    ],
    tokenizer: {
      root: [
        [/%.+$/, "comment"],
        [/@(require|import)/, "preamble"],
        [/\+[a-zA-Z][0-9a-zA-Z-]*/, "blockcommand"],
        [/\\[a-zA-Z][0-9a-zA-Z-]*/, "inlinecommand"],
        [
          /[a-z][0-9a-zA-Z-]*/,
          {
            cases: {
              "@keywords": "keyword",
              "@default": "identifier"
            }
          }
        ],
        [/`/, "string", "@string"],
        [/\*/, "item"]
      ],
      string: [[/[^`]/, "string"], [/`/, "string", "@pop"]]
    }
  });

  monaco.editor.defineTheme("satysfier", {
    base: "vs",
    inherit: true,
    rules: [
      { token: "preamble", foreground: "ff3333" },
      { token: "item", fontStyle: "bold" },
      { token: "blockcommand", foreground: "4000ff" },
      { token: "inlinecommand", foreground: "c14800" }
    ]
  });
}
