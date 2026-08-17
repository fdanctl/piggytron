import { confirmModal } from "../confirmModal";
import { showToast } from "../toast";

/**
 * Registry of account-specific actions ("account.*").
 */
export const accountActions = {
  "account.close.confirm": closeAccountConfirm,
  "account.delete.confirm": deleteAccountConfirm,
};

/**
 * Checks if the account balance is 0, and asks the user for confirmation.
 * If the account balance is not 0 it sends a errors toast.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The form element.
 * @param {DOMStringMap} param0.data - Dataset of the form; `data.balance`
 *   the account balance.
 */
function closeAccountConfirm({ ele, data }) {
  if (data.balance > 0) {
    showToast("error", "Make the account balance 0 before closing it");
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
 * Checks if exists transactions, and asks the user for confirmation.
 * If there is transactions it sends a errors toast.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The form element.
 * @param {DOMStringMap} param0.data - Dataset of the form; `data.transactions`
 *   the account transactions count.
 */
function deleteAccountConfirm({ ele, data }) {
  if (data.transactions > 0) {
    showToast("error", "Can't delete an account with historical data");
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
