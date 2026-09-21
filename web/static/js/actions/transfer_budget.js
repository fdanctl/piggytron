import { sanitizeNumberText } from "../input";

export const transferBudgetActions = {
  "transfer-budget.overspent.max-value": maxAmount,
};

function maxAmount({ ele }) {
  const form = ele.closest("form");
  const isOverspend =
    form.querySelector('input[name="type"]').value === "overspend";
  if (!isOverspend) {
    return;
  }

  const amountInput = form.querySelector('input[name="amount"]');
  const currAmount = sanitizeNumberText(amountInput.value);

  const selectedAmount = sanitizeNumberText(
    ele.closest("li").querySelector(".available").innerText,
  );

  amountInput.value = Math.min(currAmount, Math.max(selectedAmount, 0));
}
