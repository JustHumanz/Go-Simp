// eslint-disable-next-line no-undef
module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      screens: {
        xs: "454px",
      },
      backgroundImage: {
        smolame: "url(/src/assets/smolame.jpg)",
        rip: "url(/src/assets/rip.svg)",
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

      spacing: {
        42: "10,5rem",
      },

      fontSize: {
        "1.5xl": "2.5rem",
      },

      colors: {
        youtube: "#ff0000",
        bilibili: "#00a1d6",
        twitch: "#a970ff",
        twitter: "#1da1f2",
        pixiv: "#0096fa",
      },
    },
  },
  plugins: [],
}
