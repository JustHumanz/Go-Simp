<template>
  <div id="app">
    <img src="./assets/logo.png">
      <div id="nav">
        <router-link to="/">Home</router-link> |
        <router-link to="/group">group</router-link>
      </div>
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
