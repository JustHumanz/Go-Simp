<template>
  <section v-html="rendered"></section>
</template>

<script>
import { marked } from "marked"
const mdfiles = import.meta.glob("./docs/*.md", { as: "raw" })

export default {
  props: {
    page: {
      type: String,
      required: true,
    },
  },
  created() {
    this.$watch(
      () => this.$route.hash,
      async () => {
        if (!this.$route.hash) return
        await new Promise((resolve) => setTimeout(resolve, 60))
        let element = document.getElementById(this.$route.hash.replace("#", ""))

        if (element) {
          const headerOffset = 70
          const elementPosition = element.getBoundingClientRect().top
          const offsetPosition =
            elementPosition + window.pageYOffset - headerOffset

          window.scrollTo({
            top: offsetPosition,
            behavior: "smooth",
          })
        }
      },
      {
        immediate: true,
      }
    )
  },
  computed: {
    rendered() {
      if (!mdfiles[`./docs/${this.page}.md`]) return ""
      return marked(mdfiles[`./docs/${this.page}.md`])
        .replace(/(%br%)/g, '<div style="margin-top:1rem;"></div>')
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
    // @apply leading-normal whitespace-pre-line;
  }

  h1 {
    @apply text-3xl font-bold;
  }

  h2 {
    @apply text-2xl font-bold mt-6;
  }

  h3 {
    @apply text-lg font-bold mt-4;
  }

  pre {
    @apply bg-slate-200 dark:bg-slate-700 px-3 py-2 rounded-lg my-2 text-sm whitespace-pre-line;
  }

  .value-slash {
    @apply py-px px-2 bg-slate-300 dark:bg-slate-800 inline-block whitespace-nowrap rounded-md my-px font-mono;
  }

  a {
    @apply text-blue-500 dark:text-gray-300 hover:underline;
  }

  ul {
    @apply list-disc ml-5 my-2;
  }

  ol {
    @apply list-decimal ml-5 my-2;
  }

  table {
    @apply w-full table-auto border-collapse rounded-md overflow-hidden max-w-screen-sm;

    th {
      @apply bg-slate-300 dark:bg-slate-700 text-center;
    }

    td {
      @apply bg-slate-100 dark:bg-slate-500 py-px px-2 border-t-2 border-r-2 last:border-r-0 border-slate-300 dark:border-slate-700 text-center;
    }
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
