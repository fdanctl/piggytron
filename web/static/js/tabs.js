export function changeTab({ ele, data }) {
  console.log(data.tabTrigger);
  const triggers = document.querySelectorAll(
    `[data-tab-trigger="${data.tabTrigger}"]`,
  );
  const contents = document.querySelectorAll(
    `[data-tab-content="${data.tabTrigger}"]`,
  );

  triggers.forEach((t) => t.classList.remove("active"));
  ele.classList.add("active");

  contents.forEach((content) => content.classList.add("hidden"));
  const targetContent = document.getElementById(`${data.tab}-content`);
  if (targetContent) {
    targetContent.classList.remove("hidden");
  }
}
