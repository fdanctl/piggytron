import { confirmModal } from "../confirmModal";
import { showToast } from "../toast";

/**
 * Registry of account-specific actions ("category.*").
 */
export const categoryActions = {
  "category.delete.confirm": deleteCategoryConfirm,
  "category.archived.toggle": toggleArchived,
};

/**
 * Checks if exists transactions, and asks the user for confirmation.
 * If there is transactions it sends a errors toast.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The form element.
 * @param {DOMStringMap} param0.data - Dataset of the form; `data.transactions`
 *   the account transactions count.
 */
function deleteCategoryConfirm({ ele, data }) {
  if (data.transactions > 0) {
    showToast("error", "Can't delete a category with historical data");
    return;
  }

  const config = {
    title: "Warning",
    message: "This action is irreversible. Do you want to procced?",
    acceptText: "Yes",
    refuseText: "Cancel",
  };
  confirmModal(config).then(function (result) {
    if (result) {
      htmx.trigger(ele, "confirmed");
    }
  });
}

/**
 * Show/hides archived categories.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The form element.
 */
function toggleArchived({ ele }) {
  console.log(ele);
  document
    .getElementById("categories-container")
    .classList.toggle("show-archived");
}
