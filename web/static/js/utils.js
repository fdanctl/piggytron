/**
 * Formats a date as "Month Day, Year" (e.g. "January 5, 2026").
 *
 * @param {Date} d - The date to format.
 * @returns {string} The formatted date string.
 */
export function formatDate(d) {
  // TODO: locale config
  const formatter = new Intl.DateTimeFormat("en-US", {
    month: "long",
  });
  return `${formatter.format(d)} ${d.getDate()}, ${d.getFullYear()}`;
}
