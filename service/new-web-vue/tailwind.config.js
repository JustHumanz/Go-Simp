// eslint-disable-next-line no-undef
module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      backgroundImage: {
        "kizuna-ai": "url(/src/assets/kizuna-ai.png)",
        "tokino-sora": "url(/src/assets/tokino-sora.png)",
        mito: "url(/src/assets/mito.png)",
        smolame: "url(/src/assets/smolame.png)",
        "img-none": "none",
      },

      backgroundPosition: {
        "ipad-header": "-12.5rem 0",
        "ipad-agency": "60vw 0%",
        "post-default": "0 0",
      },

      boxShadow: {
        center: "0 0.5rem 1rem rgba(0, 0, 0, 0.2)",
      },

      colors: {
        youtube: "#ff0000",
        bilibili: "#00a1d6",
        twitch: "#6441a5",
        twitter: "#1da1f2",
      },
    },
  },
  plugins: [],
}
