<style scoped>
  @import "./style.css";
</style>
<template>
  <div style="min-height:100%; height:100px;min-width: 100%">
    <div id="particles-js" />
  </div>
</template>

<script>
import particlesJs from 'particles.js'
import particlesConfig from './particles.json'
export default {
  data() {
    return {}
  },
  mounted() {
    this.init()
  },
  methods: {
    init() {
      particlesJS('particles-js', particlesConfig)
      document.body.style.overflow = 'hidden'
    }
  }
}
</script>
