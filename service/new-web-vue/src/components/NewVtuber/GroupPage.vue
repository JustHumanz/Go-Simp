<script setup>
import AmeLoading from "../AmeComp/AmeLoading.vue"
</script>
<template>
  <AmeLoading v-if="!groups.length" class="!h-screen" />
  <h4 class="mt-2 text-lg font-semibold" v-if="groups.length">
    Select your group...
  </h4>
  <ul class="group-list">
    <li class="group-list-item" v-for="group in groups">
      <a
        class="group-list-item__link"
        href="#"
        :data-id="group.ID"
        onclick="return false"
        @click="setGroup"
      >
        <font-awesome-icon
          class="fa-fw group-list-item__svg"
          icon="users"
          v-if="!group.ID"
        />
        <font-awesome-icon
          class="fa-fw group-list-item__svg"
          icon="user"
          v-else-if="group.ID && group.ID == 10"
        />
        <img
          draggable="false"
          v-else-if="group.ID && group.ID != 10"
          :src="group.GroupIcon"
          :alt="group.GroupName"
          class="group-list-item__img"
        />
        <span class="group-list-item__span">
          {{
            (
              group.GroupName.charAt(0).toUpperCase() + group.GroupName.slice(1)
            ).replace("_", " ")
          }}
        </span>
      </a>
    </li>
    <li class="group-list-item">
      <a
        href="#"
        class="group-list-item__link"
        onclick="return false"
        @click="setGroup"
        v-if="groups.length"
      >
        <font-awesome-icon
          class="fa-fw group-list-item__svg fa-fw"
          icon="circle-plus"
        />
        <span class="group-list-item__span">Request new group</span>
      </a>
    </li>
  </ul>
</template>
<script>
export default {
  data() {
    return {
      toggle: false,
    }
  },
  props: {
    groups: {
      type: Array,
      default: [],
    },
  },
  emits: ["group"],
  methods: {
    async setGroup(e) {
      const id = e.target.closest(".group-list-item__link").dataset.id
      const group = this.groups.find((group) => group.ID == id)
      if (!group) this.$emit("group", null)
      else this.$emit("group", group)
    },
  },
}
</script>
<style lang="scss" scoped>
.group {
  &-list {
    @apply grid gap-2 pt-2;
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));

    &-item {
      @apply overflow-hidden rounded-md bg-blue-400 dark:bg-slate-700;

      &__link {
        @apply flex h-full items-center space-x-1 px-2 py-1 font-semibold text-white transition-all duration-200 ease-in-out hover:bg-black/10;
      }

      &__img {
        @apply w-6 object-contain drop-shadow-md;
      }

      &__svg {
        @apply w-6;
      }

      &__span {
        @apply text-sm font-semibold text-white;
      }
    }
  }
}
</style>
