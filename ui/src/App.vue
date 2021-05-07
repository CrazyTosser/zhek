<template>
  <div id="app">
    <b-navbar toggleable="lg" type="dark" variant="dark">
      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item
            ><router-link to="/contr">Контроллеры</router-link>
          </b-nav-item>
          <b-nav-item><router-link to="/proj">Проекты</router-link></b-nav-item>
          <b-nav-item
            ><router-link to="/param">Параметры</router-link></b-nav-item
          >
          <b-nav-item><router-link to="/test">Тест</router-link></b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <router-view />
  </div>
</template>

<style>
a {
  color: white !important;
}
html,
body,
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 100%;
  height: 100%;
}
</style>
