/** @type {import('tailwindcss').Config} */

export const content = ["./web/templates/**/*.{html,js}"];
export const theme = {
  extend: {
    colors: {
      main: "#010e26",
      // main: "#d8ba98",
      background: "#c3d8e2",
      test: "#5B3739",
    },
    animation: {
      carousel: "carousel 40s linear infinite",
    },
    keyframes: {
      carousel: {
        to: {
          transform: "translateX(-50%)",
        },
      },
    },
  },
  fontFamily: {
    body: ["Outfit", "sans-serif"],
  },
};
export const plugins = [];
