<template>
  <pre id="content-md">
  # Quick Start

  ### Setup Live Streaming
  ```slash
  /setup channel-type livestream slash{channel-name: channel{#hololive}} slash{vtuber-group: hololive}
  ```
  router-link{« Read More}(/docs/configuration#setup-live-streaming)

  ### Update Stage
  ```slash
  /channel-update slash{channel-name: channel{#hololive}}
  ```
  router-link{« Read More}(/docs/configuration#update-stage)

  ### Tag Roles
  ```slash
  /tag-role slash{role-name: role{@Holo Simps}} slash{vtuber-group: hololive}
  ```
  this means that the bot will mention role{@Holo Simps} when any new hololive fan arts or live streams are uploaded. 
  router-link{« Read More}(/docs/roles-and-taging#tag-roles)

  ### Get Vtuber Group and Vtuber Name

  </pre>
</template>

<script>
import { marked } from "marked"

export default {
  data() {
    return {
      narkdown: "",
      rendered: null,
    }
  },
  mounted() {
    document.title = "Quick Start - Documentation"
    this.render()
  },
  methods: {
    render() {
      this.rendered = marked(document.getElementById("content-md").innerText)
        .replace(/channel{([^}]+)}/g, '<span class="d-channel">$1</span>')
        .replace(/role{([^}]+)}/g, '<span class="d-role">$1</span>')
        .replace(/slash{([^}]+)}/g, '<span class="value-slash">$1</span>')
        .replace(
          /router-link{([^}]+)}\(([a-z\/?=&#\-]+)\)/g,
          '<a id="router-link" href="$2">$1</a>'
        )
      document.getElementById("content-md").innerHTML = this.rendered
    },
  },
}
</script>

<style lang="scss">
.content {
  font-family: "Nunito", sans-serif;
  @apply text-base leading-4;

  & > * {
    @apply leading-normal whitespace-pre-line;
  }

  h1 {
    @apply text-3xl font-bold;
  }

  h2 {
    @apply text-2xl font-bold ;
  }

  h3 {
    @apply text-xl font-bold mt-7;
  }

  pre {
    @apply bg-slate-200 dark:bg-slate-700 px-3 py-2 rounded-lg whitespace-nowrap overflow-y-scroll;

    .value-slash {
      @apply py-px px-2 bg-slate-300 dark:bg-slate-800 inline-block rounded-md my-px;

      &:last-child {
        @apply mr-3;
      }
    }
  }

  a {
    @apply text-blue-500 dark:text-gray-300 hover:underline;
  }

  .d-channel {
    @apply text-blue-500 font-semibold;
  }

  .d-role {
    @apply text-orange-500 font-semibold;
  }
}
</style>
