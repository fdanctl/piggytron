/**
 * @param {Date} d
 */
export function formatDate(d) {
  const formatter = new Intl.DateTimeFormat("en-US", {
    month: "long",
  });
  return `${formatter.format(d)} ${d.getDate()}, ${d.getFullYear()}`;
}
