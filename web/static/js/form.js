// disable child buttons
export function disableBtns({ ele }) {
  const btns = ele.querySelectorAll(".btn");
  for (let i = 0; i < btns.length; i++) {
    btns[i].disabled = true;
  }
}
