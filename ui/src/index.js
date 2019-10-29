import Vue from "vue";
import App from "./App.vue";
import router from "./router.js";

let vue = new Vue({
  el: "#app",
  router,
  render: h => h(App)
});
