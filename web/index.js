import Vue from "vue";
import App from "./src/App.vue";
import router from "./src/router.js";
import * as monaco from "monaco-editor";

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

let vue = new Vue({
  el: "#app",
  router: router,
  render: h => h(App)
});
