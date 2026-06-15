import { confirmModal } from "../htmx";

export const goalActions = {
  "goal.edit.confirm": confirmGoal,
};

function confirmGoal({ ele, evt, data }) {
  const formData = new FormData(evt.target);
  if (formData.get("category") == data.previousCategory) {
    htmx.trigger(ele, "confirmed");
    return;
  }

  const contributionCount =
    document.getElementById("contribution-count").innerText;
  if (
    Number.isNaN(Number(contributionCount)) ||
    Number(contributionCount) <= 0
  ) {
    htmx.trigger(ele, "confirmed");
    return;
  }

  const categoryName = document.querySelector(
    `[data-value="${formData.get("category")}"]`,
  ).innerText;
  const prevCategoryName = document.querySelector(
    `[data-value="${data.previousCategory}"]`,
  ).innerText;

  const config = {
    title: "Warning",
    message: `${contributionCount} contributions will change category. ${prevCategoryName} > ${categoryName}.`,
    acceptText: "Proceed",
    refuseText: "Cancel",
  };
  confirmModal(config).then(function (result) {
    if (result) {
      htmx.trigger(ele, "confirmed");
    }
  });
}
