<script setup>
import "./style.scss"
</script>
<template>
  <h3 class="content-title">Roles and Taging</h3>
  <h4 class="content-subtitle">Roles</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>Roles <span>@role-name</span></p>
  <p class="content-text">
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">/role-info <span>role-name</span></p>
  <p>Show following group or member vtuber in your roles <br><br><b>Example 1:</b></p>
  <p class="content-code">vtbot>Roles <span>@Hololive</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">/role-info <span>role-name: @nijisanji</span></p>
  <h4 class="content-subtitle">Taging Myself</h4>
  <p class="content-code">vtbot>Tag me <span>group/member</span></p>
</template>
