const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");

/**
 * Resolves the effective theme from the stored preference and the OS
 * preference (used when the stored preference is "system").
 *
 * @returns {"light"|"dark"} The effective theme.
 */
export function getPreferredTheme() {
  const theme = localStorage.getItem("theme") || "system";
  const prefersDark = mediaQuery.matches;
  const isDark = theme === "dark" || (theme === "system" && prefersDark);

  return isDark ? "dark" : "light";
}

/**
 * Applies the given theme to the document root and wires (or removes) the
 * media query listener as needed.
 *
 * @param {"light"|"dark"|"system"} theme - The theme to apply.
 */
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
      const preferredTheme = getPreferredTheme();
      if (preferredTheme === "dark") {
        document.documentElement.setAttribute("data-theme", "dark");
      } else {
        document.documentElement.removeAttribute("data-theme", "dark");
      }
  }
}

/**
 * Applies the selected theme with a view transition.
 *
 * @param {Object} param0 - Action payload.
 * @param {DOMStringMap} param0.data - Dataset of the option; `data.value` is
 *   the theme ("light", "dark" or "system").
 */
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

/**
 * Applies the OS color scheme to the document root. Registered as the
 * media query change listener while in "system" mode.
 *
 * @param {MediaQueryListEvent} e - The media query change event.
 */
function updateTheme(e) {
  document.documentElement.dataset.theme = e.matches ? "dark" : "light";
}
const initialTheme = localStorage.getItem("theme") || "system";
if (initialTheme === "system") {
  mediaQuery.addEventListener("change", updateTheme);
}

/**
 * Marks the clicked theme option as selected and reveals its hint text.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The clicked theme option.
 * @param {DOMStringMap} param0.data - Dataset of the option; `data.value` is
 *   the theme the option stands for.
 */
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
