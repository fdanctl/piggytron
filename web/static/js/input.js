document.addEventListener(
  "blur",
  (e) => {
    if (e.target.matches(".input")) {
      e.target.classList.toggle("error", !e.target.validity.valid);
      const parent = e.target.parentElement.parentElement;
      if (parent.classList.contains("input-group")) {
        parent.classList.toggle("error", !e.target.validity.valid);
      }
    }
  },
  true,
);

// password
document.addEventListener("click", (e) => {
  const iconDiv = e.target.closest(".pwdToggle");
  if (!iconDiv) return;

  const pwdInput = iconDiv.parentElement.parentElement.children[0];

  if (pwdInput.type === "password") {
    pwdInput.type = "text";
  } else if (pwdInput.type === "text") {
    pwdInput.type = "password";
  }
  iconDiv.children[0].classList.toggle("hidden");
  iconDiv.children[1].classList.toggle("hidden");
});
