import Vue from "vue";
import VueRouter from "vue-router";

import Index from "./Index.vue";
import Register from "./Register.vue";
import Login from "./Login.vue";

Vue.use(VueRouter);

export default new VueRouter({
  routes: [
    {
      path: "/",
      component: Index
    },
    {
      path: "/register",
      component: Register
    },
    {
      path: "/login",
      component: Login
    }
  ]
});
