/**
 * Builds an inline <svg>.
 *
 * @param {"success"|"warning"|"error"|"info"|"x"} icon - Icon name.
 * @param {number} [size] - Width/height in px (default 24).
 * @param {string} [color] - Stroke color (default "currentColor").
 * @param {string} [className] - CSS class for the <svg> element.
 * @returns {string} SVG markup string.
 */
export function makeSVG(
  icon,
  size = 24,
  color = "currentColor",
  className = "",
) {
  let children = "";
  switch (icon) {
    case "success":
      children = '<path d="M20 6 9 17l-5-5"></path>';
      break;

    case "warning":
      children =
        '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"></path><path d="M12 9v4"></path><path d="M12 17h.01"></path>';
      break;

    case "error":
      children =
        '<circle cx="12" cy="12" r="10"></circle><path d="m15 9-6 6"></path><path d="m9 9 6 6"></path>';
      break;

    case "info":
      children =
        '<circle cx="12" cy="12" r="10"></circle><path d="M12 16v-4"></path><path d="M12 8h.01"></path>';
      break;

    case "x":
      children = '<path d="M18 6 6 18"></path> <path d="m6 6 12 12"></path>';
      break;

    default:
      children =
        '<circle cx="12" cy="12" r="10"></circle><path d="M12 16v-4"></path><path d="M12 8h.01"></path>';
      break;
  }

  return `
	<svg
		xmlns="http://www.w3.org/2000/svg"
		width=${size}
		height=${size}
		viewBox="0 0 24 24"
		fill="none"
		stroke=${color}
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
		class=${className}
	>
		${children}
	</svg>
  `;
}
