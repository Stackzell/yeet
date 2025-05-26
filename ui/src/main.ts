import { createApp } from "vue";
import "./style.css";
import { VueQueryPlugin } from "@tanstack/vue-query";
import App from "./App.vue";
import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [{ path: "/", component: () => import("./pages/Home.vue") }],
});

const app = createApp(App);

app.use(router);
app.use(VueQueryPlugin);
app.mount("#app");
