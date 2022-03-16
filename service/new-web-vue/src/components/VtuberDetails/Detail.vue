<template>
  <div class="detail-vtuber">
    <div class="detail-vtuber-item">
      <span class="detail-vtuber-item__title"> Name </span>
      <span class="detail-vtuber-item__value">
        {{ vtuber.EnName }}
      </span>
    </div>

    <div class="detail-vtuber-item" v-if="vtuber.JpName">
      <span class="detail-vtuber-item__title"> Japanese Name </span>
      <span class="detail-vtuber-item__value">
        {{ vtuber.JpName }}
      </span>
    </div>

    <div class="detail-vtuber-item">
      <span class="detail-vtuber-item__title"> Nickname </span>
      <span class="detail-vtuber-item__value">
        {{ vtuber.NickName.toLowerCase() }}
      </span>
    </div>

    <div class="detail-vtuber-item">
      <span class="detail-vtuber-item__title"> Group/Agency </span>
      <router-link
        :to="`/vtubers/${vtuber.Group.ID}`"
        class="detail-vtuber-item__value"
      >
        <img
          :src="vtuber.Group.IconURL"
          :alt="GroupName"
          class="detail-vtuber-item__icon"
          v-if="vtuber.Group.ID !== 10"
        />
        {{ GroupName }}
      </router-link>
    </div>

    <div class="detail-vtuber-item">
      <span class="detail-vtuber-item__title"> Region </span>
      <router-link
        :to="`/vtubers?reg=${vtuber.Region}`"
        class="detail-vtuber-item__value"
      >
        <img
          :src="`/assets/flags/${vtuber.Regions.flagCode}.svg`"
          :alt="vtuber.Regions.name"
          class="detail-vtuber-item__icon"
        />
        {{ vtuber.Regions.name }}
      </router-link>
    </div>

    <div class="detail-vtuber-item">
      <span class="detail-vtuber-item__title"> Status </span>
      <span class="detail-vtuber-item__value">
        {{ vtuber.Status }}
      </span>
    </div>

    <div class="detail-vtuber-item" v-if="vtuber.Fanbase">
      <span class="detail-vtuber-item__title"> Fanbase </span>
      <span class="detail-vtuber-item__value">
        {{ vtuber.Fanbase }}
      </span>
      </div>
  </div>
  <hr class="m-2" />
</template>

<script>
export default {
  props: {
    vtuber: Object,
  },
  computed: {
    GroupName() {
      return (
        this.vtuber.Group.GroupName.charAt(0).toUpperCase() +
        this.vtuber.Group.GroupName.slice(1).replace("_", " ")
      )
    },
  },
}
</script>

<style lang="scss" scoped>
a.detail-vtuber-item__value {
  @apply hover:text-gray-600 dark:hover:text-gray-300;
}

.detail-vtuber {
  @apply flex flex-wrap p-2 justify-center;

  &-item {
    @apply inline-flex flex-col px-2;

    &__title {
      @apply text-sm font-semibold;
    }

    &__value {
      @apply font-light;
    }

    &__icon {
      @apply mr-[2px] -mt-1 w-5 object-contain inline-block drop-shadow-sm;
    }
  }
}
</style>
