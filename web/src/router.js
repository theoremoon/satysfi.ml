import Vue from "vue";
import VueRouter from "vue-router";

import Index from "./Index.vue";

Vue.use(VueRouter);
export default new VueRouter({
  routes: [
    {
      path: "/",
      component: Index
    },
    {
      path: "/project/:id",
      component: Index
    }
  ]
});
