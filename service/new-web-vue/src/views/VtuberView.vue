<template>
    <div class="filter">
      <div class="group-menu">
        <a href="#">All Vtubers</a>
        <ul class="groups">
          <li class="group"><router-link to="">All Vtubers</router-link></li>
          
        </ul>
      </div>

      <div class="flag-menu">

      </div>
    </div>
  <section class="vtuber-page">
  </section>
</template>

<style lang="scss">
  .filter {
    @apply bg-blue-400 fixed top-16 w-screen py-3 px-3 md:px-[15%] lg:px-[17%] flex items-center;

    .group-menu {
      @apply inline-block;

      & > a {
        @apply bg-slate-200 py-1 px-2 rounded-md;
      } 

      &:focus-within .groups {
        @apply visible;
      }

      .groups {
        @apply invisible bg-zinc-800 rounded-md absolute mt-1;
      }
    }
  }
.vtuber-page {
  @apply w-full md:w-[80%] lg:w-[75%] m-auto pt-[4.2rem];

}
</style>
