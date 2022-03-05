<template>
  <header class="header">
    <div class="header-image">
      <img src="/src/assets/vtuber.png" draggable="false" alt="" class="header-image__img" />
    </div>
    <h3 class="header-title">Friend for get notification from vtuber</h3>
    <a href="" class="header-button">
      <font-awesome-icon :icon="['fab', 'discord']" class="fa-fw" />
      Add to Discord
    </a>
    <div class="header-links">
      <a href="" class="header-links__link">
        <font-awesome-icon :icon="['fab', 'discord']" class="fa-fw" />
        Discord Server
      </a>
      <span class="header-links__separator">|</span>
      <a href="" class="header-links__link">
        <font-awesome-icon :icon="['fab', 'github']" class="fa-fw" />
        Source Code
      </a>
    </div>
  </header>
  <section class="feature">
    <div class="feature-image">
      <img
        src="/src/assets/graphics/SVG/live_light.svg"
        alt=""
        draggable="false"
        class="img-light"
      />
      <img
        src="/src/assets/graphics/SVG/live_dark.svg"
        alt=""
        draggable="false"
        class="img-dark"
      />
    </div>
    <div class="feature-content">
      <h3 class="feature-content__title">Live Notification</h3>
      <p class="feature-content__description">
        Get any live notification over +400 vtubers.
      </p>
    </div>
  </section>
  <section class="feature reverse">
    <div class="feature-image">
      <img
        src="/src/assets/graphics/SVG/fanart_light.svg"
        alt=""
        draggable="false"
        class="img-light"
      />
      <img
        src="/src/assets/graphics/SVG/fanart_dark.svg"
        alt=""
        draggable="false"
        class="img-dark"
      />
    </div>
    <div class="feature-content">
      <h3 class="feature-content__title">Fan Art</h3>
      <p class="feature-content__description">
        Look a cool random fan art and NSFW every members.
      </p>
    </div>
  </section>
  <section class="feature">
    <div class="feature-image">
      <img
        src="/src/assets/graphics/SVG/member_light.svg"
        alt=""
        draggable="false"
        class="img-light"
      />
      <img
        src="/src/assets/graphics/SVG/member_dark.svg"
        alt=""
        draggable="false"
        class="img-dark"
      />
    </div>
    <div class="feature-content">
      <h3 class="feature-content__title">Milestone</h3>
      <p class="feature-content__description">
        Track your favorite vtuber's for target milestone.
      </p>
    </div>
  </section>
  <section class="ending">
    <div class="ending-content">
      <h3 class="ending-content__title">
        Your favorite vtuber cannot be here?
        </h3>
      <p class="ending-content__description">
        Request your favorite vtuber to be here.
      </p>
    </div>
    <a href="" class="ending-button">Request Here</a>
  </section>
</template>

<script>
// add discord and github in font awesome
import { library } from "@fortawesome/fontawesome-svg-core"
import { faDiscord, faGithub } from "@fortawesome/free-brands-svg-icons"

library.add(faDiscord, faGithub)

export default {
  async created() {
    this.$watch(
      () => this.$route.params,
      async () => {
        document.title = "Vtbot - A vtuber scheduler from Discord Server"
      },

      { immediate: true }
    )
  },
}
</script>

<style lang="scss" scoped>
.header {
  @apply bg-indigo-300 h-[90vh] px-12 flex flex-col items-center justify-center;

  &-image {
    @apply overflow-hidden h-[19vmax] flex items-start drop-shadow-sm;

    &__img {
      @apply w-[35vmax] object-contain;
    }
  }

  &-title {
    @apply mt-3 text-center text-[3vmax] font-bold uppercase text-white tracking-tight;
  }

  &-button {
    @apply mt-2 bg-indigo-500 text-white text-lg font-bold rounded-md px-4 py-2 shadow-sm shadow-indigo-600/50 hover:shadow-lg transition-shadow;
  }

  &-links {
    @apply mt-3 flex justify-center;

    &__link {
      @apply text-white hover:text-indigo-600 text-sm font-thin tracking-tight;
    }

    &__separator {
      @apply mx-1 text-indigo-500 text-sm font-thin tracking-tight;
    }
  }
}

.feature {
  @apply h-[90vh] flex justify-center md:justify-evenly items-center flex-col md:flex-row;

  &-image {
    @apply w-[37vmax] object-contain rounded-xl shadow-sm hover:shadow-md transition-all hover:scale-105 overflow-hidden;
  }

  &.reverse {
    @apply flex-col md:flex-row-reverse;
    .feature-content {
      @apply md:text-left;
    }
  }

  &-content {
    @apply flex flex-col max-w-[37vmax] mt-12 md:mt-0 text-center md:text-right;

    &__title {
      @apply text-3xl font-bold tracking-tight;
    }

    &__description {
      @apply text-lg;
    }
  }
}

.img-light {
  @apply dark:hidden;
}

.img-dark {
  @apply hidden dark:block;
}

.ending {
  @apply h-[90vh] flex justify-center items-center flex-col bg-sky-200;

  &-content {
    @apply flex flex-col text-center;

    &__title {
      @apply text-[3vmax] font-bold tracking-tight;
    }

    &__description {
      @apply text-lg;
    }
  }

  &-button {
    @apply mt-4 bg-sky-500 text-white text-lg font-bold rounded-md px-4 py-2 shadow-sm shadow-sky-600/50 hover:shadow-lg transition-shadow;
  }
}
</style>
