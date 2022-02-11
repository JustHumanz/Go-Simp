<template>
    <div class="filter">
      <div class="menu-filters">
        <a href="#">All Vtubers</a>
        <ul class="groups">
          <li class="group"><router-link to="">All Vtubers</router-link></li>
          
        </ul>
      </div>

      <div class="menu-filters">
        <a href="#">All Regions</a>
        <ul class="groups">
          <li class="group"><router-link to="">All Vtubers</router-link></li>
          
        </ul>
      </div>

      <div class="menu-filters">
        <a href="#">More Filters</a>
        <ul class="groups">
          <li class="group"><input type="checkbox" name="" id="inactive"><label for="inactive">Inactive Vtuber</label></li>
          <li class="group"><input type="checkbox" name="" id="live"><label for="live">Live Now</label></li>
          
        </ul>
      </div>


      <div class="search">
        <input type="text" placeholder="Search..."/> 
      </div>

    </div>
  <section class="vtuber-page">
  </section>
</template>

<style lang="scss">
  .filter {
    @apply bg-blue-400 fixed top-16 w-screen py-3 px-3 md:px-[15%] lg:px-[17%] flex flex-wrap items-center;

    .menu-filters {
      @apply inline-block mx-1;

      & > a {
        @apply bg-slate-200 py-1 px-2 rounded-md;
      } 

      &:focus-within .groups {
        @apply visible;
      }

      .groups {
        @apply invisible bg-zinc-800 rounded-md absolute mt-1 z-10;
      }
    }

    .search {
      @apply inline-block mx-1;

      input {
        @apply bg-blue-200 py-1 px-2 rounded-md;
      }
    }
  }
.vtuber-page {
  @apply w-full md:w-[80%] lg:w-[75%] m-auto pt-[4.2rem];

}
</style>
