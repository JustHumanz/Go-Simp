<script setup></script>

<template>
  <section class="header-section">
    <div class="header">
      <h4>Welcome to Vtbot</h4>
      <p>Vtbot is a sceduler and sharing fan-art for vtubers.</p>
      <a href="/invite" target="_blank" class="btn-invite"> ADD ME! </a>
      <div class="text-orange-200">
        <a
          href="https://github.com/JustHumanz/Go-Simp"
          target="_blank"
          class="link-source"
        >
          Source Code
        </a>
        |
        <a
          href="https://discord.gg/ydWC5knbJT"
          target="_blank"
          class="link-source"
        >
          Join Server
        </a>
      </div>
    </div>
  </section>
  <section class="agency-section">
    <div class="section">
      <h4>remember your talent in action</h4>
      <p>Scheduling live and video over 300+ vtubers and counting</p>
    </div>
  </section>
  <section class="member-section">
    <div class="section">
      <h4>Look a gift for fans</h4>
      <p>There are also fanart and h****i made for vtuber members</p>
    </div>
  </section>
  <section class="quest-section">
    <div class="section">
      <h4>There is something missing?</h4>
      <p>
        You can request your best talent on
        <a href="https://discord.gg/ydWC5knbJT" target="_blank">Discord</a>
        or
        <a
          href="https://github.com/JustHumanz/Go-Simp/issues/new/choose"
          target="_blank"
          >Github</a
        >
      </p>
    </div>
  </section>
</template>

<style lang="scss">
.header-section {
  @apply h-full py-52 bg-img-none md:bg-kizuna-ai bg-no-repeat bg-ipad-header xl:bg-post-default;

  .header {
    @apply flex flex-col items-center justify-center ml-0 md:ml-64 lg:ml-[25rem];

    h4 {
      @apply text-4xl font-semibold text-center uppercase;
    }

    p {
      @apply text-xl text-center font-light text-gray-600;
    }

    .btn-invite {
      @apply bg-orange-400 my-3 py-2 px-4 text-xl rounded-md font-semibold text-white cursor-pointer shadow-none transition-shadow hover:shadow-md hover:shadow-orange-400/75;
    }
    .link-source {
      @apply text-sm text-center text-gray-500 hover:text-gray-700 cursor-pointer;
    }
  }
}

.agency-section {
  @apply h-full py-52 bg-sky-300 bg-img-none md:bg-tokino-sora bg-no-repeat bg-ipad-agency xl:bg-right-top;

  .section {
    @apply flex flex-col items-center justify-center mr-0 md:mr-64 lg:mr-[27rem];

    h4 {
      @apply text-4xl font-semibold text-center uppercase px-14 lg:px-0;
    }

    p {
      @apply text-xl text-center font-light text-gray-600;
    }
  }
}
.member-section {
  @apply h-full py-52 bg-pink-300 bg-img-none md:bg-mito bg-no-repeat bg-ipad-header xl:bg-post-default;

  .section {
    @apply flex flex-col items-center justify-center  ml-0 md:ml-64 lg:ml-[26rem];

    h4 {
      @apply text-4xl font-semibold text-center uppercase;
    }

    p {
      @apply text-xl text-center font-light text-gray-600;
    }
  }
}
.quest-section {
  @apply h-full py-52 bg-orange-300 bg-img-none;

  .section {
    @apply flex flex-col items-center justify-center;

    h4 {
      @apply text-4xl font-semibold text-center uppercase;
    }

    p {
      @apply text-xl text-center font-light text-gray-600;

      a {
        @apply text-gray-500 hover:text-gray-700 cursor-pointer;
      }
    }
  }
}
</style>
