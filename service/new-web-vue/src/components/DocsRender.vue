<template>
  <section v-html="rendered"></section>
</template>

<script>
import { marked } from "marked"
const mdfiles = import.meta.glob("./docs/*.md", { assert: { type: "raw" } })

export default {
  props: {
    page: {
      type: String,
      required: true,
    },
  },
  computed: {
    rendered() {
      if (!mdfiles[`./docs/${this.page}.md`]) return ""
      return marked(mdfiles[`./docs/${this.page}.md`])
        .replace(/channel{([^}]+)}/g, '<span class="d-channel">$1</span>')
        .replace(/role{([^}]+)}/g, '<span class="d-role">$1</span>')
        .replace(/slash{([^}]+)}/g, '<span class="value-slash">$1</span>')
        .replace(/center{([^}]+)}/g, '<div class="flex-center">$1</div>')
        .replace(/small{([^}]+)}/g, "<small>$1</small>")
        .replace(
          /router-link{([^}]+)}\(([a-z\/?=&#\-]+)\)/g,
          '<a id="router-link" href="$2">$1</a>'
        )
    },
  },
}
</script>

<style lang="scss">
.content {
  & > * {
    @apply leading-normal whitespace-pre-line;
  }

  h1 {
    @apply text-3xl font-bold;
  }

  h2 {
    @apply text-2xl font-bold mt-6;
  }

  h3 {
    @apply text-xl font-bold mt-4;
  }

  pre {
    @apply bg-slate-200 dark:bg-slate-700 px-3 py-2 rounded-lg whitespace-nowrap overflow-y-scroll my-2;

    .value-slash {
      &:last-child {
        @apply mr-3;
      }
    }
  }

  .value-slash {
    @apply py-px px-2 bg-slate-300 dark:bg-slate-800 inline-block whitespace-nowrap rounded-md my-px font-mono;
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

  .flex-center {
    @apply flex justify-center items-center;
  }
}
</style>
