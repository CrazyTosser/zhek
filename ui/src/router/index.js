import Vue from "vue";
import VueRouter from "vue-router";
import Controller from "@/views/Controller.vue";
import Project from "@/views/Project.vue";
import Param from "@/views/Param";

Vue.use(VueRouter);

const routes = [
  {
    path: "/contr",
    name: "Controller",
    component: Controller,
  },
  {
    path: "/proj",
    name: "Project",
    component: Project
  },
  {
    path: "/param",
    name: "Param",
    component: Param
  }
];

const router = new VueRouter({
  routes,
});

export default router;
