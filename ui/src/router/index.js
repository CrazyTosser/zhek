import Vue from "vue";
import VueRouter from "vue-router";
import Controller from "@/views/Controller.vue";
import Project from "@/views/Project.vue";
import Param from "@/views/Param";
import Test from "@/views/test";

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
    component: Project,
  },
  {
    path: "/param",
    name: "Param",
    component: Param,
  },
  {
    path: "/test",
    name: "Test",
    component: Test,
  },
];

const router = new VueRouter({
  routes,
});

export default router;
