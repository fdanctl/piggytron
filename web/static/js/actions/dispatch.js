import { handleInputOnBlur } from "../input";
import { showToast } from "../toast";
import { goalActions } from "./goal";
import { uiActions } from "./ui";

const actions = {
  ...uiActions,
  ...goalActions,
};
const eventAttributes = {
  click: "action",
  input: "input",
  change: "change",
  focusin: "focus",
  focusout: "focusout",
  keydown: "keydown",
  "htmx:beforeRequest": "beforerequest",
  "htmx:afterRequest": "afterrequest",
  animationend: "animationend",
  submit: "submit",
};

document.addEventListener("focusout", (evt) => {
  if (!(evt.target instanceof HTMLInputElement)) return;

  handleInputOnBlur(evt.target);
});

window.addEventListener("offline", () => {
  showToast("error", "You are offline");
});

window.addEventListener("online", () => {
  showToast("success", "Internet connection restored");
});

function dispatch(actionName, payload) {
  const action = actions[actionName];

  if (!action) {
    console.warn(`Unknown action: ${actionName}`);

    return;
  }

  action(payload);
}

function createListener(eventName, dataAttr) {
  document.addEventListener(eventName, (evt) => {
    const ele = evt.target.closest(`[data-${dataAttr}]`);

    if (!ele) return;

    dispatch(ele.dataset[dataAttr], {
      ele,
      evt,
      data: ele.dataset,
    });
  });
}

for (let [key, value] of Object.entries(eventAttributes)) {
  createListener(key, value);
}
