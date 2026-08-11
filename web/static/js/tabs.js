/**
 * Activates the clicked tab trigger and shows the matching tab content.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The clicked tab trigger.
 * @param {DOMStringMap} param0.data - Dataset of the trigger; `data.tab` names
 *   the content id and `data.tabTrigger` groups triggers/contents.
 */
export function changeTab({ ele, data }) {
  const triggers = document.querySelectorAll(
    `[data-tab-trigger="${data.tabTrigger}"]`,
  );
  const contents = document.querySelectorAll(
    `[data-tab-content="${data.tabTrigger}"]`,
  );

  triggers.forEach((t) => t.classList.remove("tab-trigger--active"));
  ele.classList.add("tab-trigger--active");

  contents.forEach((content) => content.classList.add("hidden"));
  const targetContent = document.getElementById(`${data.tab}-content`);
  if (targetContent) {
    targetContent.classList.remove("hidden");
  }
}
