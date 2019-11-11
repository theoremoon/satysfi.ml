import Vue from "vue";
import App from "./src/App.vue";
import router from "./src/router.js";
import store from "./src/store.js";
import initEditor from "./src/editor.js";

initEditor();

let vue = new Vue({
  el: "#app",
  router: router,
  store: store,
  render: h => h(App)
});
