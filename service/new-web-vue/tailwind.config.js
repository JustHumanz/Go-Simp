/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme")
const colors = require("tailwindcss/colors")

module.exports = {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  darkMode: "class",
  theme: {
    screens: {
      xs: "454px",
      ...defaultTheme.screens,
    },
    extend: {
      backgroundImage: {
        smolame: "url(/assets/smolame.jpg)",
        rip: "url(/src/assets/rip.svg)",
        "img-none": "none",
      },

      backgroundPosition: {
        "ipad-header": "-12.5rem 0",
        "ipad-agency": "60vw 0%",
        "post-default": "0 0",
      },

      boxShadow: {
        // center: "0 0.5rem 1rem rgba(0, 0, 0, 0.2)",
      },

      borderWidth: {
        5: "5px",
      },

      spacing: {
        42: "10,5rem",
      },

      fontSize: {
        "1.5xl": "2.5rem",
        "2xs": ["0.65rem", "0.5rem"],
      },

      colors: {
        youtube: "#ff0000",
        bilibili: "#00a1d6",
        twitch: "#a970ff",
        twitter: "#1da1f2",
        pixiv: "#0096fa",
        // 60 30 1
        rules: {
          10: colors.blue[700],
          30: colors.blue[400],
          60: colors.white,
        },
      },

      brightness: {
        125: "1.25",
        95: "0.95",
      },
    },
  },
  plugins: [],
}
