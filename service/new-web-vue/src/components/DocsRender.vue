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

    document.body.addEventListener("click", (e) => {
      if (e.target.id === "router-link") {
        e.preventDefault()
        this.$router.push(e.target.getAttribute("href"))
      }
    })
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
    @apply mt-6 text-2xl font-bold;
  }

  h3 {
    @apply mt-4 text-lg font-bold;
  }

  pre {
    @apply my-2 whitespace-pre-line rounded-lg bg-slate-200 px-3 py-2 text-sm dark:bg-slate-700;
  }

  .value-slash {
    @apply my-px inline-block whitespace-nowrap rounded-md bg-slate-300 py-px px-2 font-mono dark:bg-slate-800;
  }

  a {
    @apply text-blue-500 hover:underline dark:text-gray-300;
  }

  ul {
    @apply my-2 ml-5 list-disc;
  }

  ol {
    @apply my-2 ml-5 list-decimal;
  }

  table {
    @apply w-full max-w-screen-sm table-auto border-collapse overflow-hidden rounded-md;

    th {
      @apply bg-slate-300 text-center dark:bg-slate-700;
    }

    td {
      @apply border-t-2 border-r-2 border-slate-300 bg-slate-100 py-px px-2 text-center last:border-r-0 dark:border-slate-700 dark:bg-slate-500;
    }
  }

  .d-channel {
    @apply font-semibold text-blue-500;
  }

  .d-role {
    @apply font-semibold text-orange-500;
  }

  .flex-center {
    @apply flex items-center justify-center;
  }
}
</style>
