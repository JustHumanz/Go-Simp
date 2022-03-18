<script setup>
import AmeError from '../components/AmeComp/AmeError.vue'
</script>

<template>
<ame-error type="error"
    img="tothemon"
    :title="`Could not find page ${this.$route.path.replace(/^\//, '')}`"
    description="Actually I can make the page, but I can't because that's the rule"/>
</template>

<script>
export default {
  created() {
    document.title = 'Page not found - Vtbot'
  },
}
</script>
