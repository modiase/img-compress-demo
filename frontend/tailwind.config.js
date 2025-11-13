/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#b8d4f1",
          dark: "#9bc0e8",
        },
        secondary: "#f4c7d3",
        accent: "#c9e4ca",
        pastel: {
          blue: "#b8d4f1",
          pink: "#f4c7d3",
          green: "#c9e4ca",
          yellow: "#ffd6a5",
          red: "#ffb3b3",
        },
      },
    },
  },
  plugins: [],
};
