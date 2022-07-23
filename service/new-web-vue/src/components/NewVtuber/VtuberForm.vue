<script setup>
import VtuberPlatform from "./VtuberPlatform.vue"
</script>

<template>
  <h1 class="title-req-vtuber">
    Request new Vtuber from {{ group.GroupName }}
  </h1>
  <form @submit="requestVtuber" class="form-req-vtuber">
    <div class="form-group">
      <button type="button" name="add-vtuber" @click="addVtuber">
        Add Vtuber
      </button>
    </div>
    <div class="vtubers">
      <vtuber-platform
        v-for="vtuber in vtubers"
        :id="vtuber.id"
        :nicknames="nicknames"
        @error="checkError"
        @delete="deletePlatform"
      />
    </div>
    <!-- Submit and back btn -->
    <div class="form-group-btn">
      <button type="submit" class="submit" disabled>Next</button>
      <button type="button" class="back" @click="$emit(`back`, true)">
        Back
      </button>
    </div>
  </form>
</template>

<script>
import VtuberForm from "./VtuberForm_script.js"

export default { ...VtuberForm }
</script>

<style lang="scss" scoped>
.title-req-vtuber {
  @apply my-2 text-lg font-semibold;
}

.form-req-vtuber {
  @apply flex w-full flex-col space-y-2 md:w-3/6;
}

.form-group {
  @apply flex flex-col md:ml-2;

  button {
    @apply -translate-y-0.5 rounded-lg bg-blue-400 p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm disabled:translate-y-0.5 disabled:opacity-70 disabled:shadow-none dark:bg-slate-700;
  }
}

.form-group-btn {
  @apply grid grid-rows-2 gap-2 pb-4 md:grid-cols-2 md:grid-rows-none;

  button {
    @apply -translate-y-0.5 rounded-lg p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out first:bg-green-500 last:bg-red-600 hover:translate-y-0 hover:shadow-sm disabled:translate-y-0 disabled:opacity-70 disabled:shadow-none;
  }

  .platforms {
    @apply flex flex-col gap-2;
  }
}
</style>
