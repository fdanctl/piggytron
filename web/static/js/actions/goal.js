import { confirmModal } from "../confirmModal";

/**
 * Registry of goal-specific actions ("goal.*").
 */
export const goalActions = {
  "goal.edit.confirm": confirmGoal,
};

/**
 * Asks the user for confirmation when a goal category is changed, warning about
 * contributions category change.
 * Auto-confirms when the category is unchanged.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The form element.
 * @param {SubmitEvent} param0.evt - The submit event.
 * @param {DOMStringMap} param0.data - Dataset of the form; `data.previousCategory`
 *   is the category before the change.
 */
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
