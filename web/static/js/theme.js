const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

export function getPreferedTheme() {
  const theme = localStorage.getItem("theme") || "system";
  const prefersDark = mediaQuery.matches;
  const isDark = theme === "dark" || (theme === "system" && prefersDark);

  return isDark ? "dark" : "light";
}

function switchTheme(theme) {
  switch (theme) {
    case "light":
      mediaQuery.removeEventListener("change", updateTheme);
      document.documentElement.removeAttribute("data-theme");
      return;

    case "dark":
      mediaQuery.removeEventListener("change", updateTheme);
      document.documentElement.setAttribute("data-theme", "dark");
      return;

    case "system":
      mediaQuery.addEventListener("change", updateTheme);
      const preferedTheme = getPreferedTheme();
      if (preferedTheme === "dark") {
        document.documentElement.setAttribute("data-theme", "dark");
      }
  }
}

export function selectTheme({ data }) {
  localStorage.setItem("theme", data.value);
  document.startViewTransition({
    update() {
      document.documentElement.classList.add("theme-transition");
      switchTheme(data.value);
    },
    types: ["theme"],
  });
}

function updateTheme(e) {
  document.documentElement.dataset.theme = e.matches ? "dark" : "light";
}

const initialTheme = localStorage.getItem("theme") || "system";
if (initialTheme === "system") {
  mediaQuery.addEventListener("change", updateTheme);
}

// theme option
export function themeOptionSelect({ ele, data }) {
  const parent = ele.parentElement;
  const opts = parent.querySelectorAll(".theme-option");

  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("theme-option__selected");
  }
  ele.classList.add("theme-option__selected");

  const hints = parent.nextElementSibling.querySelectorAll("[data-hint]");
  for (let i = 0; i < opts.length; i++) {
    const e = hints[i];
    if (e.dataset.hint === data.value) {
      e.classList.remove("hidden");
    } else {
      e.classList.add("hidden");
    }
  }
}
