import { handleInputOnBlur } from "../input";
import { showToast } from "../toast";
import { goalActions } from "./goal";
import { uiActions } from "./ui";

const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      const ele = entry.target;
      dispatch(ele.dataset[eventAttributes.intersect], {
        ele,
        evt: entry,
        data: ele.dataset,
      });
    }
  });
});

document.body.addEventListener("htmx:afterOnLoad", function (ev) {
  observer.disconnect();
  const elements = document.querySelectorAll(
    `[data-${eventAttributes.intersect}]`,
  );
  elements.forEach((el) => observer.observe(el));
});

function log({ ele, evt, data }) {
  console.log(ele);
  console.log(evt);
  console.log(data);
}

const actions = {
  "dev.log.this": log,
  ...uiActions,
  ...goalActions,
};
const eventAttributes = {
  click: "action",
  input: "input",
  change: "change",
  focusin: "focus",
  focusout: "blur",
  keydown: "keydown",
  pointerdown: "pointerdown",
  "htmx:beforeRequest": "beforerequest",
  "htmx:afterRequest": "afterrequest",
  "htmx:afterOnLoad": "afteronload",
  animationend: "animationend",
  submit: "submit",
  intersect: "intersect",
  mouseover: "mouseover",
  mouseout: "mouseout",
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
  const names = actionName.trim().split(/\s+/);

  for (const n of names) {
    const action = actions[n];
    if (!action) {
      console.warn(`Unknown action: ${n}`);

      continue;
    }

    action(payload);
  }
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
  if (value === eventAttributes.intersect) {
    const elements = document.querySelectorAll(`[data-${value}]`);
    elements.forEach((el) => observer.observe(el));
  } else {
    createListener(key, value);
  }
}
