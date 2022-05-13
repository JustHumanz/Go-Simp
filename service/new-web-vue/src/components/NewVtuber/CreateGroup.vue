<script setup>
import GroupPlatform from "./GroupPlatform.vue"
</script>

<template>
  <h1 class="title-req-group">Request new Group</h1>
  <form @submit="sendGroup" class="form-req-group">
    <div class="form-group">
      <label for="group-name">Group Name</label>
      <input
        type="text"
        id="group-name"
        name="group-name"
        class="form-control"
        autocomplete="off"
        placeholder="Group Name"
      />
      <small class="error"></small>
    </div>
    <div class="form-group">
      <label for="group-icon">URL Icon</label>
      <input
        type="text"
        id="group-icon"
        name="group-icon"
        class="form-control"
        autocomplete="off"
        placeholder="URL Icon"
      />
      <small class="description"
        >(Optional) It is recommended to fill this in accordingly to make it
        easier to find the icon.</small
      >
      <small class="error"></small>
    </div>
    <div class="form-group">
      <button type="button" name="add-platform" @click="addPlatform">
        Add Platform
      </button>
    </div>
    <div class="platforms">
      <group-platform
        v-for="platform in platforms"
        :id="platform.id"
        :set="platform.set"
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
import CreateGroup_script from "./CreateGroup_script.js"

export default { ...CreateGroup_script }
</script>

<style lang="scss" scoped>
.title-req-group {
  @apply my-2 text-lg font-semibold;
}

.form-req-group {
  @apply flex w-full flex-col space-y-2 md:w-3/6;
}

.form-group {
  @apply flex flex-col md:ml-2;

  &.has-error {
    input {
      @apply bg-red-200;
    }

    .description {
      @apply hidden;
    }

    .error {
      @apply block;
    }
  }

  input {
    @apply my-1 -translate-y-0.5 rounded-lg bg-slate-200 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
  }

  small {
    @apply text-xs;
  }

  .description {
    @apply text-gray-600;
  }

  .error {
    @apply hidden text-red-500;
  }

  button {
    @apply -translate-y-0.5 rounded-lg bg-blue-400 p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm disabled:translate-y-0.5 disabled:opacity-70 disabled:shadow-none;
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

<style lang="scss">
.platforms {
  .platform {
    @apply flex flex-col md:ml-2;

    .delete-btn {
      @apply mt-2 -translate-y-0.5 rounded-lg bg-red-400 px-2 py-1 text-sm font-semibold text-white shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm disabled:translate-y-0.5 disabled:opacity-70 disabled:shadow-none;
    }

    label {
      @apply text-lg;
    }

    select {
      @apply my-1 -translate-y-0.5 rounded-lg bg-slate-200 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
    }

    input {
      @apply my-1 -translate-y-0.5 rounded-lg bg-slate-200 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
    }

    small {
      @apply text-xs;
    }

    .description {
      @apply text-gray-600;
    }

    .error {
      @apply text-red-500;
    }
  }
}
</style>
